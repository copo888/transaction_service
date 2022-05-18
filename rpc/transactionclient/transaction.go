// Code generated by goctl. DO NOT EDIT!
// Source: transaction.proto

package transactionclient

import (
	"context"

	"github.com/copo888/transaction_service/rpc/transaction"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CalculateCommissionMonthAllRequest       = transaction.CalculateCommissionMonthAllRequest
	CalculateCommissionMonthAllResponse      = transaction.CalculateCommissionMonthAllResponse
	CalculateProfit                          = transaction.CalculateProfit
	ChannelWithdraw                          = transaction.ChannelWithdraw
	ConfirmCommissionMonthReportRequest      = transaction.ConfirmCommissionMonthReportRequest
	ConfirmCommissionMonthReportResponse     = transaction.ConfirmCommissionMonthReportResponse
	ConfirmPayOrderRequest                   = transaction.ConfirmPayOrderRequest
	ConfirmPayOrderResponse                  = transaction.ConfirmPayOrderResponse
	CorrespondMerChnRate                     = transaction.CorrespondMerChnRate
	FrozenReceiptOrderRequest                = transaction.FrozenReceiptOrderRequest
	FrozenReceiptOrderResponse               = transaction.FrozenReceiptOrderResponse
	InternalOrder                            = transaction.InternalOrder
	InternalOrderRequest                     = transaction.InternalOrderRequest
	InternalOrderResponse                    = transaction.InternalOrderResponse
	InternalReviewSuccessRequest             = transaction.InternalReviewSuccessRequest
	InternalReviewSuccessResponse            = transaction.InternalReviewSuccessResponse
	MakeUpReceiptOrderRequest                = transaction.MakeUpReceiptOrderRequest
	MakeUpReceiptOrderResponse               = transaction.MakeUpReceiptOrderResponse
	MerchantBalanceUpdateRequest             = transaction.MerchantBalanceUpdateRequest
	MerchantBalanceUpdateResponse            = transaction.MerchantBalanceUpdateResponse
	MerchantOrderRateListView                = transaction.MerchantOrderRateListView
	PayCallBackRequest                       = transaction.PayCallBackRequest
	PayCallBackResponse                      = transaction.PayCallBackResponse
	PayOrder                                 = transaction.PayOrder
	PayOrderRequest                          = transaction.PayOrderRequest
	PayOrderResponse                         = transaction.PayOrderResponse
	PayOrderSwitchTestRequest                = transaction.PayOrderSwitchTestRequest
	PayOrderSwitchTestResponse               = transaction.PayOrderSwitchTestResponse
	PersonalRebundRequest                    = transaction.PersonalRebundRequest
	PersonalRebundResponse                   = transaction.PersonalRebundResponse
	ProxyOrderRequest                        = transaction.ProxyOrderRequest
	ProxyOrderResponse                       = transaction.ProxyOrderResponse
	ProxyOrderTestRequest                    = transaction.ProxyOrderTestRequest
	ProxyOrderTestResponse                   = transaction.ProxyOrderTestResponse
	ProxyOrderUI                             = transaction.ProxyOrderUI
	ProxyOrderUIRequest                      = transaction.ProxyOrderUIRequest
	ProxyOrderUIResponse                     = transaction.ProxyOrderUIResponse
	ProxyPayFailRequest                      = transaction.ProxyPayFailRequest
	ProxyPayFailResponse                     = transaction.ProxyPayFailResponse
	ProxyPayOrderRequest                     = transaction.ProxyPayOrderRequest
	RecalculateCommissionMonthReportRequest  = transaction.RecalculateCommissionMonthReportRequest
	RecalculateCommissionMonthReportResponse = transaction.RecalculateCommissionMonthReportResponse
	RecalculateProfitRequest                 = transaction.RecalculateProfitRequest
	RecalculateProfitResponse                = transaction.RecalculateProfitResponse
	RecoverReceiptOrderRequest               = transaction.RecoverReceiptOrderRequest
	RecoverReceiptOrderResponse              = transaction.RecoverReceiptOrderResponse
	UnFrozenReceiptOrderRequest              = transaction.UnFrozenReceiptOrderRequest
	UnFrozenReceiptOrderResponse             = transaction.UnFrozenReceiptOrderResponse
	WithdrawOrderRequest                     = transaction.WithdrawOrderRequest
	WithdrawOrderResponse                    = transaction.WithdrawOrderResponse
	WithdrawReviewFailRequest                = transaction.WithdrawReviewFailRequest
	WithdrawReviewFailResponse               = transaction.WithdrawReviewFailResponse
	WithdrawReviewSuccessRequest             = transaction.WithdrawReviewSuccessRequest
	WithdrawReviewSuccessResponse            = transaction.WithdrawReviewSuccessResponse

	Transaction interface {
		MerchantBalanceUpdateTranaction(ctx context.Context, in *MerchantBalanceUpdateRequest, opts ...grpc.CallOption) (*MerchantBalanceUpdateResponse, error)
		ProxyOrderTranaction_DFB(ctx context.Context, in *ProxyOrderRequest, opts ...grpc.CallOption) (*ProxyOrderResponse, error)
		ProxyOrderTranaction_XFB(ctx context.Context, in *ProxyOrderRequest, opts ...grpc.CallOption) (*ProxyOrderResponse, error)
		ProxyOrderTransactionFail_DFB(ctx context.Context, in *ProxyPayFailRequest, opts ...grpc.CallOption) (*ProxyPayFailResponse, error)
		ProxyOrderTransactionFail_XFB(ctx context.Context, in *ProxyPayFailRequest, opts ...grpc.CallOption) (*ProxyPayFailResponse, error)
		PayOrderSwitchTest(ctx context.Context, in *PayOrderSwitchTestRequest, opts ...grpc.CallOption) (*PayOrderSwitchTestResponse, error)
		ProxyOrderToTest_DFB(ctx context.Context, in *ProxyOrderTestRequest, opts ...grpc.CallOption) (*ProxyOrderTestResponse, error)
		ProxyOrderToTest_XFB(ctx context.Context, in *ProxyOrderTestRequest, opts ...grpc.CallOption) (*ProxyOrderTestResponse, error)
		ProxyTestToNormal_DFB(ctx context.Context, in *ProxyOrderTestRequest, opts ...grpc.CallOption) (*ProxyOrderTestResponse, error)
		ProxyTestToNormal_XFB(ctx context.Context, in *ProxyOrderTestRequest, opts ...grpc.CallOption) (*ProxyOrderTestResponse, error)
		PayOrderTranaction(ctx context.Context, in *PayOrderRequest, opts ...grpc.CallOption) (*PayOrderResponse, error)
		InternalOrderTransaction(ctx context.Context, in *InternalOrderRequest, opts ...grpc.CallOption) (*InternalOrderResponse, error)
		WithdrawOrderTransaction(ctx context.Context, in *WithdrawOrderRequest, opts ...grpc.CallOption) (*WithdrawOrderResponse, error)
		PayCallBackTranaction(ctx context.Context, in *PayCallBackRequest, opts ...grpc.CallOption) (*PayCallBackResponse, error)
		InternalReviewSuccessTransaction(ctx context.Context, in *InternalReviewSuccessRequest, opts ...grpc.CallOption) (*InternalReviewSuccessResponse, error)
		WithdrawReviewFailTransaction(ctx context.Context, in *WithdrawReviewFailRequest, opts ...grpc.CallOption) (*WithdrawReviewFailResponse, error)
		WithdrawReviewSuccessTransaction(ctx context.Context, in *WithdrawReviewSuccessRequest, opts ...grpc.CallOption) (*WithdrawReviewSuccessResponse, error)
		ProxyOrderUITransaction_DFB(ctx context.Context, in *ProxyOrderUIRequest, opts ...grpc.CallOption) (*ProxyOrderUIResponse, error)
		ProxyOrderUITransaction_XFB(ctx context.Context, in *ProxyOrderUIRequest, opts ...grpc.CallOption) (*ProxyOrderUIResponse, error)
		MakeUpReceiptOrderTransaction(ctx context.Context, in *MakeUpReceiptOrderRequest, opts ...grpc.CallOption) (*MakeUpReceiptOrderResponse, error)
		ConfirmPayOrderTransaction(ctx context.Context, in *ConfirmPayOrderRequest, opts ...grpc.CallOption) (*ConfirmPayOrderResponse, error)
		RecoverReceiptOrderTransaction(ctx context.Context, in *RecoverReceiptOrderRequest, opts ...grpc.CallOption) (*RecoverReceiptOrderResponse, error)
		FrozenReceiptOrderTransaction(ctx context.Context, in *FrozenReceiptOrderRequest, opts ...grpc.CallOption) (*FrozenReceiptOrderResponse, error)
		UnFrozenReceiptOrderTransaction(ctx context.Context, in *UnFrozenReceiptOrderRequest, opts ...grpc.CallOption) (*UnFrozenReceiptOrderResponse, error)
		PersonalRebundTransaction_DFB(ctx context.Context, in *PersonalRebundRequest, opts ...grpc.CallOption) (*PersonalRebundResponse, error)
		PersonalRebundTransaction_XFB(ctx context.Context, in *PersonalRebundRequest, opts ...grpc.CallOption) (*PersonalRebundResponse, error)
		RecalculateProfitTransaction(ctx context.Context, in *RecalculateProfitRequest, opts ...grpc.CallOption) (*RecalculateProfitResponse, error)
		CalculateCommissionMonthAllReport(ctx context.Context, in *CalculateCommissionMonthAllRequest, opts ...grpc.CallOption) (*CalculateCommissionMonthAllResponse, error)
		RecalculateCommissionMonthReport(ctx context.Context, in *RecalculateCommissionMonthReportRequest, opts ...grpc.CallOption) (*RecalculateCommissionMonthReportResponse, error)
		ConfirmCommissionMonthReport(ctx context.Context, in *ConfirmCommissionMonthReportRequest, opts ...grpc.CallOption) (*ConfirmCommissionMonthReportResponse, error)
	}

	defaultTransaction struct {
		cli zrpc.Client
	}
)

func NewTransaction(cli zrpc.Client) Transaction {
	return &defaultTransaction{
		cli: cli,
	}
}

func (m *defaultTransaction) MerchantBalanceUpdateTranaction(ctx context.Context, in *MerchantBalanceUpdateRequest, opts ...grpc.CallOption) (*MerchantBalanceUpdateResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.MerchantBalanceUpdateTranaction(ctx, in, opts...)
}

func (m *defaultTransaction) ProxyOrderTranaction_DFB(ctx context.Context, in *ProxyOrderRequest, opts ...grpc.CallOption) (*ProxyOrderResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.ProxyOrderTranaction_DFB(ctx, in, opts...)
}

func (m *defaultTransaction) ProxyOrderTranaction_XFB(ctx context.Context, in *ProxyOrderRequest, opts ...grpc.CallOption) (*ProxyOrderResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.ProxyOrderTranaction_XFB(ctx, in, opts...)
}

func (m *defaultTransaction) ProxyOrderTransactionFail_DFB(ctx context.Context, in *ProxyPayFailRequest, opts ...grpc.CallOption) (*ProxyPayFailResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.ProxyOrderTransactionFail_DFB(ctx, in, opts...)
}

func (m *defaultTransaction) ProxyOrderTransactionFail_XFB(ctx context.Context, in *ProxyPayFailRequest, opts ...grpc.CallOption) (*ProxyPayFailResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.ProxyOrderTransactionFail_XFB(ctx, in, opts...)
}

func (m *defaultTransaction) PayOrderSwitchTest(ctx context.Context, in *PayOrderSwitchTestRequest, opts ...grpc.CallOption) (*PayOrderSwitchTestResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.PayOrderSwitchTest(ctx, in, opts...)
}

func (m *defaultTransaction) ProxyOrderToTest_DFB(ctx context.Context, in *ProxyOrderTestRequest, opts ...grpc.CallOption) (*ProxyOrderTestResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.ProxyOrderToTest_DFB(ctx, in, opts...)
}

func (m *defaultTransaction) ProxyOrderToTest_XFB(ctx context.Context, in *ProxyOrderTestRequest, opts ...grpc.CallOption) (*ProxyOrderTestResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.ProxyOrderToTest_XFB(ctx, in, opts...)
}

func (m *defaultTransaction) ProxyTestToNormal_DFB(ctx context.Context, in *ProxyOrderTestRequest, opts ...grpc.CallOption) (*ProxyOrderTestResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.ProxyTestToNormal_DFB(ctx, in, opts...)
}

func (m *defaultTransaction) ProxyTestToNormal_XFB(ctx context.Context, in *ProxyOrderTestRequest, opts ...grpc.CallOption) (*ProxyOrderTestResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.ProxyTestToNormal_XFB(ctx, in, opts...)
}

func (m *defaultTransaction) PayOrderTranaction(ctx context.Context, in *PayOrderRequest, opts ...grpc.CallOption) (*PayOrderResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.PayOrderTranaction(ctx, in, opts...)
}

func (m *defaultTransaction) InternalOrderTransaction(ctx context.Context, in *InternalOrderRequest, opts ...grpc.CallOption) (*InternalOrderResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.InternalOrderTransaction(ctx, in, opts...)
}

func (m *defaultTransaction) WithdrawOrderTransaction(ctx context.Context, in *WithdrawOrderRequest, opts ...grpc.CallOption) (*WithdrawOrderResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.WithdrawOrderTransaction(ctx, in, opts...)
}

func (m *defaultTransaction) PayCallBackTranaction(ctx context.Context, in *PayCallBackRequest, opts ...grpc.CallOption) (*PayCallBackResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.PayCallBackTranaction(ctx, in, opts...)
}

func (m *defaultTransaction) InternalReviewSuccessTransaction(ctx context.Context, in *InternalReviewSuccessRequest, opts ...grpc.CallOption) (*InternalReviewSuccessResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.InternalReviewSuccessTransaction(ctx, in, opts...)
}

func (m *defaultTransaction) WithdrawReviewFailTransaction(ctx context.Context, in *WithdrawReviewFailRequest, opts ...grpc.CallOption) (*WithdrawReviewFailResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.WithdrawReviewFailTransaction(ctx, in, opts...)
}

func (m *defaultTransaction) WithdrawReviewSuccessTransaction(ctx context.Context, in *WithdrawReviewSuccessRequest, opts ...grpc.CallOption) (*WithdrawReviewSuccessResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.WithdrawReviewSuccessTransaction(ctx, in, opts...)
}

func (m *defaultTransaction) ProxyOrderUITransaction_DFB(ctx context.Context, in *ProxyOrderUIRequest, opts ...grpc.CallOption) (*ProxyOrderUIResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.ProxyOrderUITransaction_DFB(ctx, in, opts...)
}

func (m *defaultTransaction) ProxyOrderUITransaction_XFB(ctx context.Context, in *ProxyOrderUIRequest, opts ...grpc.CallOption) (*ProxyOrderUIResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.ProxyOrderUITransaction_XFB(ctx, in, opts...)
}

func (m *defaultTransaction) MakeUpReceiptOrderTransaction(ctx context.Context, in *MakeUpReceiptOrderRequest, opts ...grpc.CallOption) (*MakeUpReceiptOrderResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.MakeUpReceiptOrderTransaction(ctx, in, opts...)
}

func (m *defaultTransaction) ConfirmPayOrderTransaction(ctx context.Context, in *ConfirmPayOrderRequest, opts ...grpc.CallOption) (*ConfirmPayOrderResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.ConfirmPayOrderTransaction(ctx, in, opts...)
}

func (m *defaultTransaction) RecoverReceiptOrderTransaction(ctx context.Context, in *RecoverReceiptOrderRequest, opts ...grpc.CallOption) (*RecoverReceiptOrderResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.RecoverReceiptOrderTransaction(ctx, in, opts...)
}

func (m *defaultTransaction) FrozenReceiptOrderTransaction(ctx context.Context, in *FrozenReceiptOrderRequest, opts ...grpc.CallOption) (*FrozenReceiptOrderResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.FrozenReceiptOrderTransaction(ctx, in, opts...)
}

func (m *defaultTransaction) UnFrozenReceiptOrderTransaction(ctx context.Context, in *UnFrozenReceiptOrderRequest, opts ...grpc.CallOption) (*UnFrozenReceiptOrderResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.UnFrozenReceiptOrderTransaction(ctx, in, opts...)
}

func (m *defaultTransaction) PersonalRebundTransaction_DFB(ctx context.Context, in *PersonalRebundRequest, opts ...grpc.CallOption) (*PersonalRebundResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.PersonalRebundTransaction_DFB(ctx, in, opts...)
}

func (m *defaultTransaction) PersonalRebundTransaction_XFB(ctx context.Context, in *PersonalRebundRequest, opts ...grpc.CallOption) (*PersonalRebundResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.PersonalRebundTransaction_XFB(ctx, in, opts...)
}

func (m *defaultTransaction) RecalculateProfitTransaction(ctx context.Context, in *RecalculateProfitRequest, opts ...grpc.CallOption) (*RecalculateProfitResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.RecalculateProfitTransaction(ctx, in, opts...)
}

func (m *defaultTransaction) CalculateCommissionMonthAllReport(ctx context.Context, in *CalculateCommissionMonthAllRequest, opts ...grpc.CallOption) (*CalculateCommissionMonthAllResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.CalculateCommissionMonthAllReport(ctx, in, opts...)
}

func (m *defaultTransaction) RecalculateCommissionMonthReport(ctx context.Context, in *RecalculateCommissionMonthReportRequest, opts ...grpc.CallOption) (*RecalculateCommissionMonthReportResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.RecalculateCommissionMonthReport(ctx, in, opts...)
}

func (m *defaultTransaction) ConfirmCommissionMonthReport(ctx context.Context, in *ConfirmCommissionMonthReportRequest, opts ...grpc.CallOption) (*ConfirmCommissionMonthReportResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.ConfirmCommissionMonthReport(ctx, in, opts...)
}
