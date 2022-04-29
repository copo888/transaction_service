package logic

import (
	"context"
	"errors"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"gorm.io/gorm"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawOrderTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWithdrawOrderTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawOrderTransactionLogic {
	return &WithdrawOrderTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WithdrawOrderTransactionLogic) WithdrawOrderTransaction(in *transactionclient.WithdrawOrderRequest) (*transactionclient.WithdrawOrderResponse, error) {
	db := l.svcCtx.MyDB

	transferAmount := utils.FloatAdd(in.OrderAmount, in.HandlingFee)

	// 初始化订单
	txOrder := &types.Order{
		MerchantCode:         in.MerchantCode,
		MerchantOrderNo:      "COPO_" + in.OrderNo,
		OrderNo:              in.OrderNo,
		Type:                 constants.ORDER_TYPE_XF,
		Status:               constants.WAIT_PROCESS,
		IsMerchantCallback:   constants.IS_MERCHANT_CALLBACK_NOT_NEED,
		IsLock:               constants.IS_LOCK_NO,
		IsCalculateProfit:    constants.IS_CALCULATE_PROFIT_NO,
		PersonProcessStatus:  constants.PERSON_PROCESS_STATUS_NO_ROCESSING,
		CreatedBy:            in.UserAccount,
		UpdatedBy:            in.UserAccount,
		BalanceType:          "XFB",
		ChannelCode:          in.CurrencyCode,
		TransferAmount:       transferAmount,
		OrderAmount:          in.OrderAmount,
		TransferHandlingFee:  in.HandlingFee,
		HandlingFee:          in.HandlingFee,
		MerchantBankAccount:  in.MerchantBankeAccount,
		MerchantBankNo:       in.MerchantBankNo,
		MerchantBankProvince: in.MerchantBankProvince,
		MerchantBankCity:     in.MerchantBankCity,
		MerchantAccountName:  in.MerchantAccountName,
		CurrencyCode:         in.CurrencyCode,
		Source:               in.Source,
		PageUrl:              in.PageUrl,
		NotifyUrl:            in.NotifyUrl,
	}
	if len(in.MerchantOrderNo) > 0 {
		txOrder.MerchantOrderNo = in.MerchantOrderNo
	}

	tx := db.Begin()
	// 新增收支记录，更新商户余额
	updateBalance := types.UpdateBalance{
		MerchantCode:    txOrder.MerchantCode,
		CurrencyCode:    txOrder.CurrencyCode,
		OrderNo:         txOrder.OrderNo,
		MerchantOrderNo: txOrder.MerchantOrderNo,
		OrderType:       txOrder.Type,
		TransactionType: "11",
		BalanceType:     txOrder.BalanceType,
		TransferAmount:  txOrder.TransferAmount,
		CreatedBy:       in.UserAccount,
	}
	if in.Source == constants.API {
		isBlock, _ := model.NewBankBlockAccount(tx).CheckIsBlockAccount(txOrder.MerchantBankAccount)
		if isBlock { //银行账号为黑名单
			logx.Infof("交易账户%s-%s在黑名单内", txOrder.MerchantAccountName, txOrder.MerchantBankNo)
			updateBalance.TransferAmount = 0                           // 使用0元前往钱包扣款
			txOrder.ErrorType = constants.ERROR6_BANK_ACCOUNT_IS_BLACK //交易账户为黑名单
			txOrder.ErrorNote = constants.BANK_ACCOUNT_IS_BLACK        //失败原因：黑名单交易失败
			txOrder.Status = constants.FAIL                            //状态:失败
			txOrder.Fee = 0                                            //写入本次手续费(未发送到渠道的交易，都设为0元)
			txOrder.HandlingFee = 0
			//transAt = types.JsonTime{}.New()
			logx.Infof("商户 %s，代付订单 %#v ，交易账户为黑名单", txOrder.MerchantCode, txOrder)
		}
	}
	//更新钱包且新增商户钱包异动记录
	merchantBalanceRecord, err1 := merchantbalanceservice.UpdateXFBalance_Debit(tx, &updateBalance)
	if err1 != nil {
		logx.Errorf("商户:%s，更新錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err1.Error(), updateBalance)
		//TODO  IF 更新钱包错误是response.DATABASE_FAILURE THEN return SYSTEM_ERROR
		tx.Rollback()
		return nil, err1
	} else {
		logx.Infof("下发提单 %s，錢包扣款成功", merchantBalanceRecord.OrderNo)
		txOrder.BeforeBalance = merchantBalanceRecord.BeforeBalance // 商戶錢包異動紀錄
		txOrder.Balance = merchantBalanceRecord.AfterBalance
	}

	// 创建订单
	if err3 := tx.Table("tx_orders").Create(&types.OrderX{
		Order: *txOrder,
	}).Error; err3 != nil {
		logx.Errorf("新增下发提单失败，商户号: %s, 订单号: %s, err : %s", txOrder.MerchantCode, txOrder.OrderNo, err3.Error())
		tx.Rollback()
		return nil, err3
	}

	//計算商戶利潤
	calculateProfit := types.CalculateProfit{
		MerchantCode:        txOrder.MerchantCode,
		OrderNo:             txOrder.OrderNo,
		Type:                txOrder.Type,
		CurrencyCode:        txOrder.CurrencyCode,
		BalanceType:         txOrder.BalanceType,
		OrderAmount:         txOrder.ActualAmount,
	}
	if err4 := l.calculateOrderProfit(tx, calculateProfit); err4 != nil {
		logx.Errorf("计算下发利润失败，商户号: %s, 订单号: %s, err : %s", txOrder.MerchantCode, txOrder.OrderNo, err4.Error())
		tx.Rollback()
		return nil, err4
	}

	if err4 := tx.Commit().Error; err4 != nil {
		logx.Errorf("最终新增下发提单失败，商户号: %s, 订单号: %s, err : %s", txOrder.MerchantCode, txOrder.OrderNo, err4.Error())
		tx.Rollback()
		return nil, err4
	}

	// 新單新增訂單歷程 (不抱錯) TODO: 異步??
	if err5 := db.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     txOrder.OrderNo,
			Action:      "PLACE_ORDER",
			UserAccount: in.MerchantCode,
			Comment:     "",
		},
	}).Error; err5 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err5.Error())
	}

	return &transactionclient.WithdrawOrderResponse{
		OrderNo: txOrder.OrderNo,
	}, nil
}

func (l *WithdrawOrderTransactionLogic) calculateOrderProfit(db *gorm.DB, calculateProfit types.CalculateProfit) (err error) {
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

	// 3.设定商户费率
	var merchantCurrency *types.MerchantCurrency
	if err = db.Table("mc_merchant_currencies").
		Where("merchant_code = ? AND currency_code = ?", calculateProfit.MerchantCode, calculateProfit.CurrencyCode).
		Take(&merchantCurrency).Error;
		errors.Is(err, gorm.ErrRecordNotFound) {
		return errorz.New(response.RATE_NOT_CONFIGURED, err.Error())
	} else if err != nil {
		return errorz.New(response.DATABASE_FAILURE, err.Error())
	}
	orderFeeProfit.HandlingFee = merchantCurrency.WithdrawHandlingFee
	//  交易手續費總額 = 訂單金額 + 手續費
	orderFeeProfit.TransferHandlingFee =
		utils.FloatAdd(calculateProfit.OrderAmount, orderFeeProfit.HandlingFee)

	// 4. 除存费率
	if err = db.Table("tx_orders_fee_profit").Create(&types.OrderFeeProfitX{
		OrderFeeProfit: orderFeeProfit,
	}).Error; err != nil {
		return errorz.New(response.DATABASE_FAILURE, err.Error())
	}
	return nil
}