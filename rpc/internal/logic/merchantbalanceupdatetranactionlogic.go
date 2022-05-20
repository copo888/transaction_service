package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"gorm.io/gorm"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantBalanceUpdateTranactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMerchantBalanceUpdateTranactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantBalanceUpdateTranactionLogic {
	return &MerchantBalanceUpdateTranactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MerchantBalanceUpdateTranactionLogic) MerchantBalanceUpdateTranaction(req *transactionclient.MerchantBalanceUpdateRequest) (*transactionclient.MerchantBalanceUpdateResponse, error) {
	updateBalance := types.UpdateBalance{
		MerchantCode:    req.MerchantCode,
		CurrencyCode:    req.CurrencyCode,
		OrderNo:         "",
		OrderType:       "TJ",
		PayTypeCode:     "",
		PayTypeCodeNum:  "",
		TransferAmount:  req.Amount,
		TransactionType: "20", //異動類型 (20=調整)
		BalanceType:     req.BalanceType,
		Comment:         req.Comment,
		CreatedBy:       req.UserAccount,
	}

	if err := l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {
		_, err = merchantbalanceservice.UpdateBalance(db, updateBalance)
		return
	}); err != nil {
		return &transactionclient.MerchantBalanceUpdateResponse{
			Code:    response.SYSTEM_ERROR,
			Message: "更新錢包失敗",
		}, err
	}

	return &transactionclient.MerchantBalanceUpdateResponse{
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}, nil
}
