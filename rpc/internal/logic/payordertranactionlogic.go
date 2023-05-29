package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
)

type PayOrderTranactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayOrderTranactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderTranactionLogic {
	return &PayOrderTranactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayOrderTranactionLogic) PayOrderTranaction(in *transactionclient.PayOrderRequest) (resp *transactionclient.PayOrderResponse, err error) {

	var payOrderReq = in.PayOrder
	var correspondMerChnRate = in.Rate
	var isMerchantCallback string

	orderAmount, _ := strconv.ParseFloat(payOrderReq.OrderAmount, 64)

	if payOrderReq.NotifyUrl != "" {
		isMerchantCallback = constants.MERCHANT_CALL_BACK_NO
	} else {
		isMerchantCallback = constants.MERCHANT_CALL_BACK_DONT_USE
	}

	// 初始化訂單
	order := &types.Order{
		Type:                constants.ORDER_TYPE_ZF,
		MerchantCode:        payOrderReq.MerchantId,
		OrderNo:             in.OrderNo,
		MerchantOrderNo:     payOrderReq.OrderNo,
		ChannelOrderNo:      in.ChannelOrderNo,
		FrozenAmount:        0,
		Status:              constants.PROCESSING,
		IsLock:              constants.IS_LOCK_NO,
		ChannelCode:         correspondMerChnRate.ChannelCode,
		ChannelPayTypesCode: correspondMerChnRate.ChannelPayTypesCode,
		PayTypeCode:         correspondMerChnRate.PayTypeCode,
		CurrencyCode:        correspondMerChnRate.CurrencyCode,
		MerchantBankNo:      payOrderReq.BankCode,
		MerchantAccountName: payOrderReq.UserId,
		OrderAmount:         orderAmount,
		Source:              constants.API,
		CallBackStatus:      constants.CALL_BACK_STATUS_PROCESSING,
		IsMerchantCallback:  isMerchantCallback,
		NotifyUrl:           payOrderReq.NotifyUrl,
		PageUrl:             payOrderReq.PageUrl,
		PersonProcessStatus: constants.PERSON_PROCESS_STATUS_NO_ROCESSING,
		IsCalculateProfit:   constants.IS_CALCULATE_PROFIT_NO,
		IsTest:              constants.IS_TEST_NO,
		CreatedBy:           payOrderReq.MerchantId,
		UpdatedBy:           payOrderReq.MerchantId,
	}

	// 取得餘額類型
	if order.BalanceType, err = merchantbalanceservice.GetBalanceType(l.svcCtx.MyDB, order.ChannelCode, order.Type); err != nil {
		return
	}

	/****     交易開始      ****/
	txDB := l.svcCtx.MyDB.Begin()

	order.Fee = correspondMerChnRate.Fee
	order.HandlingFee = correspondMerChnRate.HandlingFee
	// 交易手續費總額 = 訂單金額 / 100 * 費率 + 手續費
	order.TransferHandlingFee = utils.FloatAdd(utils.FloatMul(utils.FloatDiv(order.OrderAmount, 100), order.Fee), order.HandlingFee)
	// 計算實際交易金額 = 訂單金額 - 手續費
	order.TransferAmount = order.OrderAmount - order.TransferHandlingFee

	if err = txDB.Table("tx_orders").Create(&types.OrderX{
		Order: *order,
	}).Error; err != nil {
		txDB.Rollback()
		return &transactionclient.PayOrderResponse{
			Code:       response.DATABASE_FAILURE,
			Message:    "数据库错误 tx_orders Create",
			PayOrderNo: order.OrderNo,
		}, nil
	}

	if err = txDB.Commit().Error; err != nil {
		txDB.Rollback()
		logx.WithContext(l.ctx).Errorf("支付提单失败，商户号: %s, 订单号: %s, err : %s", order.MerchantCode, order.OrderNo, err.Error())
		return &transactionclient.PayOrderResponse{
			Code:       response.DATABASE_FAILURE,
			Message:    "Commit 数据库错误",
			PayOrderNo: order.OrderNo,
		}, nil
	}
	/****     交易結束      ****/

	// 新單新增訂單歷程 (不抱錯) TODO: 異步??
	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     order.OrderNo,
			Action:      "PLACE_ORDER",
			UserAccount: payOrderReq.MerchantId,
			Comment:     "",
		},
	}).Error; err4 != nil {
		logx.WithContext(l.ctx).Errorf("紀錄訂單歷程出錯:%s", err4.Error())
	}

	return &transactionclient.PayOrderResponse{
		Code:       response.API_SUCCESS,
		Message:    "操作成功",
		PayOrderNo: order.OrderNo,
	}, nil
}
