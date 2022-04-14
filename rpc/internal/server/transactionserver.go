// Code generated by goctl. DO NOT EDIT!
// Source: transaction.proto

package server

import (
	"context"
	"github.com/copo888/transaction_service/rpc/internal/logic"
	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/transaction"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type TransactionServer struct {
	svcCtx *svc.ServiceContext
	transaction.UnimplementedTransactionServer
}

func NewTransactionServer(svcCtx *svc.ServiceContext) *TransactionServer {
	return &TransactionServer{
		svcCtx: svcCtx,
	}
}

func (s *TransactionServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *TransactionServer) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return nil
}

func (s *TransactionServer) ProxyOrderTranaction_DFB(ctx context.Context, in *transactionclient.ProxyOrderRequest) (*transactionclient.ProxyOrderResponse, error) {
	l := logic.NewProxyOrderTranactionDFBLogic(ctx, s.svcCtx)
	return l.ProxyOrderTranaction_DFB(in)
}

func (s *TransactionServer) ProxyOrderTranaction_XFB(ctx context.Context, in *transactionclient.ProxyOrderRequest) (*transactionclient.ProxyOrderResponse, error) {
	l := logic.NewProxyOrderTranactionXFBLogic(ctx, s.svcCtx)
	return l.ProxyOrderTranaction_XFB(in)
}

func (s *TransactionServer) PayOrderTranaction(ctx context.Context, in *transactionclient.PayOrderRequest) (*transactionclient.PayOrderResponse, error) {
	l := logic.NewPayOrderTranactionLogic(ctx, s.svcCtx)
	return l.PayOrderTranaction(in)
}

func (s *TransactionServer) PayCallBackTranaction(ctx context.Context, in *transactionclient.PayCallBackRequest) (*transactionclient.PayCallBackResponse, error) {
	l := logic.NewPayCallBackTranactionLogic(ctx, s.svcCtx)
	return l.PayCallBackTranaction(in)
}
