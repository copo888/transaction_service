package logic

import (
	"context"

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

func (l *WithdrawCommissionOrderTransactionLogic) WithdrawCommissionOrderTransaction(in *transaction.WithdrawCommissionOrderRequest) (*transaction.WithdrawCommissionOrderResponse, error) {
	// todo: add your logic here and delete this line

	return &transaction.WithdrawCommissionOrderResponse{}, nil
}
