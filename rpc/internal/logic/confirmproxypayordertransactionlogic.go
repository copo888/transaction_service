package logic

import (
	"context"
	"fmt"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"github.com/gioco-play/easy-i18n/i18n"
	"github.com/neccoys/go-zero-extension/redislock"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmProxyPayOrderTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmProxyPayOrderTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmProxyPayOrderTransactionLogic {
	return &ConfirmProxyPayOrderTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfirmProxyPayOrderTransactionLogic) ConfirmProxyPayOrderTransaction(in *transactionclient.ConfirmProxyPayOrderRequest) (*transactionclient.ConfirmProxyPayOrderResponse, error) {
	var order *types.OrderX

	myDB := l.svcCtx.MyDB

	if err := myDB.Table("tx_orders").
		Where("order_no = ?", in.OrderNo).Take(&order).Error; err != nil {
		return &transactionclient.ConfirmProxyPayOrderResponse{
			Code:    response.ORDER_NUMBER_NOT_EXIST,
			Message: "订单号不存在 ",
		}, nil
	}

	// 代付單才能
	if order.Type != constants.ORDER_TYPE_DF {
		return &transactionclient.ConfirmProxyPayOrderResponse{
			Code:    response.ORDER_TYPE_IS_WRONG,
			Message: "订单类型错误 ",
		}, nil
	}

	// 鎖定單 不可操作
	if order.IsLock == "1" {
		return &transactionclient.ConfirmProxyPayOrderResponse{
			Code:    response.ORDER_IS_STATUS_IS_LOCK,
			Message: "订单号已锁定",
		}, nil
	}

	// 失敗單才能轉成功單
	if order.Status != constants.FAIL {
		return &transactionclient.ConfirmProxyPayOrderResponse{
			Code:    response.ORDER_STATUS_WRONG,
			Message: "订单状态錯誤",
		}, nil
	}

	// 确认是否有设置费率
	var merchantChannelRate *types.MerchantChannelRate
	if err := l.svcCtx.MyDB.Table("mc_merchant_channel_rate").
		Where("merchant_code = ? AND channel_pay_types_code = ?", order.MerchantCode, order.ChannelPayTypesCode).
		Take(&merchantChannelRate).Error; err != nil {
		return &transactionclient.ConfirmProxyPayOrderResponse{
			Code:    response.RATE_NOT_CONFIGURED,
			Message: "未配置商户渠道费率",
		}, nil
	}

	updateBalance := &types.UpdateBalance{
		MerchantCode:    order.MerchantCode,
		CurrencyCode:    order.CurrencyCode,
		OrderNo:         order.OrderNo,
		MerchantOrderNo: order.MerchantOrderNo,
		OrderType:       order.Type,
		PayTypeCode:     order.PayTypeCode,
		TransferAmount:  order.TransferAmount,
		TransactionType: constants.TRANSACTION_TYPE_PROXY_PAY,
		BalanceType:     order.BalanceType,
		CreatedBy:       order.MerchantCode,
		ChannelCode:     order.ChannelCode,
		Comment:         in.Comment,
		MerPtBalanceId:  merchantChannelRate.MerchantPtBalanceId,
	}

	redisKey := fmt.Sprintf("%s-%s", updateBalance.MerchantCode, updateBalance.CurrencyCode)
	redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(8)
	if isOK, redisErr := redisLock.TryLockTimeout(8); isOK {
		defer redisLock.Release()

		/****     交易開始      ****/
		txDB := myDB.Begin()

		if order.BalanceType == constants.DF_BALANCE {
			//异动子钱包
			if merchantChannelRate.MerchantPtBalanceId > 0 {
				if _, err := merchantbalanceservice.UpdateDF_Pt_Balance_Debit(l.ctx, txDB, updateBalance); err != nil {
					logx.WithContext(l.ctx).Errorf("商户:%s，幣別: %s，更新子錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, order.CurrencyCode, err.Error(), updateBalance)
					return &transactionclient.ConfirmProxyPayOrderResponse{
						Code:    response.SYSTEM_ERROR,
						Message: "更新子钱包错误",
					}, nil
				}
			}

			if _, err := merchantbalanceservice.UpdateDFBalance_Debit(l.ctx, txDB, updateBalance); err != nil {
				logx.WithContext(l.ctx).Errorf("商户:%s，单号:%s，更新錢包紀錄錯誤:%s", order.MerchantCode, order.OrderNo, err.Error())
				txDB.Rollback()
				return &transactionclient.ConfirmProxyPayOrderResponse{
					Code:    response.SYSTEM_ERROR,
					Message: "更新钱包错误",
				}, nil
			}
		} else if order.BalanceType == constants.XF_BALANCE {

			if merchantChannelRate.MerchantPtBalanceId > 0 {
				if _, err := merchantbalanceservice.UpdateXF_Pt_Balance_Debit(l.ctx, txDB, updateBalance); err != nil {
					logx.WithContext(l.ctx).Errorf("商户:%s，幣別: %s，更新子錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, order.CurrencyCode, err.Error(), updateBalance)
					return &transactionclient.ConfirmProxyPayOrderResponse{
						Code:    response.SYSTEM_ERROR,
						Message: "更新子钱包错误",
					}, nil
				}
			}

			if _, err := merchantbalanceservice.UpdateXFBalance_Debit(l.ctx, txDB, updateBalance); err != nil {
				logx.WithContext(l.ctx).Errorf("商户:%s，单号:%s，更新錢包紀錄錯誤:%s", order.MerchantCode, order.OrderNo, err.Error())
				txDB.Rollback()
				return &transactionclient.ConfirmProxyPayOrderResponse{
					Code:    response.SYSTEM_ERROR,
					Message: "更新钱包错误",
				}, nil
			}
		} else {
			logx.WithContext(l.ctx).Errorf("商户:%s，单号:%s，錢包類型錯誤:%s", order.MerchantCode, order.BalanceType)
			txDB.Rollback()
			return &transactionclient.ConfirmProxyPayOrderResponse{
				Code:    response.SYSTEM_ERROR,
				Message: "錢包類型錯誤",
			}, nil
		}

		// 編輯訂單
		order.Status = constants.SUCCESS
		order.TransAt = types.JsonTime{}.New()
		order.Memo = in.Comment + " \n" + order.Memo
		if err := txDB.Table("tx_orders").Updates(&order).Error; err != nil {
			logx.WithContext(l.ctx).Errorf("商户:%s，单号:%s，编辑订单错误:%s", order.MerchantCode, order.OrderNo, err.Error())
			txDB.Rollback()
			return &transactionclient.ConfirmProxyPayOrderResponse{
				Code:    response.SYSTEM_ERROR,
				Message: "錢包類型錯誤",
			}, nil
		}

		if err := txDB.Commit().Error; err != nil {
			txDB.Rollback()
			logx.Errorf("支付確認收款Commit失败，商户号: %s, 订单号: %s, err : %s", order.MerchantCode, order.OrderNo, err.Error())
			return &transactionclient.ConfirmProxyPayOrderResponse{
				Code:    response.DATABASE_FAILURE,
				Message: "资料库错误 Commit失败",
			}, nil
		}
		/****     交易結束      ****/
	} else {
		logx.WithContext(l.ctx).Errorf("商户钱包处理中，Err:%s。 %s", redisErr.Error(), redisKey)
		return &transactionclient.ConfirmProxyPayOrderResponse{
			Code:    response.BALANCE_PROCESSING,
			Message: i18n.Sprintf(response.BALANCE_PROCESSING),
		}, nil
	}

	// 新單新增訂單歷程 (不抱錯)
	if err := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     order.OrderNo,
			Action:      constants.ACTION_SUCCESS,
			UserAccount: order.MerchantCode,
			Comment:     in.Comment,
		},
	}).Error; err != nil {
		logx.WithContext(l.ctx).Error("紀錄訂單歷程出錯:%s", err.Error())
	}

	return &transactionclient.ConfirmProxyPayOrderResponse{
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}, nil
}
