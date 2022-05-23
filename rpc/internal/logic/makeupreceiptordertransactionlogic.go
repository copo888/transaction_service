package logic

import (
	"context"

	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/service/orderfeeprofitservice"
	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type MakeUpReceiptOrderTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMakeUpReceiptOrderTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MakeUpReceiptOrderTransactionLogic {
	return &MakeUpReceiptOrderTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MakeUpReceiptOrderTransactionLogic) MakeUpReceiptOrderTransaction(req *transactionclient.MakeUpReceiptOrderRequest) (*transactionclient.MakeUpReceiptOrderResponse, error) {
	var order types.Order
	var newOrder types.Order
	var transferAmount float64
	var newOrderNo string

	/****     交易開始      ****/
	txDB := l.svcCtx.MyDB.Begin()

	// 1. 取得訂單
	if err := txDB.Table("tx_orders").Where("order_no = ?", req.OrderNo).Find(&order).Error; err != nil {
		txDB.Rollback()
		return &transactionclient.MakeUpReceiptOrderResponse{
			Code:    response.DATABASE_FAILURE,
			Message: "取得訂單失敗",
		}, nil
	}

	// 驗證
	if errCode := l.verifyMakeUpReceiptOrder(order, req); errCode != "" {
		txDB.Rollback()
		return &transactionclient.MakeUpReceiptOrderResponse{
			Code:    errCode,
			Message: "驗證失敗: " + errCode,
		}, nil
	}

	newOrderNo = model.GenerateOrderNo(order.Type)

	// 計算交易手續費 (金額 / 100 * 費率 + 手續費)
	transferHandlingFee := utils.FloatAdd(utils.FloatMul(utils.FloatDiv(req.Amount, 100), order.Fee), order.HandlingFee)
	// 計算實際交易金額 = 訂單金額 - 手續費
	transferAmount = req.Amount - transferHandlingFee

	// 變更 商戶餘額並記錄

	merchantBalanceRecord, err := merchantbalanceservice.UpdateBalanceForZF(txDB, types.UpdateBalance{
		MerchantCode:    order.MerchantCode,
		CurrencyCode:    order.CurrencyCode,
		OrderNo:         newOrderNo,
		MerchantOrderNo: order.MerchantOrderNo,
		OrderType:       order.Type,
		ChannelCode:     order.ChannelCode,
		PayTypeCode:     order.PayTypeCode,
		TransactionType: constants.TRANSACTION_TYPE_MAKE_UP,
		BalanceType:     order.BalanceType,
		TransferAmount:  transferAmount,
		Comment:         req.Comment,
		CreatedBy:       req.UserAccount,
	})
	if err != nil {
		txDB.Rollback()
		return &transactionclient.MakeUpReceiptOrderResponse{
			Code:    response.SYSTEM_ERROR,
			Message: "更新錢包失敗",
		}, nil
	}

	// 新增訂單
	newOrder = order
	newOrder.ID = 0
	newOrder.Status = constants.SUCCESS
	newOrder.SourceOrderNo = order.OrderNo
	newOrder.ChannelOrderNo = req.ChannelOrderNo
	newOrder.MerchantOrderNo = order.MerchantOrderNo + "#M"
	newOrder.OrderNo = newOrderNo
	newOrder.OrderAmount = req.Amount
	newOrder.ActualAmount = req.Amount
	newOrder.BeforeBalance = merchantBalanceRecord.BeforeBalance
	newOrder.TransferAmount = merchantBalanceRecord.TransferAmount
	newOrder.Balance = merchantBalanceRecord.AfterBalance
	newOrder.IsLock = constants.IS_LOCK_NO
	newOrder.CallBackStatus = constants.CALL_BACK_STATUS_PROCESSING
	newOrder.IsMerchantCallback = constants.MERCHANT_CALL_BACK_YES
	newOrder.ReasonType = req.ReasonType
	newOrder.PersonProcessStatus = constants.PERSON_PROCESS_STATUS_NO_ROCESSING
	newOrder.InternalChargeOrderPath = ""
	newOrder.HandlingFee = order.HandlingFee
	newOrder.Fee = order.Fee
	newOrder.TransferHandlingFee = transferHandlingFee
	newOrder.Memo = req.Comment
	newOrder.Source = constants.ORDER_SOURCE_BY_PLATFORM
	newOrder.IsCalculateProfit = constants.IS_CALCULATE_PROFIT_YES

	if err = txDB.Table("tx_orders").Create(&types.OrderX{
		Order:   newOrder,
		TransAt: types.JsonTime{}.New(),
	}).Error; err != nil {
		txDB.Rollback()
		return &transactionclient.MakeUpReceiptOrderResponse{
			Code:    response.SYSTEM_ERROR,
			Message: "新增訂單失敗",
		}, nil
	}

	// 舊單鎖定
	order.IsLock = "1"
	order.Memo = "补单单号:" + newOrderNo + " \n" + order.Memo
	if err = txDB.Table("tx_orders").Updates(&types.OrderX{
		Order: order,
	}).Error; err != nil {
		txDB.Rollback()
		return &transactionclient.MakeUpReceiptOrderResponse{
			Code:    response.SYSTEM_ERROR,
			Message: "舊單鎖定失敗",
		}, nil
	}

	// 計算利潤
	if err = orderfeeprofitservice.CalculateOrderProfit(txDB, types.CalculateProfit{
		MerchantCode:        newOrder.MerchantCode,
		OrderNo:             newOrder.OrderNo,
		Type:                newOrder.Type,
		CurrencyCode:        newOrder.CurrencyCode,
		BalanceType:         newOrder.BalanceType,
		ChannelCode:         newOrder.ChannelCode,
		ChannelPayTypesCode: newOrder.ChannelPayTypesCode,
		OrderAmount:         newOrder.ActualAmount,
	}); err != nil {
		txDB.Rollback()
		return &transactionclient.MakeUpReceiptOrderResponse{
			Code:    response.SYSTEM_ERROR,
			Message: "計算利潤出錯",
		}, nil
	}

	if err := txDB.Commit().Error; err != nil {
		txDB.Rollback()
		logx.Errorf("支付補單失败，商户号: %s, 订单号: %s, err : %s", order.MerchantCode, order.OrderNo, err.Error())
		return &transactionclient.MakeUpReceiptOrderResponse{
			Code:    response.DATABASE_FAILURE,
			Message: "资料库错误 Commit失败",
		}, nil
	}
	/****     交易結束      ****/

	// 舊單新增歷程
	if err := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     order.OrderNo,
			Action:      constants.ACTION_MAKE_UP_LOCK_ORDER,
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
			Action:      constants.ACTION_MAKE_UP_ORDER,
			UserAccount: req.UserAccount,
			Comment:     req.Comment,
		},
	}).Error; err != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err.Error())
	}

	return &transactionclient.MakeUpReceiptOrderResponse{
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}, nil
}

func (l *MakeUpReceiptOrderTransactionLogic) verifyMakeUpReceiptOrder(order types.Order, req *transactionclient.MakeUpReceiptOrderRequest) string {

	// 支付單才能補單
	if order.Type != constants.ORDER_TYPE_ZF {
		return response.ORDER_TYPE_IS_WRONG
	}

	// 檢查訂單狀態 (處理中 成功 失敗) 才能補單
	if !(order.Status == "1" || order.Status == "20" || order.Status == "30") {
		return response.ORDER_STATUS_WRONG_CANNOT_MAKE_UP
	}

	if order.IsLock == constants.IS_LOCK_YES {
		return response.ORDER_IS_STATUS_IS_LOCK
	}

	// 成功單還未計算傭金,請稍後
	if order.Status == "20" && order.IsCalculateProfit != constants.IS_CALCULATE_PROFIT_YES {
		return response.ORIGINAL_ORDER_NOT_CALCULATED_COMMISSION
	}

	if req.Amount <= 0 {
		return response.AMOUNT_MUST_BE_GREATER_THAN_ZERO
	}

	return ""
}
