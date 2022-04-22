package orderfeeprofitservice

import (
	"errors"
	"github.com/copo888/transaction_service/common/constants"
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
		if err = calculateProfitLoop(db, &calculateProfit, &orderFeeProfits, true); err != nil {
			logx.Errorf("計算利潤錯誤: %s ", err.Error())
			return err
		}
		logx.Infof("計算利潤: %#v ", orderFeeProfits)

		if err = updateOrderByIsCalculateProfit(db, calculateProfit.OrderNo); err != nil {
			logx.Errorf("計算利潤錯誤: %s ", err.Error())
			return err
		}
		return
	})
}

// CalculateOrderProfitForRecover 追回單計算利潤 (方式不同 開專用func)
func CalculateOrderProfitForRecover(db *gorm.DB, calculateProfit types.CalculateProfit, isCalculateCommission bool) (err error) {
	return db.Transaction(func(db *gorm.DB) (err error) {
		var orderFeeProfits []types.OrderFeeProfit
		if err = calculateProfitLoop(db, &calculateProfit, &orderFeeProfits, isCalculateCommission); err != nil {
			logx.Errorf("計算利潤錯誤: %s ", err.Error())
			return err
		}
		logx.Infof("計算利潤: %#v ", orderFeeProfits)

		if err = updateOrderByIsCalculateProfit(db, calculateProfit.OrderNo); err != nil {
			logx.Errorf("計算利潤錯誤: %s ", err.Error())
			return err
		}
		return
	})
}

// calculateProfitLoop 計算利潤迴圈
func calculateProfitLoop(db *gorm.DB, calculateProfit *types.CalculateProfit, orderFeeProfits *[]types.OrderFeeProfit, isCalculateCommission bool) (err error) {

	var merchant *types.Merchant
	var agentLayerCode string
	var agentParentCode string

	// 1. 不是算系統利潤時 要取當前計算商戶(或代理商戶)
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
		// MerchantCode 不是 00000000 要取商戶費率
		if err = setMerchantFee(db, calculateProfit, &orderFeeProfit); err != nil {
			return
		}
	} else {
		// MerchantCode 是 00000000 要取渠道費率
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

	// 6.判斷是否計算下一筆利潤傭金
	*orderFeeProfits = append(*orderFeeProfits, orderFeeProfit)
	if calculateProfit.MerchantCode == "00000000" {
		// 已計算系統利潤 結束迴圈
		return
	} else if agentParentCode != "" && isCalculateCommission {
		// 有上層代理 且 需計算傭金 要接下去算代理傭金
		calculateProfit.MerchantCode = agentParentCode
		return calculateProfitLoop(db, calculateProfit, orderFeeProfits, isCalculateCommission)
	} else {
		// 沒上層代理 開始計算系統利潤
		calculateProfit.MerchantCode = "00000000"
		return calculateProfitLoop(db, calculateProfit, orderFeeProfits, isCalculateCommission)
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

func updateOrderByIsCalculateProfit(db *gorm.DB, orderNo string) error {
	return db.Table("tx_orders").
		Where("order_no = ?", orderNo).
		Updates(map[string]interface{}{"is_calculate_profit": constants.IS_CALCULATE_PROFIT_YES}).Error
}

