package logic

import (
	"context"
	"errors"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/transaction"

	"github.com/zeromicro/go-zero/core/logx"
)

type CryptoPayOrderTranactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCryptoPayOrderTranactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CryptoPayOrderTranactionLogic {
	return &CryptoPayOrderTranactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CryptoPayOrderTranactionLogic) CryptoPayOrderTranaction(in *transaction.CryptoPayOrderRequest) (resp *transaction.CryptoPayOrderResponse, err error) {

	var payOrderReq = in.CryptoPayOrder
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
		MerchantBankNo:      "",
		MerchantBankName:    "",
		OrderAmount:         orderAmount,
		Source:              constants.API,
		CallBackStatus:      constants.CALL_BACK_STATUS_PROCESSING,
		IsMerchantCallback:  isMerchantCallback,
		NotifyUrl:           payOrderReq.NotifyUrl,
		PageUrl:             "",
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

	// TODO: 綁定地址
	// 1. 取得可綁定地址
	var walletAddress types.WalletAddress
	if err = txDB.Table("ch_wallet_address").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("channel_pay_types_code = ? AND status = ? AND order_no = '' ", in.Rate.ChannelPayTypesCode, constants.WALLET_ADDRESS_STATUS_IDLE).
		Take(&walletAddress).Error; errors.Is(err, gorm.ErrRecordNotFound) { // 無閒置可用地址
		txDB.Rollback()
		logx.WithContext(l.ctx).Errorf("虛擬貨幣支付提单失败,無閒置可用地址，商户号: %s, 订单号: %s, err : %s", order.MerchantCode, order.OrderNo, err.Error())
		return &transactionclient.CryptoPayOrderResponse{
			Code:       response.NO_AVAILABLE_CHANNEL_ERROR,
			Message:    "無可用地址",
			PayOrderNo: order.OrderNo,
		}, nil
	} else if err != nil {
		txDB.Rollback()
		logx.WithContext(l.ctx).Errorf("虛擬貨幣支付提单失败,数据库错误，商户号: %s, 订单号: %s, err : %s", order.MerchantCode, order.OrderNo, err.Error())
		return &transactionclient.CryptoPayOrderResponse{
			Code:       response.DATABASE_FAILURE,
			Message:    "数据库错误 ch_wallet_address select",
			PayOrderNo: order.OrderNo,
		}, nil
	}

	order.ChannelBankAccount = walletAddress.Address // 訂單綁定地址
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
		logx.WithContext(l.ctx).Errorf("虛擬貨幣支付提单失败,Create error，商户号: %s, 订单号: %s, err : %s", order.MerchantCode, order.OrderNo, err.Error())
		return &transactionclient.CryptoPayOrderResponse{
			Code:       response.DATABASE_FAILURE,
			Message:    "数据库错误 tx_orders Create error",
			PayOrderNo: order.OrderNo,
		}, nil
	}

	// 抵制改為綁定狀態
	if err = l.svcCtx.MyDB.Table("ch_wallet_address").Where("id = ?", walletAddress.ID).
		Updates(map[string]interface{}{
			"order_no": order.OrderNo,
			"status":  constants.WALLET_ADDRESS_STATUS_USE,
			"usage_at":  types.JsonTime{}.New(),
		}).Error; err != nil {
		logx.WithContext(l.ctx).Error("編輯錢包地址失敗: %s", err.Error())
		return &transactionclient.CryptoPayOrderResponse{
			Code:       response.DATABASE_FAILURE,
			Message:    "数据库错误, ch_wallet_address error",
			PayOrderNo: order.OrderNo,
		}, nil
	}

	if err = txDB.Commit().Error; err != nil {
		txDB.Rollback()
		logx.WithContext(l.ctx).Errorf("虛擬貨幣支付提单失败，商户号: %s, 订单号: %s, err : %s", order.MerchantCode, order.OrderNo, err.Error())
		return &transactionclient.CryptoPayOrderResponse{
			Code:       response.DATABASE_FAILURE,
			Message:    "Commit 数据库错误",
			PayOrderNo: order.OrderNo,
		}, nil
	}
	/****     交易結束      ****/

	// 新單新增訂單歷程 (不抱錯)
	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     order.OrderNo,
			Action:      "PLACE_ORDER",
			UserAccount: payOrderReq.MerchantId,
			Comment:     "",
		},
	}).Error; err4 != nil {
		logx.WithContext(l.ctx).Errorf("虛擬貨幣紀錄訂單歷程出錯:%s", err4.Error())
	}

	return &transactionclient.CryptoPayOrderResponse{
		Code:       response.API_SUCCESS,
		Message:    "操作成功",
		PayOrderNo: order.OrderNo,
	}, nil
}
