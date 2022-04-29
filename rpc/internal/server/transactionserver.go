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

func (s *TransactionServer) MerchantBalanceUpdateTranaction(ctx context.Context, in *transactionclient.MerchantBalanceUpdateRequest) (*transactionclient.MerchantBalanceUpdateResponse, error) {
	l := logic.NewMerchantBalanceUpdateTranactionLogic(ctx, s.svcCtx)
	return l.MerchantBalanceUpdateTranaction(in)
}

func (s *TransactionServer) ProxyOrderTranaction_DFB(ctx context.Context, in *transactionclient.ProxyOrderRequest) (*transactionclient.ProxyOrderResponse, error) {
	l := logic.NewProxyOrderTranactionDFBLogic(ctx, s.svcCtx)
	return l.ProxyOrderTranaction_DFB(in)
}

func (s *TransactionServer) ProxyOrderTranaction_XFB(ctx context.Context, in *transactionclient.ProxyOrderRequest) (*transactionclient.ProxyOrderResponse, error) {
	l := logic.NewProxyOrderTranactionXFBLogic(ctx, s.svcCtx)
	return l.ProxyOrderTranaction_XFB(in)
}

func (s *TransactionServer) ProxyOrderTransactionFail_DFB(ctx context.Context, in *transactionclient.ProxyPayFailRequest) (*transactionclient.ProxyPayFailResponse, error) {
	l := logic.NewProxyOrderTransactionFailDFBLogic(ctx, s.svcCtx)
	return l.ProxyOrderTransactionFail_DFB(in)
}

func (s *TransactionServer) ProxyOrderTransactionFail_XFB(ctx context.Context, in *transactionclient.ProxyPayFailRequest) (*transactionclient.ProxyPayFailResponse, error) {
	l := logic.NewProxyOrderTransactionFailXFBLogic(ctx, s.svcCtx)
	return l.ProxyOrderTransactionFail_XFB(in)
}

func (s *TransactionServer) ProxyOrderToTest_DFB(ctx context.Context, in *transactionclient.ProxyOrderTestRequest) (*transactionclient.ProxyOrderTestResponse, error) {
	l := logic.NewProxyOrderToTestDFBLogic(ctx, s.svcCtx)
	return l.ProxyOrderToTest_DFB(in)
}

func (s *TransactionServer) ProxyOrderToTest_XFB(ctx context.Context, in *transactionclient.ProxyOrderTestRequest) (*transactionclient.ProxyOrderTestResponse, error) {
	l := logic.NewProxyOrderToTestXFBLogic(ctx, s.svcCtx)
	return l.ProxyOrderToTest_XFB(in)
}

func (s *TransactionServer) PayOrderTranaction(ctx context.Context, in *transactionclient.PayOrderRequest) (*transactionclient.PayOrderResponse, error) {
	l := logic.NewPayOrderTranactionLogic(ctx, s.svcCtx)
	return l.PayOrderTranaction(in)
}

func (s *TransactionServer) InternalOrderTransaction(ctx context.Context, in *transactionclient.InternalOrderRequest) (*transactionclient.InternalOrderResponse, error) {
	l := logic.NewInternalOrderTransactionLogic(ctx, s.svcCtx)
	return l.InternalOrderTransaction(in)
}

func (s *TransactionServer) WithdrawOrderTransaction(ctx context.Context, in *transactionclient.WithdrawOrderRequest) (*transactionclient.WithdrawOrderResponse, error) {
	l := logic.NewWithdrawOrderTransactionLogic(ctx, s.svcCtx)
	return l.WithdrawOrderTransaction(in)
}

func (s *TransactionServer) PayCallBackTranaction(ctx context.Context, in *transactionclient.PayCallBackRequest) (*transactionclient.PayCallBackResponse, error) {
	l := logic.NewPayCallBackTranactionLogic(ctx, s.svcCtx)
	return l.PayCallBackTranaction(in)
}

func (s *TransactionServer) InternalReviewSuccessTransaction(ctx context.Context, in *transactionclient.InternalReviewSuccessRequest) (*transactionclient.InternalReviewSuccessResponse, error) {
	l := logic.NewInternalReviewSuccessTransactionLogic(ctx, s.svcCtx)
	return l.InternalReviewSuccessTransaction(in)
}

func (s *TransactionServer) WithdrawReviewFailTransaction(ctx context.Context, in *transactionclient.WithdrawReviewFailRequest) (*transactionclient.WithdrawReviewFailResponse, error) {
	l := logic.NewWithdrawReviewFailTransactionLogic(ctx, s.svcCtx)
	return l.WithdrawReviewFailTransaction(in)
}

func (s *TransactionServer) WithdrawReviewSuccessTransaction(ctx context.Context, in *transactionclient.WithdrawReviewSuccessRequest) (*transactionclient.WithdrawReviewSuccessResponse, error) {
	l := logic.NewWithdrawReviewSuccessTransactionLogic(ctx, s.svcCtx)
	return l.WithdrawReviewSuccessTransaction(in)
}

func (s *TransactionServer) ProxyOrderUITransaction_DFB(ctx context.Context, in *transactionclient.ProxyOrderUIRequest) (*transactionclient.ProxyOrderUIResponse, error) {
	l := logic.NewProxyOrderUITransactionDFBLogic(ctx, s.svcCtx)
	return l.ProxyOrderUITransaction_DFB(in)
}

func (s *TransactionServer) ProxyOrderUITransaction_XFB(ctx context.Context, in *transactionclient.ProxyOrderUIRequest) (*transactionclient.ProxyOrderUIResponse, error) {
	l := logic.NewProxyOrderUITransactionXFBLogic(ctx, s.svcCtx)
	return l.ProxyOrderUITransaction_XFB(in)
}

func (s *TransactionServer) MakeUpReceiptOrderTransaction(ctx context.Context, in *transactionclient.MakeUpReceiptOrderRequest) (*transactionclient.MakeUpReceiptOrderResponse, error) {
	l := logic.NewMakeUpReceiptOrderTransactionLogic(ctx, s.svcCtx)
	return l.MakeUpReceiptOrderTransaction(in)
}

func (s *TransactionServer) ConfirmPayOrderTransaction(ctx context.Context, in *transactionclient.ConfirmPayOrderRequest) (*transactionclient.ConfirmPayOrderResponse, error) {
	l := logic.NewConfirmPayOrderTransactionLogic(ctx, s.svcCtx)
	return l.ConfirmPayOrderTransaction(in)
}

func (s *TransactionServer) RecoverReceiptOrderTransaction(ctx context.Context, in *transactionclient.RecoverReceiptOrderRequest) (*transactionclient.RecoverReceiptOrderResponse, error) {
	l := logic.NewRecoverReceiptOrderTransactionLogic(ctx, s.svcCtx)
	return l.RecoverReceiptOrderTransaction(in)
}

func (s *TransactionServer) FrozenReceiptOrderTransaction(ctx context.Context, in *transactionclient.FrozenReceiptOrderRequest) (*transactionclient.FrozenReceiptOrderResponse, error) {
	l := logic.NewFrozenReceiptOrderTransactionLogic(ctx, s.svcCtx)
	return l.FrozenReceiptOrderTransaction(in)
}

func (s *TransactionServer) UnFrozenReceiptOrderTransaction(ctx context.Context, in *transactionclient.UnFrozenReceiptOrderRequest) (*transactionclient.UnFrozenReceiptOrderResponse, error) {
	l := logic.NewUnFrozenReceiptOrderTransactionLogic(ctx, s.svcCtx)
	return l.UnFrozenReceiptOrderTransaction(in)
}

func (s *TransactionServer) PersonalRebundTransaction_DFB(ctx context.Context, in *transactionclient.PersonalRebundRequest) (*transactionclient.PersonalRebundResponse, error) {
	l := logic.NewPersonalRebundTransactionDFBLogic(ctx, s.svcCtx)
	return l.PersonalRebundTransaction_DFB(in)
}

func (s *TransactionServer) PersonalRebundTransaction_XFB(ctx context.Context, in *transactionclient.PersonalRebundRequest) (*transactionclient.PersonalRebundResponse, error) {
	l := logic.NewPersonalRebundTransactionXFBLogic(ctx, s.svcCtx)
	return l.PersonalRebundTransaction_XFB(in)
}
