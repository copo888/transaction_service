package logic

import (
	"context"
	"fmt"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"github.com/gioco-play/easy-i18n/i18n"
	"github.com/neccoys/go-zero-extension/redislock"

	"github.com/zeromicro/go-zero/core/logx"
)

type MerchantBalanceFreezeTranactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMerchantBalanceFreezeTranactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MerchantBalanceFreezeTranactionLogic {
	return &MerchantBalanceFreezeTranactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MerchantBalanceFreezeTranactionLogic) MerchantBalanceFreezeTranaction(req *transactionclient.MerchantBalanceFreezeRequest) (*transactionclient.MerchantBalanceFreezeResponse, error) {
	newOrderNo := model.GenerateOrderNo("TJ")
	var transactionType string

	if req.Amount > 0 {
		transactionType = constants.TRANSACTION_TYPE_FREEZE
	} else {
		transactionType = constants.TRANSACTION_TYPE_UNFREEZE
	}

	redisKey := fmt.Sprintf("%s-%s", req.MerchantCode, req.CurrencyCode)
	redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(8)
	if isOK, redisErr := redisLock.TryLockTimeout(8); isOK {
		defer redisLock.Release()
		/****     交易開始      ****/
		txDB := l.svcCtx.MyDB.Begin()

		// 變更 商戶餘額&凍結金額 並記錄
		if _, err := merchantbalanceservice.FrozenManually(txDB, types.FrozenManually{
			MerchantCode:    req.MerchantCode,
			CurrencyCode:    req.CurrencyCode,
			OrderNo:         newOrderNo,
			OrderType:       "TJ",
			TransactionType: transactionType,
			BalanceType:     req.BalanceType,
			FrozenAmount:    req.Amount,
			Comment:         req.Comment,
			CreatedBy:       req.UserAccount,
		}, req.MerchantPtBalanceId); err != nil {
			txDB.Rollback()
			return &transactionclient.MerchantBalanceFreezeResponse{
				Code:    err.Error(),
				Message: "凍結金額異動失敗",
			}, nil
		}

		if err := txDB.Commit().Error; err != nil {
			txDB.Rollback()
			logx.Errorf("調整凍結Commit失败，MerchantCode: %s, CurrencyCode: %s, err : %s", req.MerchantCode, req.CurrencyCode, err.Error())
			return &transactionclient.MerchantBalanceFreezeResponse{
				Code:    response.DATABASE_FAILURE,
				Message: "资料库错误 Commit失败",
			}, nil
		}
	} else {
		logx.WithContext(l.ctx).Errorf("商户钱包处理中，Err:%s。 %s", redisErr.Error(), redisKey)
		return &transactionclient.MerchantBalanceFreezeResponse{
			Code:    response.BALANCE_PROCESSING,
			Message: i18n.Sprintf(response.BALANCE_PROCESSING),
		}, nil
	}

	return &transactionclient.MerchantBalanceFreezeResponse{
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}, nil
}
