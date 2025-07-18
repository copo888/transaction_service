package logic

import (
	"context"
	"fmt"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/redislock"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"github.com/gioco-play/easy-i18n/i18n"
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
	newOrderNo := model.GenerateOrderNo("TJ")
	updateBalance := types.UpdateBalance{
		MerchantCode:    req.MerchantCode,
		CurrencyCode:    req.CurrencyCode,
		OrderNo:         newOrderNo,
		MerchantOrderNo: "",
		OrderType:       "TJ",
		PayTypeCode:     "",
		PayTypeCodeNum:  "",
		TransferAmount:  req.Amount,
		TransactionType: constants.TRANSACTION_TYPE_ADJUST, //異動類型 (20=調整)
		BalanceType:     req.BalanceType,
		Comment:         req.Comment,
		CreatedBy:       req.UserAccount,
	}

	redisKey := fmt.Sprintf("%s-%s", updateBalance.MerchantCode, updateBalance.CurrencyCode)
	redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(8)
	if isOK, redisErr := redisLock.TryLockTimeout(8); isOK {
		defer redisLock.Release()

		if err := l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {
			_, err = merchantbalanceservice.UpdateBalance(db, updateBalance)
			return
		}); err != nil {
			return &transactionclient.MerchantBalanceUpdateResponse{
				Code:    response.SYSTEM_ERROR,
				Message: "更新錢包失敗",
			}, err
		}

	} else {
		logx.WithContext(l.ctx).Errorf("商户钱包处理中，Err:%s。 %s", redisErr.Error(), redisKey)
		return &transactionclient.MerchantBalanceUpdateResponse{
			Code:    response.BALANCE_PROCESSING,
			Message: i18n.Sprintf(response.BALANCE_PROCESSING),
		}, nil
	}

	return &transactionclient.MerchantBalanceUpdateResponse{
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}, nil
}
