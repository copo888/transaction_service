package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/service/orderfeeprofitservice"
	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transaction"
	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecalculateProfitTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecalculateProfitTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecalculateProfitTransactionLogic {
	return &RecalculateProfitTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RecalculateProfitTransactionLogic) RecalculateProfitTransaction(in *transaction.RecalculateProfitRequest) (*transaction.RecalculateProfitResponse, error) {
	errNum := 0
	okNum := 0
	for _, profit := range in.List {
		var calculateProfit types.CalculateProfit
		copier.Copy(&calculateProfit, &profit)
		if err := orderfeeprofitservice.CalculateOrderProfit(l.svcCtx.MyDB, calculateProfit); err != nil {
			errNum += 1
		} else {
			okNum += 1
		}
	}
	logx.Infof("補算利潤總筆數:%d, 成功數:%d, 失敗數:%d", len(in.List), okNum, errNum)

	return &transaction.RecalculateProfitResponse{
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}, nil
}
