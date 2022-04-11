package logic

import (
	"context"
	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/transaction"
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

func (l *ProxyOrderTranactionLogic) ProxyOrderTranaction(in *transaction.ProxyOrderReq_DFB) (*transaction.ProxyOrderResp_DFB, error) {
	// todo: add your logic here and delete this line

	return &transaction.ProxyOrderResp_DFB{}, nil
}
