package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"

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
	tx := l.svcCtx.MyDB

	transferAmount := utils.FloatAdd(in.OrderAmount, in.HandlingFee)

	// 初始化订单
	txOrder := &types.Order{
		MerchantCode: in.MerchantCode,
		MerchantOrderNo: "COPO_"+in.OrderNo,
		OrderNo: in.OrderNo,
		Type: constants.ORDER_TYPE_XF,
		Status: constants.WAIT_PROCESS,
		CreatedBy: in.UserAccount,
		UpdatedBy: in.UserAccount,
		BalanceType: "XFB",
		ChannelCode: in.CurrencyCode,
		TransferAmount: transferAmount,
		OrderAmount: in.OrderAmount,
		TransferHandlingFee: in.HandlingFee,
		HandlingFee: in.HandlingFee,
		MerchantBankAccount: in.MerchantBankeAccount,
		MerchantBankNo: in.MerchantBankNo,
		MerchantBankProvince: in.MerchantBankProvince,
		MerchantBankCity: in.MerchantBankCity,
		MerchantAccountName: in.MerchantAccountName,
		IsMerchantCallback: constants.IS_MERCHANT_CALLBACK_NO,
		IsLock: constants.IS_LOCK_NO,
		CurrencyCode: in.CurrencyCode,
		Source: in.Source,
		PageUrl: in.PageUrl,
		NotifyUrl: in.NotifyUrl,
	}
	if len(in.MerchantOrderNo) > 0 {
		txOrder.MerchantOrderNo = in.MerchantOrderNo
	}

	tx.Begin()
	// 新增收支记录，更新商户余额
	updateBalance := types.UpdateBalance{
		MerchantCode:    txOrder.MerchantCode,
		CurrencyCode:    txOrder.CurrencyCode,
		OrderNo:         txOrder.OrderNo,
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
			txOrder.Status = constants.PROXY_PAY_FAIL                  //状态:失败
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
		Order:   *txOrder,
		}).Error; err3 != nil {
		logx.Errorf("新增下发提单失败，商户号: %s, 订单号: %s, err : %s", txOrder.MerchantCode, txOrder.OrderNo, err3.Error())
		tx.Rollback()
		return nil, err3
	}
	if err4 := tx.Commit().Error; err4 != nil {
		logx.Errorf("最终新增下发提单失败，商户号: %s, 订单号: %s, err : %s", txOrder.MerchantCode, txOrder.OrderNo, err4.Error())
		tx.Rollback()
		return nil, err4
	}

	// 新單新增訂單歷程 (不抱錯) TODO: 異步??
	if err5 := tx.Table("tx_order_actions").Create(&types.OrderActionX{
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
