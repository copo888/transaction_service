package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/service/orderfeeprofitservice"
	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transaction"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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

func (l *PayOrderTranactionLogic) PayOrderTranaction(in *transaction.PayOrderRequest) (resp *transaction.PayOrderResponse, err error) {

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
		PayTypeCodeNum:      correspondMerChnRate.PayTypeCodeNum,
		PayTypeNum:          correspondMerChnRate.PayTypeCode + correspondMerChnRate.PayTypeCodeNum,
		CurrencyCode:        correspondMerChnRate.CurrencyCode,
		MerchantBankNo:      payOrderReq.BankCode,
		MerchantBankName:    payOrderReq.UserId,
		OrderAmount:         orderAmount,
		Source:              constants.API,
		CallBackStatus:      constants.CALL_BACK_STATUS_PROCESSING,
		IsMerchantCallback:  isMerchantCallback,
		NotifyUrl:           payOrderReq.NotifyUrl,
		PageUrl:             payOrderReq.PageUrl,
		PersonProcessStatus: constants.PERSON_PROCESS_STATUS_NO_ROCESSING,
		CreatedBy:           payOrderReq.MerchantId,
		UpdatedBy:           payOrderReq.MerchantId,
	}

	// 取得餘額類型
	if order.BalanceType, err = merchantbalanceservice.GetBalanceType(l.svcCtx.MyDB, order.ChannelCode, order.Type); err != nil {
		return
	}

	if err = l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

		order.Fee = correspondMerChnRate.Fee
		order.HandlingFee = correspondMerChnRate.HandlingFee
		// 交易手續費總額 = 訂單金額 / 100 * 費率 + 手續費
		order.TransferHandlingFee = utils.FloatAdd(utils.FloatMul(utils.FloatDiv(order.OrderAmount, 100), order.Fee), order.HandlingFee)
		// 計算實際交易金額 = 訂單金額 - 手續費
		order.TransferAmount = order.OrderAmount - order.TransferHandlingFee

		if err = db.Table("tx_orders").Create(&types.OrderX{
			Order: *order,
		}).Error; err != nil {
			return
		}

		return nil
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
		OrderAmount:         order.OrderAmount,
	}); err4 != nil {
		logx.Error("計算利潤出錯:%s", err4.Error())
	}

	// 新單新增訂單歷程 (不抱錯) TODO: 異步??
	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     order.OrderNo,
			Action:      "PLACE_ORDER",
			UserAccount: payOrderReq.MerchantId,
			Comment:     "",
		},
	}).Error; err4 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err4.Error())
	}

	return &transaction.PayOrderResponse{
		PayOrderNo: order.OrderNo,
	}, nil
}
