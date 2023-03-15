package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantPtBalanceService"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/service/orderfeeprofitservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type RecoverReceiptOrderTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecoverReceiptOrderTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecoverReceiptOrderTransactionLogic {
	return &RecoverReceiptOrderTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RecoverReceiptOrderTransactionLogic) RecoverReceiptOrderTransaction(req *transactionclient.RecoverReceiptOrderRequest) (*transactionclient.RecoverReceiptOrderResponse, error) {
	var order types.Order
	var newOrder types.Order
	var transferAmount float64

	/****     交易開始      ****/
	txDB := l.svcCtx.MyDB.Begin()

	// 取得訂單
	if err := txDB.Table("tx_orders").Where("order_no = ?", req.OrderNo).Find(&order).Error; err != nil {
		txDB.Rollback()
		return &transactionclient.RecoverReceiptOrderResponse{
			Code:    response.DATABASE_FAILURE,
			Message: "取得訂單失敗",
		}, nil
	}

	// 驗證
	if errCode := l.verifyOrder(order, req); errCode != "" {
		txDB.Rollback()
		return &transactionclient.RecoverReceiptOrderResponse{
			Code:    errCode,
			Message: "驗證失敗: " + errCode,
		}, nil
	}

	newOrderNo := model.GenerateOrderNo(order.Type)
	// 計算交易手續費 (金額 / 100 * 費率 + 手續費)
	transferHandlingFee := -utils.FloatAdd(utils.FloatMul(utils.FloatDiv(req.Amount, 100), order.Fee), order.HandlingFee)
	// 計算實際交易金額 = 訂單金額 + 手續費
	transferAmount = -req.Amount - transferHandlingFee
	merchantOrderNo := order.MerchantOrderNo + "#R"

	updateBalance := types.UpdateBalance{
		MerchantCode:    order.MerchantCode,
		CurrencyCode:    order.CurrencyCode,
		OrderNo:         newOrderNo,
		MerchantOrderNo: merchantOrderNo,
		OrderType:       order.Type,
		ChannelCode:     order.ChannelCode,
		PayTypeCode:     order.PayTypeCode,
		TransactionType: constants.TRANSACTION_TYPE_RECOVER,
		BalanceType:     order.BalanceType,
		TransferAmount:  transferAmount,
		Comment:         req.Comment,
		CreatedBy:       req.UserAccount,
	}

	// 變更 商戶餘額並記錄
	merchantBalanceRecord, err := merchantbalanceservice.UpdateBalanceForZF(txDB, l.svcCtx.RedisClient, updateBalance)
	if err != nil {
		txDB.Rollback()
		logx.WithContext(l.ctx).Errorf("更新錢包失敗:%s", err.Error())
		return &transactionclient.RecoverReceiptOrderResponse{
			Code:    response.SYSTEM_ERROR,
			Message: "更新錢包失敗",
		}, nil
	}

	var isDisplayBalance string
	err = txDB.Table("mc_merchant_channel_rate").
		Select("is_display_balance").
		Where("merchant_code = ? AND currency_code = ?", order.MerchantCode, order.CurrencyCode).
		Where("channel_code = ? AND pay_type_code = ?", order.ChannelCode, order.PayTypeCode).
		Find(&isDisplayBalance).Error

	// 若有啟用顯示子錢包
	if isDisplayBalance == "1" {
		logx.WithContext(l.ctx).Infof("需異動子錢包: %s,%s,%s.%s", order.MerchantCode, order.CurrencyCode, order.ChannelCode, order.PayTypeCode)
		// 變更 商戶子錢包餘額
		_, err = merchantPtBalanceService.UpdatePtBalanceForZF(txDB, l.svcCtx.RedisClient, updateBalance)
		if err != nil {
			txDB.Rollback()
			logx.WithContext(l.ctx).Errorf("異動子錢包失敗:%s", err.Error())
			return &transactionclient.RecoverReceiptOrderResponse{
				Code:    response.SYSTEM_ERROR,
				Message: "更新子錢包失敗",
			}, nil
		}
	}

	// 新增訂單
	newOrder = order
	newOrder.ID = 0
	newOrder.Status = constants.SUCCESS
	newOrder.SourceOrderNo = order.OrderNo
	newOrder.ChannelOrderNo = ""
	newOrder.MerchantOrderNo = merchantOrderNo
	newOrder.OrderNo = newOrderNo
	newOrder.OrderAmount = order.OrderAmount
	newOrder.ActualAmount = -req.Amount
	newOrder.BeforeBalance = merchantBalanceRecord.BeforeBalance
	newOrder.TransferAmount = merchantBalanceRecord.TransferAmount
	newOrder.Balance = merchantBalanceRecord.AfterBalance
	newOrder.IsLock = constants.IS_LOCK_NO
	newOrder.CallBackStatus = constants.CALL_BACK_STATUS_PROCESSING
	newOrder.IsMerchantCallback = constants.IS_MERCHANT_CALLBACK_NOT_NEED
	newOrder.ReasonType = "11"
	newOrder.PersonProcessStatus = constants.PERSON_PROCESS_STATUS_NO_ROCESSING
	newOrder.InternalChargeOrderPath = ""
	newOrder.HandlingFee = order.HandlingFee
	newOrder.Fee = order.Fee
	newOrder.TransferHandlingFee = transferHandlingFee
	newOrder.Memo = "原始单号:" + order.OrderNo + " \n" + req.Comment
	newOrder.Source = constants.ORDER_SOURCE_BY_PLATFORM
	newOrder.IsCalculateProfit = constants.IS_CALCULATE_PROFIT_YES

	if err = txDB.Table("tx_orders").Create(&types.OrderX{
		Order:   newOrder,
		TransAt: types.JsonTime{}.New(),
	}).Error; err != nil {
		txDB.Rollback()
		return &transactionclient.RecoverReceiptOrderResponse{
			Code:    response.SYSTEM_ERROR,
			Message: "新增訂單失敗",
		}, nil
	}

	// 舊單鎖定
	order.IsLock = "1"
	order.Memo = "追回单号:" + newOrderNo + " \n" + order.Memo
	if err = txDB.Table("tx_orders").Updates(&types.OrderX{
		Order: order,
	}).Error; err != nil {
		txDB.Rollback()
		return &transactionclient.RecoverReceiptOrderResponse{
			Code:    response.SYSTEM_ERROR,
			Message: "舊單鎖定失敗",
		}, nil
	}

	// 追回單利潤計算方式不同 要放在交易裡
	if err = orderfeeprofitservice.CalculateSubOrderProfit(txDB, types.CalculateSubOrderProfit{
		OldOrderNo:            order.OrderNo,
		NewOrderNo:            newOrderNo,
		OrderAmount:           -req.Amount,
		IsCalculateCommission: req.IsCalculateCommission,
	}); err != nil {
		txDB.Rollback()
		return &transactionclient.RecoverReceiptOrderResponse{
			Code:    response.SYSTEM_ERROR,
			Message: "計算利潤出錯",
		}, nil
	}

	if err := txDB.Commit().Error; err != nil {
		txDB.Rollback()
		logx.Errorf("支付追回失败，商户号: %s, 订单号: %s, err : %s", order.MerchantCode, order.OrderNo, err.Error())
		return &transactionclient.RecoverReceiptOrderResponse{
			Code:    response.DATABASE_FAILURE,
			Message: "资料库错误 Commit失败",
		}, nil
	}
	/****     交易結束      ****/

	// 舊單新增歷程
	if err := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     order.OrderNo,
			Action:      "MAKE_UP_LOCK_ORDER",
			UserAccount: req.UserAccount,
			Comment:     "",
		},
	}).Error; err != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err.Error())
	}

	// 新單新增訂單歷程
	if err := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     newOrder.OrderNo,
			Action:      "MAKE_UP_ORDER",
			UserAccount: req.UserAccount,
			Comment:     req.Comment,
		},
	}).Error; err != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err.Error())
	}

	return &transactionclient.RecoverReceiptOrderResponse{
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}, nil
}

func (l *RecoverReceiptOrderTransactionLogic) verifyOrder(order types.Order, req *transactionclient.RecoverReceiptOrderRequest) string {

	// 支付單才能追回
	if order.Type != constants.ORDER_TYPE_ZF {
		return response.ORDER_TYPE_IS_WRONG
	}

	// 檢查訂單狀態 成功單才能追回
	if order.Status != constants.SUCCESS {
		return response.ORDER_STATUS_WRONG
	}

	// 鎖定單不可追回 (已補單 已追回)
	if order.IsLock == constants.IS_LOCK_YES {
		return response.ORDER_IS_STATUS_IS_LOCK
	}

	// 追回單不可再追回
	if order.ReasonType == constants.ORDER_REASON_TYPE_RECOVER {
		return response.ORDER_STATUS_WRONG
	}

	// 訂單還未計算傭金,請稍後
	if order.IsCalculateProfit != constants.IS_CALCULATE_PROFIT_YES {
		return response.ORIGINAL_ORDER_NOT_CALCULATED_COMMISSION
	}

	if req.Amount <= 0 {
		return response.AMOUNT_MUST_BE_GREATER_THAN_ZERO
	}

	return ""
}
