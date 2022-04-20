package logic

import (
	"context"
	"fmt"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transaction"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmPayOrderTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmPayOrderTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmPayOrderTransactionLogic {
	return &ConfirmPayOrderTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfirmPayOrderTransactionLogic) ConfirmPayOrderTransaction(in *transaction.ConfirmPayOrderRequest) (*transaction.ConfirmPayOrderResponse, error) {
	var order *types.OrderX

	if err := l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

		if err = l.svcCtx.MyDB.Table("tx_orders").
			Where("order_no = ?", in.OrderNo).Take(&order).Error; err != nil {
			return errorz.New(response.ORDER_NUMBER_NOT_EXIST, err.Error())
		}

		// 處理中的且非鎖定訂單 才能確認收款
		if order.Status != "1" || order.IsLock == "1" {
			return errorz.New(response.TRANSACTION_FAILURE, fmt.Sprintf("(status/isLock): %s/%s", order.Status, order.IsLock))
		}

		// 編輯訂單 異動錢包和餘額
		if err = l.updateOrderAndBalance(db, in, order); err != nil {
			return errorz.New(response.SYSTEM_ERROR, err.Error())
		}

		return
	}); err != nil {
		return nil, err
	}

	return &transaction.ConfirmPayOrderResponse{}, nil
}

func (l *ConfirmPayOrderTransactionLogic) updateOrderAndBalance(db *gorm.DB, req *transactionclient.ConfirmPayOrderRequest, order *types.OrderX) (err error) {

	var merchantBalanceRecord types.MerchantBalanceRecord



	// 異動錢包
	if merchantBalanceRecord, err = merchantbalanceservice.UpdateBalanceForZF(db, types.UpdateBalance{
		MerchantCode:    order.MerchantCode,
		CurrencyCode:    order.CurrencyCode,
		OrderNo:         order.OrderNo,
		OrderType:       order.Type,
		ChannelCode:     order.ChannelCode,
		PayTypeCode:     order.PayTypeCode,
		PayTypeCodeNum:  order.PayTypeCodeNum,
		TransactionType: "1",
		BalanceType:     order.BalanceType,
		TransferAmount:  order.TransferAmount,
		Comment:         order.Memo,
		CreatedBy:       order.MerchantCode,
	}); err != nil {
		return
	}

	// 手動確認收款 實際金額 = 訂單金額
	order.ActualAmount = order.OrderAmount
	order.BeforeBalance = merchantBalanceRecord.BeforeBalance
	order.Balance = merchantBalanceRecord.AfterBalance
	order.TransAt = types.JsonTime{}.New()
	order.Status = "20"
	order.Memo = req.Comment

	// 編輯訂單
	if err = db.Table("tx_orders").Updates(&order).Error; err != nil {
		return
	}

	return
}
