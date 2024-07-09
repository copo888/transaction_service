// Code generated by goctl. DO NOT EDIT!
// Source: transaction.proto

package server

import (
	"context"
	"github.com/copo888/transaction_service/rpc/internal/logic"
	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/transaction"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "MerchantBalanceUpdateTranaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.MerchantBalanceUpdateTranaction(in)
}

func (s *TransactionServer) MerchantBalanceFreezeTranaction(ctx context.Context, in *transactionclient.MerchantBalanceFreezeRequest) (*transactionclient.MerchantBalanceFreezeResponse, error) {
	l := logic.NewMerchantBalanceFreezeTranactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "MerchantBalanceFreezeTranaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.MerchantBalanceFreezeTranaction(in)
}

func (s *TransactionServer) ProxyOrderTranaction_DFB(ctx context.Context, in *transactionclient.ProxyOrderRequest) (*transactionclient.ProxyOrderResponse, error) {
	l := logic.NewProxyOrderTranactionDFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "ProxyOrderTranaction_DFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.ProxyOrderTranaction_DFB(in)
}

func (s *TransactionServer) ProxyOrderTranaction_XFB(ctx context.Context, in *transactionclient.ProxyOrderRequest) (*transactionclient.ProxyOrderResponse, error) {
	l := logic.NewProxyOrderTranactionXFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "ProxyOrderTranaction_XFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.ProxyOrderTranaction_XFB(in)
}

func (s *TransactionServer) ProxyOrderSmartTranaction_DFB(ctx context.Context, in *transactionclient.ProxyOrderSmartRequest) (*transactionclient.ProxyOrderResponse, error) {
	l := logic.NewProxyOrderSmartTranactionDFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "ProxyOrderSmartTranaction_DFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.ProxyOrderSmartTranaction_DFB(in)
}

func (s *TransactionServer) ProxyOrderSmartTranaction_XFB(ctx context.Context, in *transactionclient.ProxyOrderSmartRequest) (*transactionclient.ProxyOrderResponse, error) {
	l := logic.NewProxyOrderSmartTranactionXFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "ProxyOrderSmartTranaction_XFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.ProxyOrderSmartTranaction_XFB(in)
}

func (s *TransactionServer) ProxyOrderTransactionFail_DFB(ctx context.Context, in *transactionclient.ProxyPayFailRequest) (*transactionclient.ProxyPayFailResponse, error) {
	l := logic.NewProxyOrderTransactionFailDFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "ProxyOrderTransactionFail_DFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.ProxyOrderTransactionFail_DFB(in)
}

func (s *TransactionServer) ProxyOrderTransactionFail_XFB(ctx context.Context, in *transactionclient.ProxyPayFailRequest) (*transactionclient.ProxyPayFailResponse, error) {
	l := logic.NewProxyOrderTransactionFailXFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "ProxyOrderTransactionFail_XFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.ProxyOrderTransactionFail_XFB(in)
}

func (s *TransactionServer) PayOrderSwitchTest(ctx context.Context, in *transactionclient.PayOrderSwitchTestRequest) (*transactionclient.PayOrderSwitchTestResponse, error) {
	l := logic.NewPayOrderSwitchTestLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "PayOrderSwitchTest",
		Value: attribute.StringValue(in.String()),
	})
	return l.PayOrderSwitchTest(in)
}

func (s *TransactionServer) ProxyOrderToTest_DFB(ctx context.Context, in *transactionclient.ProxyOrderTestRequest) (*transactionclient.ProxyOrderTestResponse, error) {
	l := logic.NewProxyOrderToTestDFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "ProxyOrderToTest_DFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.ProxyOrderToTest_DFB(in)
}

func (s *TransactionServer) ProxyOrderToTest_XFB(ctx context.Context, in *transactionclient.ProxyOrderTestRequest) (*transactionclient.ProxyOrderTestResponse, error) {
	l := logic.NewProxyOrderToTestXFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "ProxyOrderToTest_XFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.ProxyOrderToTest_XFB(in)
}

func (s *TransactionServer) ProxyTestToNormal_DFB(ctx context.Context, in *transactionclient.ProxyOrderTestRequest) (*transactionclient.ProxyOrderTestResponse, error) {
	l := logic.NewProxyTestToNormalDFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "ProxyTestToNormal_DFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.ProxyTestToNormal_DFB(in)
}

func (s *TransactionServer) ProxyTestToNormal_XFB(ctx context.Context, in *transactionclient.ProxyOrderTestRequest) (*transactionclient.ProxyOrderTestResponse, error) {
	l := logic.NewProxyTestToNormalXFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "ProxyTestToNormal_XFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.ProxyTestToNormal_XFB(in)
}

func (s *TransactionServer) PayOrderTranaction(ctx context.Context, in *transactionclient.PayOrderRequest) (*transactionclient.PayOrderResponse, error) {
	l := logic.NewPayOrderTranactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "PayOrderTranaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.PayOrderTranaction(in)
}

func (s *TransactionServer) InternalOrderTransaction(ctx context.Context, in *transactionclient.InternalOrderRequest) (*transactionclient.InternalOrderResponse, error) {
	l := logic.NewInternalOrderTransactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "InternalOrderTransaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.InternalOrderTransaction(in)
}

func (s *TransactionServer) WithdrawOrderTransaction(ctx context.Context, in *transactionclient.WithdrawOrderRequest) (*transactionclient.WithdrawOrderResponse, error) {
	l := logic.NewWithdrawOrderTransactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "WithdrawOrderTransaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.WithdrawOrderTransaction(in)
}

func (s *TransactionServer) PayCallBackTranaction(ctx context.Context, in *transactionclient.PayCallBackRequest) (*transactionclient.PayCallBackResponse, error) {
	l := logic.NewPayCallBackTranactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "PayCallBackTranaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.PayCallBackTranaction(in)
}

func (s *TransactionServer) InternalReviewSuccessTransaction(ctx context.Context, in *transactionclient.InternalReviewSuccessRequest) (*transactionclient.InternalReviewSuccessResponse, error) {
	l := logic.NewInternalReviewSuccessTransactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "InternalReviewSuccessTransaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.InternalReviewSuccessTransaction(in)
}

func (s *TransactionServer) WithdrawReviewFailTransaction(ctx context.Context, in *transactionclient.WithdrawReviewFailRequest) (*transactionclient.WithdrawReviewFailResponse, error) {
	l := logic.NewWithdrawReviewFailTransactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "WithdrawReviewFailTransaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.WithdrawReviewFailTransaction(in)
}

func (s *TransactionServer) WithdrawReviewSuccessTransaction(ctx context.Context, in *transactionclient.WithdrawReviewSuccessRequest) (*transactionclient.WithdrawReviewSuccessResponse, error) {
	l := logic.NewWithdrawReviewSuccessTransactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "WithdrawReviewSuccessTransaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.WithdrawReviewSuccessTransaction(in)
}

func (s *TransactionServer) ProxyOrderUITransaction_DFB(ctx context.Context, in *transactionclient.ProxyOrderUIRequest) (*transactionclient.ProxyOrderUIResponse, error) {
	l := logic.NewProxyOrderUITransactionDFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "ProxyOrderUITransaction_DFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.ProxyOrderUITransaction_DFB(in)
}

func (s *TransactionServer) ProxyOrderUITransaction_XFB(ctx context.Context, in *transactionclient.ProxyOrderUIRequest) (*transactionclient.ProxyOrderUIResponse, error) {
	l := logic.NewProxyOrderUITransactionXFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "ProxyOrderUITransaction_XFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.ProxyOrderUITransaction_XFB(in)
}

func (s *TransactionServer) MakeUpReceiptOrderTransaction(ctx context.Context, in *transactionclient.MakeUpReceiptOrderRequest) (*transactionclient.MakeUpReceiptOrderResponse, error) {
	l := logic.NewMakeUpReceiptOrderTransactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "MakeUpReceiptOrderTransaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.MakeUpReceiptOrderTransaction(in)
}

func (s *TransactionServer) ConfirmPayOrderTransaction(ctx context.Context, in *transactionclient.ConfirmPayOrderRequest) (*transactionclient.ConfirmPayOrderResponse, error) {
	l := logic.NewConfirmPayOrderTransactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "ConfirmPayOrderTransaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.ConfirmPayOrderTransaction(in)
}

func (s *TransactionServer) ConfirmProxyPayOrderTransaction(ctx context.Context, in *transactionclient.ConfirmProxyPayOrderRequest) (*transactionclient.ConfirmProxyPayOrderResponse, error) {
	l := logic.NewConfirmProxyPayOrderTransactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "ConfirmProxyPayOrderTransaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.ConfirmProxyPayOrderTransaction(in)
}

func (s *TransactionServer) RecoverReceiptOrderTransaction(ctx context.Context, in *transactionclient.RecoverReceiptOrderRequest) (*transactionclient.RecoverReceiptOrderResponse, error) {
	l := logic.NewRecoverReceiptOrderTransactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "RecoverReceiptOrderTransaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.RecoverReceiptOrderTransaction(in)
}

func (s *TransactionServer) FrozenReceiptOrderTransaction(ctx context.Context, in *transactionclient.FrozenReceiptOrderRequest) (*transactionclient.FrozenReceiptOrderResponse, error) {
	l := logic.NewFrozenReceiptOrderTransactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "FrozenReceiptOrderTransaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.FrozenReceiptOrderTransaction(in)
}

func (s *TransactionServer) UnFrozenReceiptOrderTransaction(ctx context.Context, in *transactionclient.UnFrozenReceiptOrderRequest) (*transactionclient.UnFrozenReceiptOrderResponse, error) {
	l := logic.NewUnFrozenReceiptOrderTransactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "UnFrozenReceiptOrderTransaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.UnFrozenReceiptOrderTransaction(in)
}

func (s *TransactionServer) PersonalRebundTransaction_DFB(ctx context.Context, in *transactionclient.PersonalRebundRequest) (*transactionclient.PersonalRebundResponse, error) {
	l := logic.NewPersonalRebundTransactionDFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "PersonalRebundTransaction_DFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.PersonalRebundTransaction_DFB(in)
}

func (s *TransactionServer) PersonalRebundTransaction_XFB(ctx context.Context, in *transactionclient.PersonalRebundRequest) (*transactionclient.PersonalRebundResponse, error) {
	l := logic.NewPersonalRebundTransactionXFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "PersonalRebundTransaction_XFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.PersonalRebundTransaction_XFB(in)
}

func (s *TransactionServer) RecalculateProfitTransaction(ctx context.Context, in *transactionclient.RecalculateProfitRequest) (*transactionclient.RecalculateProfitResponse, error) {
	l := logic.NewRecalculateProfitTransactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "RecalculateProfitTransaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.RecalculateProfitTransaction(in)
}

func (s *TransactionServer) CalculateCommissionMonthAllReport(ctx context.Context, in *transactionclient.CalculateCommissionMonthAllRequest) (*transactionclient.CalculateCommissionMonthAllResponse, error) {
	l := logic.NewCalculateCommissionMonthAllReportLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "CalculateCommissionMonthAllReport",
		Value: attribute.StringValue(in.String()),
	})
	return l.CalculateCommissionMonthAllReport(in)
}

func (s *TransactionServer) RecalculateCommissionMonthReport(ctx context.Context, in *transactionclient.RecalculateCommissionMonthReportRequest) (*transactionclient.RecalculateCommissionMonthReportResponse, error) {
	l := logic.NewRecalculateCommissionMonthReportLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "RecalculateCommissionMonthReport",
		Value: attribute.StringValue(in.String()),
	})
	return l.RecalculateCommissionMonthReport(in)
}

func (s *TransactionServer) ConfirmCommissionMonthReport(ctx context.Context, in *transactionclient.ConfirmCommissionMonthReportRequest) (*transactionclient.ConfirmCommissionMonthReportResponse, error) {
	l := logic.NewConfirmCommissionMonthReportLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "ConfirmCommissionMonthReport",
		Value: attribute.StringValue(in.String()),
	})
	return l.ConfirmCommissionMonthReport(in)
}

func (s *TransactionServer) CalculateMonthProfitReport(ctx context.Context, in *transactionclient.CalculateMonthProfitReportRequest) (*transactionclient.CalculateMonthProfitReportResponse, error) {
	l := logic.NewCalculateMonthProfitReportLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "CalculateMonthProfitReport",
		Value: attribute.StringValue(in.String()),
	})
	return l.CalculateMonthProfitReport(in)
}

func (s *TransactionServer) WithdrawCommissionOrderTransaction(ctx context.Context, in *transactionclient.WithdrawCommissionOrderRequest) (*transactionclient.WithdrawCommissionOrderResponse, error) {
	l := logic.NewWithdrawCommissionOrderTransactionLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "WithdrawCommissionOrderTransaction",
		Value: attribute.StringValue(in.String()),
	})
	return l.WithdrawCommissionOrderTransaction(in)
}

func (s *TransactionServer) WithdrawOrderToTest_XFB(ctx context.Context, in *transactionclient.WithdrawOrderTestRequest) (*transactionclient.WithdrawOrderTestResponse, error) {
	l := logic.NewWithdrawOrderToTestXFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "WithdrawOrderToTest_XFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.WithdrawOrderToTest_XFB(in)
}

func (s *TransactionServer) WithdrawTestToNormal_XFB(ctx context.Context, in *transactionclient.WithdrawOrderTestRequest) (*transactionclient.WithdrawOrderTestResponse, error) {
	l := logic.NewWithdrawTestToNormalXFBLogic(ctx, s.svcCtx)
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "WithdrawTestToNormal_XFB",
		Value: attribute.StringValue(in.String()),
	})
	return l.WithdrawTestToNormal_XFB(in)
}
