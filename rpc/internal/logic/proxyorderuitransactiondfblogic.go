package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/service/orderfeeprofitservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"gorm.io/gorm"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProxyOrderUITransactionDFBLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProxyOrderUITransactionDFBLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyOrderUITransactionDFBLogic {
	return &ProxyOrderUITransactionDFBLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProxyOrderUITransactionDFBLogic) ProxyOrderUITransaction_DFB(in *transactionclient.ProxyOrderUIRequest) (resp *transactionclient.ProxyOrderUIResponse, err error) {
	req := in.ProxyOrderUI
	rate := in.MerchantOrderRateListView
	merchantBalanceRecord := types.MerchantBalanceRecord{}

	txOrder := &types.Order{
		MerchantCode:         req.MerchantCode,
		CreatedBy:            req.UserAccount,
		MerchantOrderNo:      "COPO_" + req.OrderNo,
		OrderNo:              req.OrderNo,
		OrderAmount:          req.OrderAmount,
		BalanceType:          constants.DF_BALANCE,
		Type:                 constants.ORDER_TYPE_DF,
		Status:               constants.WAIT_PROCESS,
		Source:               constants.UI,
		IsMerchantCallback:   constants.MERCHANT_CALL_BACK_DONT_USE,
		IsCalculateProfit:    constants.IS_CALCULATE_PROFIT_NO,
		PersonProcessStatus:  constants.PERSON_PROCESS_STATUS_NO_ROCESSING,
		MerchantBankAccount:  req.MerchantBankAccount,
		MerchantBankNo:       req.MerchantBankNo,
		MerchantBankName:     req.MerchantBankName,
		MerchantAccountName:  req.MerchantAccountName,
		MerchantBankProvince: req.MerchantBankProvince,
		MerchantBankCity:     req.MerchantBankCity,
		CurrencyCode:         req.CurrencyCode,
		ChannelCode:          rate.ChannelCode,
		ChannelPayTypesCode:  rate.ChannelPayTypesCode,
		Fee:                  rate.MerFee,
		HandlingFee:          rate.MerHandlingFee,
		PayTypeCode:          rate.PayTypeCode,
		PayTypeCodeNum:       rate.PayTypeCodeNum,
		PayTypeNum:           rate.PayTypeCode + rate.PayTypeCodeNum,
		IsLock:               "0",
	}

	// 新增收支记录，与更新商户余额(商户账户号是黑名单，把交易金额为设为 0)
	updateBalance := &types.UpdateBalance{
		MerchantCode:    txOrder.MerchantCode,
		CurrencyCode:    txOrder.CurrencyCode,
		OrderNo:         txOrder.OrderNo,
		OrderType:       txOrder.Type,
		PayTypeCode:     txOrder.PayTypeCode,
		PayTypeCodeNum:  txOrder.PayTypeCodeNum,
		TransferAmount:  txOrder.TransferAmount,
		TransactionType: "11", //異動類型 (1=收款; 2=解凍; 3=沖正; 11=出款 ; 12=凍結)
		BalanceType:     constants.DF_BALANCE,
		CreatedBy:       txOrder.MerchantCode,
	}

	//交易金额 = 订单金额 + 商户手续费
	txOrder.TransferAmount = utils.FloatAdd(txOrder.OrderAmount, txOrder.TransferHandlingFee)
	updateBalance.TransferAmount = txOrder.TransferAmount //扣款依然傳正值

	// 判断单笔最大最小金额
	if  txOrder.OrderAmount < rate.SingleMinCharge {
		//金额超过上限
		logx.Errorf("錯誤:代付金額未達下限")
		return nil, errorz.New(response.ORDER_AMOUNT_LIMIT_MAX)
	} else if txOrder.OrderAmount > rate.SingleMaxCharge {
		//下发金额未达下限
		logx.Errorf("錯誤:代付金額超過上限")
		return nil, errorz.New(response.ORDER_AMOUNT_LIMIT_MIN)
	}


	if err = l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

		//更新钱包且新增商户钱包异动记录
		if merchantBalanceRecord, err = merchantbalanceservice.UpdateDFBalance_Debit(db, updateBalance); err != nil {
			logx.Errorf("商户:%s，更新錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err.Error(), updateBalance)
			return errorz.New(response.SYSTEM_ERROR, err.Error())
		} else {
			logx.Infof("代付UI提单 %s，錢包扣款成功", merchantBalanceRecord.OrderNo)
			txOrder.BeforeBalance = merchantBalanceRecord.BeforeBalance // 商戶錢包異動紀錄
			txOrder.Balance = merchantBalanceRecord.AfterBalance
		}

		// 创建订单
		if err = db.Table("tx_orders").Create(&types.OrderX{
			Order:   *txOrder,
			}).Error; err != nil {
			logx.Errorf("新增代付UI提单失败，商户号: %s, 订单号: %s, err : %s", txOrder.MerchantCode, txOrder.OrderNo, err.Error())
			return
		}

		return nil
	}); err != nil {
		return
	}

	// 計算利潤(不報錯) TODO: 異步??
	if err4 := orderfeeprofitservice.CalculateOrderProfit(l.svcCtx.MyDB, types.CalculateProfit{
		MerchantCode:        txOrder.MerchantCode,
		OrderNo:             txOrder.OrderNo,
		Type:                txOrder.Type,
		CurrencyCode:        txOrder.CurrencyCode,
		BalanceType:         txOrder.BalanceType,
		ChannelCode:         txOrder.ChannelCode,
		ChannelPayTypesCode: txOrder.ChannelPayTypesCode,
		OrderAmount:         txOrder.OrderAmount,
	}); err4 != nil {
		logx.Error("計算利潤出錯:%s", err4.Error())
	}

	// 新單新增訂單歷程 (不抱錯) TODO: 異步??
	if err5 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     txOrder.OrderNo,
			Action:      "PLACE_ORDER",
			UserAccount: req.UserAccount,
			Comment:     "",
		},
	}).Error; err5 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err5.Error())
	}

	resp = &transactionclient.ProxyOrderUIResponse{
		ProxyOrderNo: txOrder.OrderNo,
	}

	return resp, nil
}
