package logic

import (
	"context"
	"fmt"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/service/orderfeeprofitservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"gorm.io/gorm"
	"math"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type PayCallBackTranactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayCallBackTranactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayCallBackTranactionLogic {
	return &PayCallBackTranactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayCallBackTranactionLogic) PayCallBackTranaction(in *transactionclient.PayCallBackRequest) (resp *transactionclient.PayCallBackResponse, err error) {
	var order *types.OrderX

	if err = l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

		if err = l.svcCtx.MyDB.Table("tx_orders").
			Where("order_no = ?", in.PayOrderNo).Take(&order).Error; err != nil {
			return errorz.New(response.ORDER_NUMBER_NOT_EXIST, err.Error())
		}

		// 處理中的且非鎖定訂單 才能回調
		if order.Status != "1" || order.IsLock == "1" {
			return errorz.New(response.TRANSACTION_FAILURE, fmt.Sprintf("(status/isLock): %s/%s", order.Status, order.IsLock))
		}

		// 下單金額及實付金額差異風控 (差異超過5% 且 超過1元)
		limit := utils.FloatMul(order.OrderAmount, 0.05)
		diff := math.Abs(utils.FloatSub(order.OrderAmount, in.OrderAmount))
		if diff > limit && diff > 1 {
			return errorz.New(response.ORDER_AMOUNT_ERROR, fmt.Sprintf("(orderAmount/payAmount): %f/%f", order.OrderAmount, in.OrderAmount))
		}

		// 編輯訂單 異動錢包和餘額
		if err = l.updateOrderAndBalance(db, in, order); err != nil {
			return errorz.New(response.SYSTEM_ERROR, err.Error())
		}

		return
	}); err != nil {
		return
	}

	// 計算利潤 (不抱錯) TODO: 異步??
	if err4 := orderfeeprofitservice.CalculateOrderProfit(l.svcCtx.MyDB, types.CalculateProfit{
		MerchantCode:        order.MerchantCode,
		OrderNo:             order.OrderNo,
		Type:                order.Type,
		CurrencyCode:        order.CurrencyCode,
		BalanceType:         order.BalanceType,
		ChannelCode:         order.ChannelCode,
		ChannelPayTypesCode: order.ChannelPayTypesCode,
		OrderAmount:         order.ActualAmount,
	}); err4 != nil {
		logx.Error("計算利潤出錯:%s", err4.Error())
	}

	// 新單新增訂單歷程 (不抱錯)
	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     order.OrderNo,
			Action:      "SUCCESS",
			UserAccount: order.MerchantCode,
			Comment:     "",
		},
	}).Error; err4 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err4.Error())
	}

	return &transactionclient.PayCallBackResponse{
		MerchantCode:        order.MerchantCode,
		MerchantOrderNo:     order.MerchantOrderNo,
		OrderNo:             order.OrderNo,
		OrderAmount:         order.OrderAmount,
		ActualAmount:         order.ActualAmount,
		TransferHandlingFee: order.TransferHandlingFee,
		NotifyUrl:           order.NotifyUrl,
		OrderTime:           order.CreatedAt.Format("20060102150405000"),
		PayOrderTime:        order.TransAt.Time().Format("20060102150405000"),
		Status:              order.Status,
	}, nil
}

func (l *PayCallBackTranactionLogic) updateOrderAndBalance(db *gorm.DB, req *transactionclient.PayCallBackRequest, order *types.OrderX) (err error) {

	var merchantBalanceRecord types.MerchantBalanceRecord

	// 回調成功
	if req.OrderStatus == "20" {
		// 回调金额 才是实际收款金额
		order.ActualAmount = req.OrderAmount
		// (更改为实际收款金额) 交易手續費總額 = 訂單金額 / 100 * 費率 + 手續費
		order.TransferHandlingFee = utils.FloatAdd(utils.FloatMul(utils.FloatDiv(order.ActualAmount, 100), order.Fee), order.HandlingFee)
		// (更改为实际收款金额) 計算實際交易金額 = 訂單金額 - 手續費
		order.TransferAmount = order.ActualAmount - order.TransferHandlingFee

		// 異動錢包
		if merchantBalanceRecord, err = l.UpdateBalance(db, types.UpdateBalance{
			MerchantCode:    order.MerchantCode,
			CurrencyCode:    order.CurrencyCode,
			OrderNo:         order.OrderNo,
			OrderType:       order.Type,
			ChannelCode:     order.ChannelCode,
			PayTypeCode:     order.PayTypeCode,
			PayTypeCodeNum:  order.PayTypeCodeNum,
			TransactionType: "1",
			BalanceType:     order.BalanceType,
			TransferAmount:  order.TransferAmount,
			Comment:         order.Memo,
			CreatedBy:       order.MerchantCode,
		}); err != nil {
			return
		}

		order.BeforeBalance = merchantBalanceRecord.BeforeBalance
		order.Balance = merchantBalanceRecord.AfterBalance
		order.TransAt = types.JsonTime{}.New()
	}

	order.ChannelOrderNo = req.ChannelOrderNo
	order.Status = req.OrderStatus
	order.CallBackStatus = "1"

	// 編輯訂單
	if err = db.Table("tx_orders").Updates(&order).Error; err != nil {
		return
	}

	return
}

// UpdateBalance TransferAmount
func (l *PayCallBackTranactionLogic) UpdateBalance(db *gorm.DB, updateBalance types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {

	var beforeBalance float64
	var afterBalance float64

	// 1. 取得 商戶餘額表
	var merchantBalance types.MerchantBalance
	if err = db.Table("mc_merchant_balances").
		Where("merchant_code = ? AND currency_code = ? AND balance_type = ?", updateBalance.MerchantCode, updateBalance.CurrencyCode, updateBalance.BalanceType).
		Take(&merchantBalance).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 2. 計算
	var selectBalance string
	if utils.FloatAdd(merchantBalance.Balance, updateBalance.TransferAmount) < 0 {
		logx.Errorf("商户:%s，余额类型:%s，余额:%s，交易金额:%s", merchantBalance.MerchantCode, merchantBalance.BalanceType, fmt.Sprintf("%f", merchantBalance.Balance), fmt.Sprintf("%f", updateBalance.TransferAmount))
		return merchantBalanceRecord, errorz.New(response.MERCHANT_INSUFFICIENT_DF_BALANCE)
	}
	selectBalance = "balance"
	beforeBalance = merchantBalance.Balance
	afterBalance = utils.FloatAdd(beforeBalance, updateBalance.TransferAmount)
	merchantBalance.Balance = afterBalance

	// 3. 變更 商戶餘額
	if err = db.Table("mc_merchant_balances").Select(selectBalance).Updates(types.MerchantBalanceX{
		MerchantBalance: merchantBalance,
	}).Error; err != nil {
		logx.Error(err.Error())
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 4. 新增 餘額紀錄
	merchantBalanceRecord = types.MerchantBalanceRecord{
		MerchantBalanceId: merchantBalance.ID,
		MerchantCode:      merchantBalance.MerchantCode,
		CurrencyCode:      merchantBalance.CurrencyCode,
		OrderNo:           updateBalance.OrderNo,
		OrderType:         updateBalance.OrderType,
		ChannelCode:       updateBalance.ChannelCode,
		PayTypeCode:       updateBalance.PayTypeCode,
		PayTypeCodeNum:    updateBalance.PayTypeCodeNum,
		TransactionType:   updateBalance.TransactionType,
		BalanceType:       updateBalance.BalanceType,
		BeforeBalance:     beforeBalance,
		TransferAmount:    updateBalance.TransferAmount,
		AfterBalance:      afterBalance,
		Comment:           updateBalance.Comment,
		CreatedBy:         updateBalance.CreatedBy,
	}

	if err = db.Table("mc_merchant_balance_records").Create(&types.MerchantBalanceRecordX{
		MerchantBalanceRecord: merchantBalanceRecord,
	}).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	return
}
