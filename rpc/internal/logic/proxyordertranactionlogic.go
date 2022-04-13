package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/service/orderfeeprofitservice"
	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProxyOrderTranactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProxyOrderTranactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyOrderTranactionLogic {
	return &ProxyOrderTranactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

/*
	@PARAM ProxyOrderReqest
    @param1 req 	    商戶代付提單請求參數
    @param2 rate 	    商戶配置費率、手續費
	@param3 balanceType 商戶錢包類型

	@return
*/
func (l *ProxyOrderTranactionLogic) ProxyOrderTranaction(in *transactionclient.ProxyOrderReq_DFB) (*transactionclient.ProxyOrderResp_DFB, error) {
	userAccount := "TEST0001"

	tx := l.svcCtx.MyDB
	req := in.Req
	rate := in.Rate
	balanceType := in.BalanceType

	// 依商户是否给回调网址，决定是否回调商户flag
	var isMerchantCallback string
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
		CreatedBy:            userAccount,
		MerchantCode:         req.MerchantId,
		IsLock:               "0", //是否锁定状态 (0=否;1=是) 预设否,\
		//API 要填的参数
		NotifyUrl:          req.NotifyUrl,
		IsMerchantCallback: isMerchantCallback,
	}

	// 新增收支记录，与更新商户余额(商户账户号是黑名单，把交易金额为设为 0)
	updateBalance := types.UpdateBalance{
		MerchantCode:    txOrder.MerchantCode,
		CurrencyCode:    txOrder.CurrencyCode,
		OrderNo:         txOrder.OrderNo,
		OrderType:       txOrder.Type,
		PayTypeCode:     txOrder.PayTypeCode,
		PayTypeCodeNum:  txOrder.PayTypeCodeNum,
		TransferAmount:  txOrder.TransferAmount,
		TransactionType: "11", //異動類型 (1=收款; 2=解凍; 3=沖正; 11=出款 ; 12=凍結)
		BalanceType:     balanceType,
		CreatedBy:       userAccount,
	}

	var transAt types.JsonTime

	tx.Begin()
	//判断是否是银行账号是否是黑名单
	//是。1. 失败单 2. 手续费、费率设为0 3.不在txOrder计算利润 4.交易金额设为0 更动钱包
	isBlock, _ := model.NewBankBlockAccount(tx).CheckIsBlockAccount(txOrder.MerchantBankAccount)
	if isBlock { //银行账号为黑名单
		logx.Infof("交易账户%s-%s在黑名单内，使用0元假扣款", txOrder.MerchantAccountName, txOrder.MerchantBankNo)
		updateBalance.TransferAmount = 0                           // 使用0元前往钱包扣款
		txOrder.ErrorType = constants.ERROR6_BANK_ACCOUNT_IS_BLACK //交易账户为黑名单
		txOrder.ErrorNote = constants.BANK_ACCOUNT_IS_BLACK        //失败原因：黑名单交易失败
		txOrder.Status = constants.PROXY_PAY_FAIL                  //状态:失败
		txOrder.Fee = 0                                            //写入本次手续费(未发送到渠道的交易，都设为0元)
		txOrder.HandlingFee = 0
		//txOrder.TransAt     = time.Now().Format("20060102150405") //失败单,一并写入交易日期
		transAt = types.JsonTime{}.New()
		logx.Infof("商户 %s，代付订单 %#v ，交易账户为黑名单", txOrder.MerchantCode, txOrder)
	} else { //不为黑名单

		//交易金额 = 订单金额 + 商户手续费
		txOrder.TransferAmount = utils.FloatAdd(txOrder.OrderAmount, txOrder.TransferHandlingFee)
		updateBalance.TransferAmount = txOrder.TransferAmount //扣款依然傳正值

		//更新钱包且新增商户钱包异动记录
		merchantBalanceRecord, err1 := merchantbalanceservice.UpdateDFBalance_Debit(tx, updateBalance)
		if err1 != nil {
			logx.Errorf("商户:%s，更新錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err1.Error(), updateBalance)
			//TODO  IF 更新钱包错误是response.DATABASE_FAILURE THEN return SYSTEM_ERROR
			tx.Rollback()
			return nil, err1
		} else {
			logx.Infof("代付API提单 %s，錢包扣款成功", merchantBalanceRecord.OrderNo)
			txOrder.BalanceType = balanceType                           // 餘額類型(下发余额、代付余额)
			txOrder.BeforeBalance = merchantBalanceRecord.BeforeBalance // 商戶錢包異動紀錄
			txOrder.Balance = merchantBalanceRecord.AfterBalance
		}
	}

	// 创建订单
	if err3 := tx.Table("tx_orders").Create(&types.OrderX{
		Order:   *txOrder,
		TransAt: transAt}).Error; err3 != nil {
		logx.Errorf("新增代付API提单失败，商户号: %s, 订单号: %s, err : %s", txOrder.MerchantCode, txOrder.OrderNo, err3.Error())
		tx.Rollback()
		return nil, err3
	}
	tx.Commit()

	// 計算利潤 TODO: 異步??
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

	proxyOrderResp := &transactionclient.ProxyOrderResp_DFB{
		ProxyOrderNo: txOrder.OrderNo,
	}

	return proxyOrderResp, nil
}
