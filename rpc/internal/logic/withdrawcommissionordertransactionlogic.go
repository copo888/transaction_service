package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"github.com/jinzhu/copier"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/transaction"

	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawCommissionOrderTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWithdrawCommissionOrderTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawCommissionOrderTransactionLogic {
	return &WithdrawCommissionOrderTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WithdrawCommissionOrderTransactionLogic) WithdrawCommissionOrderTransaction(in *transaction.WithdrawCommissionOrderRequest) (resp *transaction.WithdrawCommissionOrderResponse, err error) {

	/****     交易開始      ****/
	txDB := l.svcCtx.MyDB.Begin()


	var order types.CommissionWithdrawOrder
	copier.Copy(&order, &in)
	newOrderNo := model.GenerateOrderNo("YJ")

	payAt, errJ := types.JsonTime{}.Parse(in.PayAt)
	if errJ != nil {
		txDB.Rollback()
		return &transactionclient.WithdrawCommissionOrderResponse{
			Code:    response.SYSTEM_ERROR,
			Message: "新增訂單失敗",
		}, nil
	}

	var merchantCommissionRecord types.MerchantCommissionRecord
	if merchantCommissionRecord, err = merchantbalanceservice.UpdateCommissionAmount(txDB, types.UpdateCommissionAmount{
		MerchantCode:            in.MerchantCode,
		CurrencyCode:            in.WithdrawCurrencyCode,
		OrderNo:                 newOrderNo,
		TransactionType:         constants.COMMISSION_TRANSACTION_TYPE_WITHDRAWAL,
		TransferAmount:          -in.WithdrawAmount,
		Comment:                 "",
		CreatedBy:               in.CreatedBy,
	}); err != nil {
		txDB.Rollback()
		return &transactionclient.WithdrawCommissionOrderResponse{
			Code:    response.MERCHANT_WALLET_RECORD_ERROR,
			Message: "更新錢包失敗",
		}, nil
	} else if merchantCommissionRecord.AfterCommission < 0 {
		txDB.Rollback()
		return &transactionclient.WithdrawCommissionOrderResponse{
			Code:    response.MERCHANT_COMMISSION_WALLET_NO_ENOUGH,
			Message: "佣金錢包餘額不足",
		}, nil
	}

	order.OrderNo = newOrderNo
	order.AfterCommission = merchantCommissionRecord.AfterCommission
	if err = txDB.Table("cm_withdraw_order").Create(&types.CommissionWithdrawOrderX{
		CommissionWithdrawOrder:   order,
		PayAt: payAt,
	}).Error; err != nil {
		txDB.Rollback()
		return &transactionclient.WithdrawCommissionOrderResponse{
			Code:    response.SYSTEM_ERROR,
			Message: "新增訂單失敗",
		}, nil
	}

	if err = txDB.Commit().Error; err != nil {
		txDB.Rollback()
		logx.Errorf("Commit失败，商户号: %s, 订单号: %s, err : %s", order.MerchantCode, order.OrderNo, err.Error())
		return &transactionclient.WithdrawCommissionOrderResponse{
			Code:    response.DATABASE_FAILURE,
			Message: "资料库错误 Commit失败",
		}, nil
	}
	/****     交易結束      ****/

	return &transactionclient.WithdrawCommissionOrderResponse{
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}, nil

}
