// Code generated by goctl. DO NOT EDIT!
// Source: transaction.proto

package transaction

import (
	"context"

	"github.com/copo888/transaction_service/rpc/transaction"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CorrespondMerChnRate          = transaction.CorrespondMerChnRate
	InternalOrder                 = transaction.InternalOrder
	InternalOrderRequest          = transaction.InternalOrderRequest
	InternalOrderResponse         = transaction.InternalOrderResponse
	InternalReviewSuccessRequest  = transaction.InternalReviewSuccessRequest
	InternalReviewSuccessResponse = transaction.InternalReviewSuccessResponse
	MerchantOrderRateListView     = transaction.MerchantOrderRateListView
	PayCallBackRequest            = transaction.PayCallBackRequest
	PayCallBackResponse           = transaction.PayCallBackResponse
	PayOrder                      = transaction.PayOrder
	PayOrderRequest               = transaction.PayOrderRequest
	PayOrderResponse              = transaction.PayOrderResponse
	ProxyOrderRequest             = transaction.ProxyOrderRequest
	ProxyOrderResp_XFB            = transaction.ProxyOrderResp_XFB
	ProxyOrderResponse            = transaction.ProxyOrderResponse
	ProxyPayOrderRequest          = transaction.ProxyPayOrderRequest
	WithdrawOrderRequest          = transaction.WithdrawOrderRequest
	WithdrawOrderResponse         = transaction.WithdrawOrderResponse
	WithdrawReviewFailRequest     = transaction.WithdrawReviewFailRequest
	WithdrawReviewFailResponse    = transaction.WithdrawReviewFailResponse

	Transaction interface {
		ProxyOrderTranaction_DFB(ctx context.Context, in *ProxyOrderRequest, opts ...grpc.CallOption) (*ProxyOrderResponse, error)
		ProxyOrderTranaction_XFB(ctx context.Context, in *ProxyOrderRequest, opts ...grpc.CallOption) (*ProxyOrderResponse, error)
		PayOrderTranaction(ctx context.Context, in *PayOrderRequest, opts ...grpc.CallOption) (*PayOrderResponse, error)
		InternalOrderTransaction(ctx context.Context, in *InternalOrderRequest, opts ...grpc.CallOption) (*InternalOrderResponse, error)
		WithdrawOrderTransaction(ctx context.Context, in *WithdrawOrderRequest, opts ...grpc.CallOption) (*WithdrawOrderResponse, error)
		PayCallBackTranaction(ctx context.Context, in *PayCallBackRequest, opts ...grpc.CallOption) (*PayCallBackResponse, error)
		InternalReviewSuccessTransaction(ctx context.Context, in *InternalReviewSuccessRequest, opts ...grpc.CallOption) (*InternalReviewSuccessResponse, error)
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

func (m *defaultTransaction) ProxyOrderTranaction_DFB(ctx context.Context, in *ProxyOrderRequest, opts ...grpc.CallOption) (*ProxyOrderResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.ProxyOrderTranaction_DFB(ctx, in, opts...)
}

func (m *defaultTransaction) ProxyOrderTranaction_XFB(ctx context.Context, in *ProxyOrderRequest, opts ...grpc.CallOption) (*ProxyOrderResponse, error) {
	client := transaction.NewTransactionClient(m.cli.Conn())
	return client.ProxyOrderTranaction_XFB(ctx, in, opts...)
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
