package logic

import (
	"context"
	"fmt"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/service/orderfeeprofitservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"github.com/neccoys/go-zero-extension/redislock"
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

	if in.OrderStatus == "20" {
		return l.PayCallBackTranactionForSuccess(l.ctx, in)
	} else if in.OrderStatus == "30" {
		return l.PayCallBackTranactionForFailure(l.ctx, in)
	}
	return &transactionclient.PayCallBackResponse{
		Code:    response.ORDER_STATUS_WRONG,
		Message: fmt.Sprintf("訂單:%s, 回調狀態:%s. 異常", in.PayOrderNo, in.OrderStatus),
	}, nil

}

func (l *PayCallBackTranactionLogic) PayCallBackTranactionForSuccess(ctx context.Context, in *transactionclient.PayCallBackRequest) (resp *transactionclient.PayCallBackResponse, err error) {
	var order *types.OrderX

	myDB := l.svcCtx.MyDB

	// 取得tx_order
	if err = myDB.Table("tx_orders").
		Where("order_no = ?", in.PayOrderNo).Take(&order).Error; err != nil {
		return &transactionclient.PayCallBackResponse{
			Code:    response.ORDER_NUMBER_NOT_EXIST,
			Message: "平台订单号不存在",
		}, nil
	}

	// 这里谨用于确认是否有设置费率
	var merchantChannelRate *types.MerchantChannelRate
	if err := myDB.Table("mc_merchant_channel_rate").
		Where("merchant_code = ? AND channel_pay_types_code = ?", order.MerchantCode, order.ChannelPayTypesCode).
		Take(&merchantChannelRate).Error; err != nil {
		return &transactionclient.PayCallBackResponse{
			Code:    response.RATE_NOT_CONFIGURED,
			Message: "未配置商户渠道费率",
		}, nil
	}

	// 處理中的且非鎖定訂單 才能回調
	if order.Status != "1" || order.IsLock == "1" {
		return &transactionclient.PayCallBackResponse{
			Code:    response.TRANSACTION_FAILURE,
			Message: "交易失败 订单号已锁定 或 订单状态非处理中",
		}, nil
	}

	// 下單金額及實付金額差異風控 (差異超過5% 且 超過1元)
	limit := utils.FloatMul(order.OrderAmount, 0.05)
	diff := math.Abs(utils.FloatSub(order.OrderAmount, in.OrderAmount))
	if diff > limit && diff > 1 {
		return &transactionclient.PayCallBackResponse{
			Code:    response.ORDER_AMOUNT_ERROR,
			Message: "商户下单金额和回調金額不符" + fmt.Sprintf("(orderAmount/payAmount): %f/%f", order.OrderAmount, in.OrderAmount),
		}, nil
	}

	redisKey := fmt.Sprintf("%s-%s", order.MerchantCode, order.CurrencyCode)
	redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(5)
	if isOK, _ := redisLock.TryLockTimeout(5); isOK {
		defer redisLock.Release()
		/****     交易開始      ****/
		txDB := l.svcCtx.MyDB.Begin()
		// 編輯訂單 異動錢包和餘額
		if err = l.updateOrderAndBalance(txDB, in, order); err != nil {
			txDB.Rollback()
			return &transactionclient.PayCallBackResponse{
				Code:    response.SYSTEM_ERROR,
				Message: "钱包异动失败",
			}, nil
		}

		if err = txDB.Commit().Error; err != nil {
			txDB.Rollback()
			logx.Errorf("支付回調失败，商户号: %s, 订单号: %s, err : %s", order.MerchantCode, order.OrderNo, err.Error())
			return &transactionclient.PayCallBackResponse{
				Code:    response.DATABASE_FAILURE,
				Message: "资料库错误 Commit失败",
			}, nil
		}
		/****     交易結束      ****/
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
		logx.WithContext(ctx).Error("計算利潤出錯:%s", err4.Error())
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
		logx.WithContext(ctx).Error("紀錄訂單歷程出錯:%s", err4.Error())
	}

	return &transactionclient.PayCallBackResponse{
		Code:                response.API_SUCCESS,
		Message:             "操作成功",
		MerchantCode:        order.MerchantCode,
		MerchantOrderNo:     order.MerchantOrderNo,
		OrderNo:             order.OrderNo,
		OrderAmount:         order.OrderAmount,
		ActualAmount:        order.ActualAmount,
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

		updateBalance := types.UpdateBalance{
			MerchantCode:    order.MerchantCode,
			CurrencyCode:    order.CurrencyCode,
			OrderNo:         order.OrderNo,
			MerchantOrderNo: order.MerchantOrderNo,
			OrderType:       order.Type,
			ChannelCode:     order.ChannelCode,
			PayTypeCode:     order.PayTypeCode,
			TransactionType: constants.TRANSACTION_TYPE_PAY,
			BalanceType:     order.BalanceType,
			TransferAmount:  order.TransferAmount,
			Comment:         order.Memo,
			CreatedBy:       order.MerchantCode,
		}

		// 異動錢包
		if merchantBalanceRecord, err = merchantbalanceservice.UpdateBalanceForZF(db, l.ctx, l.svcCtx.RedisClient, updateBalance); err != nil {
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

func (l *PayCallBackTranactionLogic) PayCallBackTranactionForFailure(ctx context.Context, in *transactionclient.PayCallBackRequest) (resp *transactionclient.PayCallBackResponse, err error) {
	var order *types.OrderX

	myDB := l.svcCtx.MyDB

	if err = myDB.Table("tx_orders").
		Where("order_no = ?", in.PayOrderNo).Take(&order).Error; err != nil {
		return &transactionclient.PayCallBackResponse{
			Code:    response.ORDER_NUMBER_NOT_EXIST,
			Message: "平台订单号不存在",
		}, nil
	}

	// 處理中的且非鎖定訂單 才能回調
	if order.Status != "1" || order.IsLock == "1" {
		return &transactionclient.PayCallBackResponse{
			Code:    response.TRANSACTION_FAILURE,
			Message: "交易失败 订单号已锁定 或 订单状态非处理中",
		}, nil
	}

	redisKey := fmt.Sprintf("%s-%s", order.MerchantCode, order.CurrencyCode)
	redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(5)
	if isOK, _ := redisLock.TryLockTimeout(5); isOK {
		defer redisLock.Release()
		/****     交易開始      ****/
		txDB := myDB.Begin()

		order.TransAt = types.JsonTime{}.New()
		order.ChannelOrderNo = in.ChannelOrderNo
		order.Status = in.OrderStatus
		order.CallBackStatus = "1"
		order.TransferAmount = 0
		order.TransferHandlingFee = 0

		// 編輯訂單
		if err = txDB.Select(
			"trans_at",
			"channel_order_no",
			"status",
			"callback_status",
			"transfer_amount",
			"transfer_handling_fee").
			Table("tx_orders").Updates(&order).Error; err != nil {
			txDB.Rollback()
			return &transactionclient.PayCallBackResponse{
				Code:    response.SYSTEM_ERROR,
				Message: "訂單失敗狀態 編輯失败",
			}, nil
		}

		if err = txDB.Commit().Error; err != nil {
			txDB.Rollback()
			logx.Errorf("支付回調失败，商户号: %s, 订单号: %s, err : %s", order.MerchantCode, order.OrderNo, err.Error())
			return &transactionclient.PayCallBackResponse{
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
			Action:      constants.ACTION_FAILURE,
			UserAccount: order.MerchantCode,
			Comment:     "",
		},
	}).Error; err4 != nil {
		logx.WithContext(ctx).Error("紀錄訂單歷程出錯:%s", err4.Error())
	}

	return &transactionclient.PayCallBackResponse{
		Code:                response.API_SUCCESS,
		Message:             "操作成功",
		MerchantCode:        order.MerchantCode,
		MerchantOrderNo:     order.MerchantOrderNo,
		OrderNo:             order.OrderNo,
		OrderAmount:         order.OrderAmount,
		ActualAmount:        order.ActualAmount,
		TransferHandlingFee: order.TransferHandlingFee,
		NotifyUrl:           order.NotifyUrl,
		OrderTime:           order.CreatedAt.Format("20060102150405000"),
		PayOrderTime:        order.TransAt.Time().Format("20060102150405000"),
		Status:              order.Status,
	}, nil
}
