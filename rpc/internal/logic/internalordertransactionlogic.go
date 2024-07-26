package logic

import (
	"context"
	"fmt"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/response"
	"github.com/gioco-play/easy-i18n/i18n"
	"github.com/neccoys/go-zero-extension/redislock"

	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"gorm.io/gorm"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type InternalOrderTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInternalOrderTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalOrderTransactionLogic {
	return &InternalOrderTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InternalOrderTransactionLogic) InternalOrderTransaction(in *transactionclient.InternalOrderRequest) (resp *transactionclient.InternalOrderResponse, err error) {

	var internalOrderReq = in.InternalOrder
	//var merchantOrderRateListView = in.MerchantOrderRateListView

	//// 交易手續費總額 = 訂單金額 / 100 * 費率
	//transferHandling := utils.FloatMulC(utils.FloatDivC(internalOrderReq.OrderAmount, 100), merchantOrderRateListView.MerFee)
	//
	//// 計算實際交易金額 = 訂單金額 - 手續費
	//transferAmount := utils.FloatSubC(internalOrderReq.OrderAmount, transferHandling)

	//初始化订单
	txOrder := &types.Order{
		OrderNo:                 model.GenerateOrderNo("NC"),
		MerchantCode:            internalOrderReq.MerchantCode,
		Type:                    constants.ORDER_TYPE_NC,
		Status:                  constants.PROCESSING,
		Source:                  constants.UI,
		IsMerchantCallback:      constants.IS_MERCHANT_CALLBACK_NOT_NEED,
		IsCalculateProfit:       constants.IS_CALCULATE_PROFIT_NO,
		IsTest:                  constants.IS_TEST_NO, //是否測試單
		PersonProcessStatus:     constants.PERSON_PROCESS_STATUS_NO_ROCESSING,
		InternalChargeOrderPath: internalOrderReq.Imgurl,
		//BalanceType:             "DFB",
		OrderAmount: internalOrderReq.OrderAmount,
		//TransferHandlingFee:     transferHandling,
		//TransferAmount:          transferAmount,
		CreatedBy:            internalOrderReq.UserAccount,
		UpdatedBy:            internalOrderReq.UserAccount,
		IsLock:               "0", //是否锁定状态 (0=否;1=是) 预设否
		CurrencyCode:         internalOrderReq.CurrencyCode,
		MerchantAccountName:  internalOrderReq.MerchantAccountName,
		MerchantBankAccount:  internalOrderReq.MerchantBankAccount,
		MerchantBankCity:     internalOrderReq.MerchantBankCity,
		MerchantBankProvince: internalOrderReq.MerchantBankProvince,
		MerchantBankNo:       internalOrderReq.MerchantBankNo,
		MerchantBankName:     internalOrderReq.MerchantBankName,
		ChannelBankName:      internalOrderReq.ChannelBankName,
		ChannelAccountName:   internalOrderReq.ChannelAccountName,
		ChannelBankAccount:   internalOrderReq.ChannelBankAccount,
		ChannelBankNo:        internalOrderReq.ChannelBankNo,
		//ChannelCode:             merchantOrderRateListView.ChannelCode,
		//ChannelPayTypesCode:     merchantOrderRateListView.ChannelPayTypesCode,
		//PayTypeCode:             merchantOrderRateListView.PayTypeCode,
		//Fee:                     merchantOrderRateListView.MerFee,
		//HandlingFee:             merchantOrderRateListView.MerHandlingFee,
	}

	redisKey := fmt.Sprintf("%s-%s", txOrder.MerchantCode, txOrder.CurrencyCode)
	redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(8)
	if isOK, redisErr := redisLock.TryLockTimeout(8); isOK {
		defer redisLock.Release()

		if err = l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {
			txOrder.MerchantOrderNo = "COPO_" + txOrder.OrderNo

			if err = db.Table("tx_orders").Create(&types.OrderX{
				Order: *txOrder,
			}).Error; err != nil {
				logx.Errorf("新增内充提单失败，商户号: %s, 订单号: %s, err : %s", txOrder.MerchantCode, txOrder.OrderNo, err.Error())
				return
			}

			//// 計算利潤 ,修改内充功能，利润改在审核才计算
			//if err = orderfeeprofitservice.CalculateOrderProfit(db, types.CalculateProfit{
			//	MerchantCode:        txOrder.MerchantCode,
			//	OrderNo:             txOrder.OrderNo,
			//	Type:                txOrder.Type,
			//	CurrencyCode:        txOrder.CurrencyCode,
			//	BalanceType:         txOrder.BalanceType,
			//	ChannelCode:         txOrder.ChannelCode,
			//	ChannelPayTypesCode: txOrder.ChannelPayTypesCode,
			//	OrderAmount:         txOrder.OrderAmount,
			//}); err != nil {
			//	logx.Error("計算利潤出錯:%s", err.Error())
			//	return err
			//}

			return nil
		}); err != nil {
			return &transactionclient.InternalOrderResponse{
				Code:    response.DATABASE_FAILURE,
				Message: "数据库错误 tx_orders Create internal charge，err : " + err.Error(),
			}, nil
		}
	} else {
		logx.WithContext(l.ctx).Errorf("商户钱包处理中，Err:%s。 %s", redisErr.Error(), redisKey)
		return &transactionclient.InternalOrderResponse{
			Code:    response.BALANCE_PROCESSING,
			Message: i18n.Sprintf(response.BALANCE_PROCESSING),
		}, nil
	}

	// 新單新增訂單歷程 (不抱錯) TODO: 異步??
	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     txOrder.OrderNo,
			Action:      "PLACE_ORDER",
			UserAccount: internalOrderReq.UserAccount,
			Comment:     "",
		},
	}).Error; err4 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err4.Error())
	}
	return &transactionclient.InternalOrderResponse{
		OrderNo: txOrder.OrderNo,
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}, nil
}
