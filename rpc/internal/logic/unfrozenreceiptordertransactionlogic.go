package logic

import (
	"context"
	"fmt"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"github.com/neccoys/go-zero-extension/redislock"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UnFrozenReceiptOrderTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnFrozenReceiptOrderTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnFrozenReceiptOrderTransactionLogic {
	return &UnFrozenReceiptOrderTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnFrozenReceiptOrderTransactionLogic) UnFrozenReceiptOrderTransaction(req *transactionclient.UnFrozenReceiptOrderRequest) (*transactionclient.UnFrozenReceiptOrderResponse, error) {
	var order types.Order
	myDB := l.svcCtx.MyDB

	// 1. 取得訂單
	if err := myDB.Table("tx_orders").Where("order_no = ?", req.OrderNo).Find(&order).Error; err != nil {
		return &transactionclient.UnFrozenReceiptOrderResponse{
			Code:    response.DATABASE_FAILURE,
			Message: "取得訂單失敗",
		}, nil
	}

	// 驗證
	if errCode := l.verify(order, req); errCode != "" {
		return &transactionclient.UnFrozenReceiptOrderResponse{
			Code:    errCode,
			Message: "驗證失敗: " + errCode,
		}, nil
	}

	redisKey := fmt.Sprintf("%s-%s", order.MerchantCode, order.CurrencyCode)
	redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(5)
	if isOK, _ := redisLock.TryLockTimeout(5); isOK {
		defer redisLock.Release()
		/****     交易開始      ****/
		txDB := myDB.Begin()
		// 變更 商戶餘額&凍結金額 並記錄
		if _, err := merchantbalanceservice.UpdateFrozenAmount(txDB, types.UpdateFrozenAmount{
			MerchantCode:    order.MerchantCode,
			CurrencyCode:    order.CurrencyCode,
			OrderNo:         order.OrderNo,
			MerchantOrderNo: order.MerchantOrderNo,
			OrderType:       order.Type,
			ChannelCode:     order.ChannelCode,
			PayTypeCode:     order.PayTypeCode,
			TransactionType: constants.TRANSACTION_TYPE_UNFREEZE,
			BalanceType:     order.BalanceType,
			FrozenAmount:    -order.FrozenAmount,
			Comment:         "订单解冻",
			CreatedBy:       req.UserAccount,
		}); err != nil {
			txDB.Rollback()
			return &transactionclient.UnFrozenReceiptOrderResponse{
				Code:    response.DATABASE_FAILURE,
				Message: "錢包異動失敗",
			}, nil
		}

		// 變更訂單狀態
		order.Status = constants.SUCCESS // 恢復成成功單
		order.FrozenAmount = 0
		order.Memo = "订单解冻 \n" + order.Memo

		if err := txDB.Select("Status", "FrozenAmount", "Memo").Table("tx_orders").Updates(&types.OrderX{
			Order: order,
		}).Error; err != nil {
			txDB.Rollback()
			return &transactionclient.UnFrozenReceiptOrderResponse{
				Code:    response.DATABASE_FAILURE,
				Message: "編輯訂單失敗",
			}, nil
		}

		if err := txDB.Commit().Error; err != nil {
			txDB.Rollback()
			logx.Errorf("解凍收款Commit失败，商户号: %s, 订单号: %s, err : %s", order.MerchantCode, order.OrderNo, err.Error())
			return &transactionclient.UnFrozenReceiptOrderResponse{
				Code:    response.DATABASE_FAILURE,
				Message: "资料库错误 Commit失败",
			}, nil
		}
		/****     交易結束      ****/
	}

	// 新單新增訂單歷程 (不抱錯)
	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     order.OrderNo,
			Action:      constants.ACTION_UNFROZEN,
			UserAccount: order.MerchantCode,
			Comment:     "订单解冻",
		},
	}).Error; err4 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err4.Error())
	}

	return &transactionclient.UnFrozenReceiptOrderResponse{
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}, nil
}

func (l *UnFrozenReceiptOrderTransactionLogic) verify(order types.Order, req *transactionclient.UnFrozenReceiptOrderRequest) string {

	// 凍結狀態 才能解凍
	if order.Status != constants.FROZEN {
		return response.ORDER_STATUS_WRONG_CANNOT_UNFROZEN
	}

	return ""
}
