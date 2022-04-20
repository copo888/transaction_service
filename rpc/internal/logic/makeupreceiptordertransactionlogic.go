package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/service/orderfeeprofitservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"gorm.io/gorm"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/transaction"

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

func (l *MakeUpReceiptOrderTransactionLogic) MakeUpReceiptOrderTransaction(req *transaction.MakeUpReceiptOrderRequest) (*transaction.MakeUpReceiptOrderResponse, error) {
	var order types.Order
	var newOrder types.Order

	var merchantBalanceRecord types.MerchantBalanceRecord
	var transferAmount float64
	var newOrderNo string
	var comment string

	 l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

		// 1. 取得訂單
		if err = db.Table("tx_orders").Where("order_no = ?", req.OrderNo).Find(&order).Error; err != nil {
			return errorz.New(response.DATABASE_FAILURE, err.Error())
		}

		// 驗證
		if err = l.verifyMakeUpReceiptOrder(order, req); err != nil {
			return
		}

		newOrderNo = model.GenerateOrderNo(order.Type)

		// 計算交易金額 (訂單金額 - 手續費)
		transferHandlingFee := utils.FloatAdd(utils.FloatMul(utils.FloatDiv(req.Amount, 100), order.Fee), order.HandlingFee)
		 // 計算實際交易金額 = 訂單金額 - 手續費
		transferAmount = req.Amount - transferHandlingFee

		// 變更 商戶餘額並記錄

		if merchantBalanceRecord, err = merchantbalanceservice.UpdateBalance(db, types.UpdateBalance{
			MerchantCode:    order.MerchantCode,
			CurrencyCode:    order.CurrencyCode,
			OrderNo:         order.OrderNo,
			OrderType:       order.Type,
			ChannelCode:     order.ChannelCode,
			PayTypeCode:     order.PayTypeCode,
			PayTypeCodeNum:  order.PayTypeCodeNum,
			TransactionType: "1",
			BalanceType:     order.BalanceType,
			TransferAmount:  transferAmount,
			Comment:         comment,
			CreatedBy:       "AAA00061", // TODO: JWT取得
		}); err != nil {
			return
		}

		// 新增訂單
		newOrder = order
		newOrder.ID = 0
		newOrder.Status = "20"
		newOrder.SourceOrderNo = order.OrderNo
		newOrder.ChannelOrderNo = req.ChannelOrderNo
		newOrder.OrderNo = newOrderNo
		newOrder.OrderAmount = req.Amount
		newOrder.ActualAmount = req.Amount
		newOrder.BeforeBalance = merchantBalanceRecord.BeforeBalance
		newOrder.TransferAmount = merchantBalanceRecord.TransferAmount
		newOrder.Balance = merchantBalanceRecord.AfterBalance
		newOrder.IsLock = "1"
		newOrder.CallBackStatus = "1"
		newOrder.IsMerchantCallback = "2"
		newOrder.MakeUpType = req.MakeUpType
		newOrder.PersonProcessStatus = "10"
		newOrder.InternalChargeOrderPath = ""
		newOrder.HandlingFee = order.HandlingFee
		newOrder.Fee = order.Fee
		newOrder.TransferHandlingFee = transferHandlingFee
		newOrder.Memo = req.Comment
		newOrder.Source = "1"

		if err = db.Table("tx_orders").Create(&types.OrderX{
			Order:   newOrder,
			TransAt: types.JsonTime{}.New(),
		}).Error; err != nil {
			return errorz.New(response.DATABASE_FAILURE, err.Error())
		}

		// 新單新增訂單歷程
		if err = db.Table("tx_order_actions").Create(&types.OrderActionX{
			OrderAction: types.OrderAction{
				OrderNo:     newOrder.OrderNo,
				Action:      "MAKE_UP_ORDER",
				UserAccount: "AAA00061", // TODO: JWT取得
				Comment:     comment,
			},
		}).Error; err != nil {
			return errorz.New(response.DATABASE_FAILURE, err.Error())
		}

		// 舊單鎖定
		order.IsLock = "1"
		order.Memo = "补单锁定"
		if err = db.Table("tx_orders").Updates(&types.OrderX{
			Order: order,
		}).Error; err != nil {
			return errorz.New(response.DATABASE_FAILURE, err.Error())
		}

		// 舊單新增歷程
		if err = db.Table("tx_order_actions").Create(&types.OrderActionX{
			OrderAction: types.OrderAction{
				OrderNo:     order.OrderNo,
				Action:      "MAKE_UP_LOCK_ORDER",
				UserAccount: "AAA00061", // TODO: JWT取得
				Comment:     "",
			},
		}).Error; err != nil {
			return errorz.New(response.DATABASE_FAILURE, err.Error())
		}

		return nil
	})

	// 計算利潤
	if err := orderfeeprofitservice.CalculateOrderProfit(l.svcCtx.MyDB, types.CalculateProfit{
		MerchantCode:        order.MerchantCode,
		OrderNo:             newOrderNo,
		Type:                order.Type,
		CurrencyCode:        order.CurrencyCode,
		BalanceType:         order.BalanceType,
		ChannelCode:         order.ChannelCode,
		ChannelPayTypesCode: order.ChannelPayTypesCode,
		OrderAmount:         req.Amount,
	}); err != nil {
		logx.Error("計算利潤出錯:%s", err.Error())
	}

	return &transaction.MakeUpReceiptOrderResponse{}, nil
}


func (l *MakeUpReceiptOrderTransactionLogic) verifyMakeUpReceiptOrder(order types.Order, req *transaction.MakeUpReceiptOrderRequest) error {
	// 檢查訂單狀態 (處理中 成功 失敗) 才能補單
	if !(order.Status == "1" || order.Status == "20" || order.Status == "30") {
		return errorz.New(response.ORDER_STATUS_WRONG_CANNOT_MAKE_UP, "订单状态错误,不可补单")
	}

	if len(order.SourceOrderNo) > 0 {
		return errorz.New(response.ORDER_IS_MAKE_UP_ORDER_DONT_MAKE_UP, "此订单为补单, 不可再补单")
	}

	if req.Amount <= 0 {
		return errorz.New(response.AMOUNT_MUST_BE_GREATER_THAN_ZERO, "金额需大於零:")
	}

	return nil
}