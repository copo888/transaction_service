package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/service/commissionService"
	"github.com/copo888/transaction_service/rpc/transactionclient"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CalculateCommissionMonthAllReportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCalculateCommissionMonthAllReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CalculateCommissionMonthAllReportLogic {
	return &CalculateCommissionMonthAllReportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CalculateCommissionMonthAllReportLogic) CalculateCommissionMonthAllReport(in *transactionclient.CalculateCommissionMonthAllRequest) (*transactionclient.CalculateCommissionMonthAllResponse, error) {

	err := commissionService.CalculateMonthAllReport(l.svcCtx.MyDB, in.Month)
	if err != nil {
		return &transactionclient.CalculateCommissionMonthAllResponse{
			Code:    response.DATABASE_FAILURE,
			Message: err.Error(),
		}, nil
	}

	return &transactionclient.CalculateCommissionMonthAllResponse{
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}, nil
}
