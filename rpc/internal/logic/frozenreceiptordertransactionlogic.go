package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/transaction"

	"github.com/zeromicro/go-zero/core/logx"
)

type FrozenReceiptOrderTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFrozenReceiptOrderTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FrozenReceiptOrderTransactionLogic {
	return &FrozenReceiptOrderTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FrozenReceiptOrderTransactionLogic) FrozenReceiptOrderTransaction(req *transaction.FrozenReceiptOrderRequest) (*transaction.FrozenReceiptOrderResponse, error) {
	var order types.Order

	/****     交易開始      ****/
	txDB := l.svcCtx.MyDB.Begin()

	// 1. 取得訂單
	if err := txDB.Table("tx_orders").Where("order_no = ?", req.OrderNo).Find(&order).Error; err != nil {
		txDB.Rollback()
		return &transactionclient.FrozenReceiptOrderResponse{
			Code: response.DATABASE_FAILURE,
			Message: "取得訂單失敗",
		}, nil
	}

	// 驗證
	if errCode := l.verify(order, req); errCode != "" {
		txDB.Rollback()
		return &transactionclient.FrozenReceiptOrderResponse{
			Code: errCode,
			Message: "驗證失敗: " + errCode,
		}, nil
	}

	// 變更 商戶餘額&凍結金額 並記錄
	if _, err := merchantbalanceservice.UpdateFrozenAmount(txDB, types.UpdateFrozenAmount{
		MerchantCode:    order.MerchantCode,
		CurrencyCode:    order.CurrencyCode,
		OrderNo:         order.OrderNo,
		OrderType:       order.Type,
		ChannelCode:     order.ChannelCode,
		PayTypeCode:     order.PayTypeCode,
		PayTypeCodeNum:  order.PayTypeCodeNum,
		TransactionType: constants.TRANSACTION_TYPE_FREEZE,
		BalanceType:     order.BalanceType,
		FrozenAmount:    req.FrozenAmount,
		Comment:         req.Comment,
		CreatedBy:       "AAA00061", // TODO: JWT取得
	}); err != nil {
		txDB.Rollback()
		return &transactionclient.FrozenReceiptOrderResponse{
			Code: response.DATABASE_FAILURE,
			Message: "錢包異動失敗",
		}, nil
	}

	// 變更訂單狀態
	order.Status = constants.FROZEN
	order.FrozenAmount = req.FrozenAmount
	order.Memo = req.Comment + " \n" + order.Memo

	if err := txDB.Table("tx_orders").Updates(&types.OrderX{
		Order: order,
		FrozenAt: types.JsonTime{}.New(),
	}).Error; err != nil {
		txDB.Rollback()
		return &transactionclient.FrozenReceiptOrderResponse{
			Code: response.DATABASE_FAILURE,
			Message: "編輯訂單失敗",
		}, nil
	}

	if err := txDB.Commit().Error; err != nil {
		txDB.Rollback()
		logx.Errorf("凍結收款膽Commit失败，商户号: %s, 订单号: %s, err : %s", order.MerchantCode, order.OrderNo, err.Error())
		return &transactionclient.FrozenReceiptOrderResponse{
			Code:    response.DATABASE_FAILURE,
			Message: "资料库错误 Commit失败",
		}, nil
	}
	/****     交易結束      ****/


	// 新單新增訂單歷程 (不抱錯)
	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     order.OrderNo,
			Action:      constants.ACTION_FROZEN,
			UserAccount: order.MerchantCode,
			Comment:     req.Comment,
		},
	}).Error; err4 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err4.Error())
	}


	return &transaction.FrozenReceiptOrderResponse{
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}, nil
}

func (l *FrozenReceiptOrderTransactionLogic) verify(order types.Order, req *transactionclient.FrozenReceiptOrderRequest) string {

	// 收款單才能凍結
	if order.Type != constants.ORDER_TYPE_ZF && order.Type != constants.ORDER_TYPE_NC  {
		return response.ORDER_STATUS_WRONG_CANNOT_FROZEN
	}

	// 成功單才能凍結
	if order.Status != constants.SUCCESS {
		return response.ORDER_STATUS_WRONG_CANNOT_FROZEN
	}

	// 凍結金額需大於交易金額
	if req.FrozenAmount < order.TransferAmount {
		return response.FROZEN_AMOUNT_NOT_LESS_THAN_ORDER_AMOUNT
	}

	return ""
}
