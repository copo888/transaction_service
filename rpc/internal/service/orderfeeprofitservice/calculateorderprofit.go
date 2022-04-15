package orderfeeprofitservice

import (
	"errors"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// CalculateOrderProfit 計算利潤 (下發不可用) (回傳 商戶,代理,系統 費率利潤)
func CalculateOrderProfit(db *gorm.DB, calculateProfit types.CalculateProfit) (err error) {
	return db.Transaction(func(db *gorm.DB) (err error) {
		var orderFeeProfits []types.OrderFeeProfit
		if err = calculateProfitLoop(db, &calculateProfit, &orderFeeProfits); err != nil {
			logx.Errorf("計算利潤錯誤: %s ", err.Error())
			return err
		}
		logx.Infof("計算利潤: %#v ", orderFeeProfits)
		//TODO: tx_order 多個欄位判斷是否已計算利潤
		return
	})
}

// calculateProfitLoop 計算利潤迴圈
func calculateProfitLoop(db *gorm.DB, calculateProfit *types.CalculateProfit, orderFeeProfits *[]types.OrderFeeProfit) (err error) {

	var merchant *types.Merchant
	var agentLayerCode string
	var agentParentCode string

	// 1. 取得商戶
	if calculateProfit.MerchantCode != "00000000" {
		if err = db.Table("mc_merchants").Where("code = ?", calculateProfit.MerchantCode).Take(&merchant).Error; err != nil {
			return errorz.New(response.DATABASE_FAILURE, err.Error())
		}
		agentLayerCode = merchant.AgentLayerCode
		agentParentCode = merchant.AgentParentCode
	}

	// 2. 設定初始資料
	orderFeeProfit := types.OrderFeeProfit{
		OrderNo:             calculateProfit.OrderNo,
		MerchantCode:        calculateProfit.MerchantCode,
		AgentLayerNo:        agentLayerCode,
		AgentParentCode:     agentParentCode,
		BalanceType:         calculateProfit.BalanceType,
		Fee:                 0,
		HandlingFee:         0,
		TransferHandlingFee: 0,
		ProfitAmount:        0,
	}

	// 3. 設置手續費
	if calculateProfit.MerchantCode != "00000000" {
		if err = setMerchantFee(db, calculateProfit, &orderFeeProfit); err != nil {
			return
		}
	} else {
		if err = setChannelFee(db, calculateProfit, &orderFeeProfit); err != nil {
			return
		}
	}

	// 4. 計算利潤 (上一筆交易手續費 - 當前交易手續費 = 利潤)
	if len(*orderFeeProfits) > 0 {
		previousFeeProfits := (*orderFeeProfits)[len(*orderFeeProfits)-1]
		orderFeeProfit.ProfitAmount = previousFeeProfits.TransferHandlingFee - orderFeeProfit.TransferHandlingFee
	}

	// 5. 保存利潤
	if err = db.Table("tx_orders_fee_profit").Create(&types.OrderFeeProfitX{
		OrderFeeProfit: orderFeeProfit,
	}).Error; err != nil {
		return errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 6.判斷是否計算下一筆利潤
	*orderFeeProfits = append(*orderFeeProfits, orderFeeProfit)
	if calculateProfit.MerchantCode == "00000000" {
		// 已計算系統利潤 END
		return
	} else if agentParentCode != "" {
		// 有上層代理
		calculateProfit.MerchantCode = agentParentCode
		return calculateProfitLoop(db, calculateProfit, orderFeeProfits)
	} else {
		// 沒上層代理 開始計算系統利潤
		calculateProfit.MerchantCode = "00000000"
		return calculateProfitLoop(db, calculateProfit, orderFeeProfits)
	}
}

// 設置費率
func setMerchantFee(db *gorm.DB, calculateProfit *types.CalculateProfit, orderFeeProfit *types.OrderFeeProfit) (err error) {

	var merchantChannelRate *types.MerchantChannelRate

	if err = db.Table("mc_merchant_channel_rate").
		Where("merchant_code = ? AND channel_pay_types_code = ?", orderFeeProfit.MerchantCode, calculateProfit.ChannelPayTypesCode).
		Take(&merchantChannelRate).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errorz.New(response.RATE_NOT_CONFIGURED, err.Error())
	} else if err != nil {
		return errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	if calculateProfit.Type == "DF" {
		// "代付" 只收手續費
		orderFeeProfit.HandlingFee = merchantChannelRate.HandlingFee

	} else if calculateProfit.Type == "ZF" || calculateProfit.Type == "NC" {
		// "支付,內充" 收計算費率&手續費
		orderFeeProfit.HandlingFee = merchantChannelRate.HandlingFee
		orderFeeProfit.Fee = merchantChannelRate.Fee

	} else {
		return errorz.New(response.ILLEGAL_PARAMETER, err.Error())
	}

	//  交易手續費總額 = 訂單金額 / 100 * 費率 + 手續費
	orderFeeProfit.TransferHandlingFee =
		utils.FloatAdd(utils.FloatMul(utils.FloatDiv(calculateProfit.OrderAmount, 100), orderFeeProfit.Fee), orderFeeProfit.HandlingFee)

	return
}

// 設置費率
func setChannelFee(db *gorm.DB, calculateProfit *types.CalculateProfit, orderFeeProfit *types.OrderFeeProfit) (err error) {
	var channelPayType *types.ChannelPayType

	if err = db.Table("ch_channel_pay_types").
		Where("code = ?", calculateProfit.ChannelPayTypesCode).
		Take(&channelPayType).Error; err != nil {
		return errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	if calculateProfit.Type == "DF" {
		// "代付" 只收手續費
		orderFeeProfit.HandlingFee = channelPayType.HandlingFee

	} else if calculateProfit.Type == "ZF" || calculateProfit.Type == "NC" {
		// "支付,內充" 收計算費率&手續費
		orderFeeProfit.HandlingFee = channelPayType.HandlingFee
		orderFeeProfit.Fee = channelPayType.Fee

	} else {
		return errorz.New(response.ILLEGAL_PARAMETER, err.Error())
	}

	//  交易手續費總額 = 訂單金額 / 100 * 費率 + 手續費
	orderFeeProfit.TransferHandlingFee =
		utils.FloatAdd(utils.FloatMul(utils.FloatDiv(calculateProfit.OrderAmount, 100), orderFeeProfit.Fee), orderFeeProfit.HandlingFee)

	return
}
