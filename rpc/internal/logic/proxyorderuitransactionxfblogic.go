package logic

import (
	"context"
	"fmt"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/redislock"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/service/orderfeeprofitservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"github.com/gioco-play/easy-i18n/i18n"
	"gorm.io/gorm"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProxyOrderUITransactionXFBLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProxyOrderUITransactionXFBLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyOrderUITransactionXFBLogic {
	return &ProxyOrderUITransactionXFBLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProxyOrderUITransactionXFBLogic) ProxyOrderUITransaction_XFB(in *transactionclient.ProxyOrderUIRequest) (resp *transactionclient.ProxyOrderUIResponse, err error) {
	req := in.ProxyOrderUI
	rate := in.MerchantOrderRateListView
	merchantBalanceRecord := &types.MerchantBalanceRecord{}

	var transferHandlingFee float64
	if rate.IsRate == "1" { // 是否算費率，0:否 1:是
		//  交易手續費總額 = 訂單金額 / 100 * 費率 + 手續費
		transferHandlingFee =
			utils.FloatAddCWithTrancated(utils.FloatMulC(utils.FloatDivC(req.OrderAmount, 100, req.CurrencyCode), rate.MerFee, req.CurrencyCode), rate.MerHandlingFee, req.CurrencyCode)
	} else {
		//  交易手續費總額 = 訂單金額 / 100 * 費率 + 手續費
		transferHandlingFee =
			utils.FloatAddCWithTrancated(utils.FloatMulC(utils.FloatDivC(req.OrderAmount, 100, req.CurrencyCode), 0, req.CurrencyCode), rate.MerHandlingFee, req.CurrencyCode)
	}
	txOrder := &types.Order{
		MerchantCode:         req.MerchantCode,
		CreatedBy:            req.UserAccount,
		MerchantOrderNo:      "COPO_" + req.OrderNo,
		OrderNo:              req.OrderNo,
		OrderAmount:          req.OrderAmount,
		BalanceType:          constants.XF_BALANCE,
		Type:                 constants.ORDER_TYPE_DF,
		Status:               constants.WAIT_PROCESS,
		Source:               constants.UI,
		IsMerchantCallback:   constants.MERCHANT_CALL_BACK_DONT_USE,
		IsCalculateProfit:    constants.IS_CALCULATE_PROFIT_NO,
		PersonProcessStatus:  constants.PERSON_PROCESS_STATUS_NO_ROCESSING,
		RepaymentStatus:      constants.REPAYMENT_NOT,
		IsTest:               constants.IS_TEST_NO, //是否測試單
		MerchantBankAccount:  req.MerchantBankAccount,
		MerchantBankNo:       req.MerchantBankNo,
		MerchantBankName:     req.MerchantBankName,
		MerchantAccountName:  req.MerchantAccountName,
		MerchantBankProvince: req.MerchantBankProvince,
		MerchantBankCity:     req.MerchantBankCity,
		CurrencyCode:         req.CurrencyCode,
		ChannelCode:          rate.ChannelCode,
		ChannelPayTypesCode:  rate.ChannelPayTypesCode,
		Fee:                  rate.MerFee,
		HandlingFee:          rate.MerHandlingFee,
		TransferHandlingFee:  transferHandlingFee,
		PayTypeCode:          rate.PayTypeCode,
		IsLock:               "0",
	}

	// 新增收支记录，与更新商户余额(商户账户号是黑名单，把交易金额为设为 0)
	updateBalance := &types.UpdateBalance{
		MerchantCode:    txOrder.MerchantCode,
		CurrencyCode:    txOrder.CurrencyCode,
		OrderNo:         txOrder.OrderNo,
		MerchantOrderNo: txOrder.MerchantOrderNo,
		OrderType:       txOrder.Type,
		PayTypeCode:     txOrder.PayTypeCode,
		TransferAmount:  txOrder.TransferAmount,
		TransactionType: "11", //異動類型 (1=收款; 2=解凍; 3=沖正; 11=出款 ; 12=凍結)
		BalanceType:     constants.XF_BALANCE,
		CreatedBy:       txOrder.MerchantCode,
		ChannelCode:     txOrder.ChannelCode,
		MerPtBalanceId:  rate.PtBalanceId,
	}

	//交易金额 = 订单金额 + 商户手续费
	txOrder.TransferAmount = utils.FloatAddCWithTrancated(txOrder.OrderAmount, txOrder.TransferHandlingFee, txOrder.CurrencyCode)
	updateBalance.TransferAmount = txOrder.TransferAmount //扣款依然傳正值

	// 判断单笔最大最小金额
	if txOrder.OrderAmount < rate.SingleMinCharge {
		//金额超过上限
		logx.WithContext(l.ctx).Errorf("錯誤:代付金額未達下限")
		return &transactionclient.ProxyOrderUIResponse{
			Code:    response.ORDER_AMOUNT_LIMIT_MIN,
			Message: "代付金額未達下限，orderNo : " + txOrder.OrderNo,
		}, nil
	} else if txOrder.OrderAmount > rate.SingleMaxCharge {
		//下发金额未达下限
		logx.WithContext(l.ctx).Errorf("錯誤:代付金額超過上限")
		return &transactionclient.ProxyOrderUIResponse{
			Code:    response.ORDER_AMOUNT_LIMIT_MAX,
			Message: "代付金額超过上限，orderNo : " + txOrder.OrderNo,
		}, nil
	}

	redisKey := fmt.Sprintf("%s-%s", updateBalance.MerchantCode, updateBalance.CurrencyCode)
	redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(8)
	if isOK, redisErr := redisLock.TryLockTimeout(8); isOK {
		defer redisLock.Release()

		if err = l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

			//更新商户子钱包且新增记录
			if rate.PtBalanceId > 0 {
				if _, err = merchantbalanceservice.UpdateXF_Pt_Balance_Debit(l.ctx, db, updateBalance); err != nil {
					logx.WithContext(l.ctx).Errorf("商户:%s，更新錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err.Error(), updateBalance)
					return errorz.New(response.SYSTEM_ERROR, err.Error())
				}
			}

			//更新钱包且新增商户钱包异动记录
			if merchantBalanceRecord, err = merchantbalanceservice.UpdateXFBalance_Debit(l.ctx, db, updateBalance); err != nil {
				logx.WithContext(l.ctx).Errorf("商户:%s，更新錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err.Error(), updateBalance)
				return errorz.New(response.SYSTEM_ERROR, err.Error())
			} else {
				logx.WithContext(l.ctx).Infof("代付UI提单 %s，錢包扣款成功", merchantBalanceRecord.OrderNo)
				txOrder.BeforeBalance = merchantBalanceRecord.BeforeBalance // 商戶錢包異動紀錄
				txOrder.Balance = merchantBalanceRecord.AfterBalance
			}

			// 创建订单
			if err = db.Table("tx_orders").Create(&types.OrderX{
				Order: *txOrder,
			}).Error; err != nil {
				logx.WithContext(l.ctx).Errorf("新增代付UI提单失败，商户号: %s, 订单号: %s, err : %s", txOrder.MerchantCode, txOrder.OrderNo, err.Error())
				return
			}

			return nil
		}); err != nil {
			return &transactionclient.ProxyOrderUIResponse{
				Code:         response.DATABASE_FAILURE,
				Message:      "数据库错误 tx_orders Create ui proxy-pay，err : " + err.Error(),
				ProxyOrderNo: txOrder.OrderNo,
			}, nil
		}

	} else {
		logx.WithContext(l.ctx).Errorf("商户钱包处理中，Err:%s。 %s", redisErr.Error(), redisKey)
		return &transactionclient.ProxyOrderUIResponse{
			Code:         response.BALANCE_PROCESSING,
			Message:      i18n.Sprintf(response.BALANCE_PROCESSING),
			ProxyOrderNo: txOrder.OrderNo,
		}, nil
	}

	// 計算利潤(不報錯) TODO: 異步??
	if err4 := orderfeeprofitservice.CalculateOrderProfit(l.svcCtx.MyDB, types.CalculateProfit{
		MerchantCode:        txOrder.MerchantCode,
		OrderNo:             txOrder.OrderNo,
		Type:                txOrder.Type,
		CurrencyCode:        txOrder.CurrencyCode,
		BalanceType:         txOrder.BalanceType,
		ChannelCode:         txOrder.ChannelCode,
		ChannelPayTypesCode: txOrder.ChannelPayTypesCode,
		OrderAmount:         txOrder.OrderAmount,
		IsRate:              rate.IsRate,
	}); err4 != nil {
		logx.WithContext(l.ctx).Errorf("計算利潤出錯:%s", err4.Error())
	}

	// 新單新增訂單歷程 (不抱錯) TODO: 異步??
	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     txOrder.OrderNo,
			Action:      "PLACE_ORDER",
			UserAccount: req.UserAccount,
			Comment:     "",
		},
	}).Error; err4 != nil {
		logx.WithContext(l.ctx).Errorf("紀錄訂單歷程出錯:%s", err4.Error())
	}

	resp = &transactionclient.ProxyOrderUIResponse{
		ProxyOrderNo: txOrder.OrderNo,
		Code:         response.API_SUCCESS,
		Message:      "操作成功",
	}

	return resp, nil
}
