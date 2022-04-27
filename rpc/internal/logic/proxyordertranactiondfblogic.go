package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/service/orderfeeprofitservice"
	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ProxyOrderTranactionDFBLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProxyOrderTranactionDFBLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyOrderTranactionDFBLogic {
	return &ProxyOrderTranactionDFBLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProxyOrderTranactionDFBLogic) ProxyOrderTranaction_DFB(in *transactionclient.ProxyOrderRequest) (resp *transactionclient.ProxyOrderResponse, err error) {
	tx := l.svcCtx.MyDB
	req := in.Req
	rate := in.Rate

	// 依商户是否给回调网址，决定是否回调商户flag
	var isMerchantCallback string //0：否、1:是、2:不需回调
	//var calBackStatus string //渠道回調狀態(0:處理中1:成功2:失敗)
	merchantBalanceRecord := types.MerchantBalanceRecord{}
	if req.NotifyUrl != "" {
		isMerchantCallback = constants.MERCHANT_CALL_BACK_NO
	} else {
		isMerchantCallback = constants.MERCHANT_CALL_BACK_DONT_USE
	}
	//初始化订单
	txOrder := &types.Order{
		OrderNo:              model.GenerateOrderNo("DF"),
		MerchantOrderNo:      req.OrderNo,
		OrderAmount:          req.OrderAmount,
		BalanceType:          constants.DF_BALANCE,
		Type:                 constants.ORDER_TYPE_DF,
		Status:               constants.WAIT_PROCESS,
		ChannelCode:          rate.ChannelCode,
		MerchantBankAccount:  req.BankNo,
		MerchantBankNo:       req.BankId,
		MerchantBankName:     req.BankName,
		Fee:                  rate.Fee,
		HandlingFee:          rate.HandlingFee,
		MerchantAccountName:  req.DefrayName,
		CurrencyCode:         req.Currency,
		MerchantBankProvince: req.BankProvince,
		MerchantBankCity:     req.BankCity,
		Source:               constants.API,
		ChannelPayTypesCode:  rate.ChannelPayTypesCode,
		PayTypeCode:          rate.PayTypeCode,
		PayTypeCodeNum:       rate.PayTypeCodeNum,
		CreatedBy:            req.MerchantId,
		MerchantCode:         req.MerchantId,
		IsLock:               constants.IS_LOCK_NO, //是否锁定状态 (0=否;1=是) 预设否,\
		PersonProcessStatus:  constants.PERSON_PROCESS_STATUS_NO_ROCESSING,
		IsCalculateProfit:    constants.IS_CALCULATE_PROFIT_NO,
		RepaymentStatus:      constants.REPAYMENT_NOT,
		//API 要填的参数
		NotifyUrl:          req.NotifyUrl,
		IsMerchantCallback: isMerchantCallback,
	}

	// 新增收支记录，与更新商户余额(商户账户号是黑名单，把交易金额为设为 0)
	updateBalance := &types.UpdateBalance{
		MerchantCode:    txOrder.MerchantCode,
		CurrencyCode:    txOrder.CurrencyCode,
		OrderNo:         txOrder.OrderNo,
		MerchantOrderNo: txOrder.MerchantOrderNo,
		OrderType:       txOrder.Type,
		PayTypeCode:     txOrder.PayTypeCode,
		PayTypeCodeNum:  txOrder.PayTypeCodeNum,
		TransferAmount:  txOrder.TransferAmount,
		TransactionType: "11", //異動類型 (1=收款; 2=解凍; 3=沖正; 11=出款 ; 12=凍結)
		BalanceType:     constants.DF_BALANCE,
		CreatedBy:       txOrder.MerchantCode,
	}

	//判断是否是银行账号是否是黑名单
	//是。1. 失败单 2. 手续费、费率设为0 3.不在txOrder计算利润 4.交易金额设为0 更动钱包
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

	if err = l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

		//交易金额 = 订单金额 + 商户手续费
		txOrder.TransferAmount = utils.FloatAdd(txOrder.OrderAmount, txOrder.TransferHandlingFee)
		updateBalance.TransferAmount = txOrder.TransferAmount //扣款依然傳正值
		//更新钱包且新增商户钱包异动记录
		if merchantBalanceRecord, err = merchantbalanceservice.UpdateDFBalance_Debit(db, updateBalance); err != nil {
			logx.Errorf("商户:%s，更新錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err.Error(), updateBalance)
			return errorz.New(response.SYSTEM_ERROR, err.Error())
		} else {
			logx.Infof("代付API提单 %s，錢包扣款成功", merchantBalanceRecord.OrderNo)
			txOrder.BeforeBalance = merchantBalanceRecord.BeforeBalance // 商戶錢包異動紀錄
			txOrder.Balance = merchantBalanceRecord.AfterBalance
		}

		// 创建订单
		if err = db.Table("tx_orders").Create(&types.OrderX{
			Order: *txOrder,
		}).Error; err != nil {
			logx.Errorf("新增代付API提单失败，商户号: %s, 订单号: %s, err : %s", txOrder.MerchantCode, txOrder.OrderNo, err.Error())
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
	} else {
		txOrder.IsCalculateProfit = constants.IS_CALCULATE_PROFIT_YES
	}

	// 新單新增訂單歷程 (不抱錯) TODO: 異步??
	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     txOrder.OrderNo,
			Action:      "PLACE_ORDER",
			UserAccount: req.MerchantId,
			Comment:     "",
		},
	}).Error; err4 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err4.Error())
	}

	proxyOrderResp := &transactionclient.ProxyOrderResponse{
		ProxyOrderNo: txOrder.OrderNo,
	}

	return proxyOrderResp, nil
}
