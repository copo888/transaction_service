package logic

import (
	"context"
	"fmt"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"github.com/gioco-play/easy-i18n/i18n"
	"github.com/neccoys/go-zero-extension/redislock"
	"gorm.io/gorm"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProxyTestToNormalDFBLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProxyTestToNormalDFBLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyTestToNormalDFBLogic {
	return &ProxyTestToNormalDFBLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProxyTestToNormalDFBLogic) ProxyTestToNormal_DFB(in *transactionclient.ProxyOrderTestRequest) (*transactionclient.ProxyOrderTestResponse, error) {
	txOrder := &types.OrderX{}
	var err error
	if txOrder, err = model.QueryOrderByOrderNo(l.svcCtx.MyDB, in.ProxyOrderNo, ""); err != nil {
		return nil, errorz.New(response.DATABASE_FAILURE, err.Error())
	} else if txOrder == nil {
		return nil, errorz.New(response.ORDER_NUMBER_NOT_EXIST)
	}

	//改非測試單
	txOrder.IsTest = "0"
	txOrder.Memo = "代付订单轉正式單\n" + txOrder.Memo

	redisKey := fmt.Sprintf("%s-%s", txOrder.MerchantCode, txOrder.CurrencyCode)
	redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(5)
	if isOK, redisErr := redisLock.TryLockTimeout(5); isOK {
		defer redisLock.Release()

		l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

			var merchantPtBalanceId int64
			if err = db.Table("mc_merchant_channel_rate").
				Select("merchant_pt_balance_id").
				Where("merchant_code = ? AND channel_pay_types_code = ?", txOrder.MerchantCode, txOrder.ChannelPayTypesCode).
				Find(&merchantPtBalanceId).Error; err != nil {
				logx.WithContext(l.ctx).Errorf("捞取子钱錢包錯誤，商户号:%s，ChannelPayTypesCode:%s，err:%s", txOrder.MerchantCode, txOrder.ChannelPayTypesCode, err.Error())
				return err
			}

			merchantBalanceRecord := types.MerchantBalanceRecord{}

			// 新增收支记录，与更新商户余额(商户账户号是黑名单，把交易金额为设为 0)
			updateBalance := &types.UpdateBalance{
				MerchantCode:    txOrder.MerchantCode,
				CurrencyCode:    txOrder.CurrencyCode,
				OrderNo:         txOrder.OrderNo,
				MerchantOrderNo: txOrder.MerchantOrderNo,
				OrderType:       txOrder.Type,
				PayTypeCode:     txOrder.PayTypeCode,
				TransferAmount:  txOrder.TransferAmount,
				TransactionType: constants.TRANSACTION_TYPE_PROXY_PAY, //異動類型 (1=收款 ; 2=解凍;  3=沖正 4=還款;  5=補單; 11=出款 ; 12=凍結 ; 13=追回; 20=調整)
				BalanceType:     constants.DF_BALANCE,
				Comment:         "代付轉正式單",
				CreatedBy:       txOrder.MerchantCode,
				ChannelCode:     txOrder.ChannelCode,
				MerPtBalanceId:  merchantPtBalanceId,
			}

			if merchantPtBalanceId > 0 {
				if _, err = merchantbalanceservice.DoUpdateDF_Pt_Balance_Debit(l.ctx, l.svcCtx, db, updateBalance); err != nil {
					logx.WithContext(l.ctx).Errorf("商户:%s，幣別: %s，更新子錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, txOrder.CurrencyCode, err.Error(), updateBalance)
					return err
				}
			}

			if merchantBalanceRecord, err = merchantbalanceservice.DoUpdateDFBalance_Debit(l.ctx, l.svcCtx, db, updateBalance); err != nil {
				logx.Errorf("商户:%s，更新錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err.Error(), updateBalance)
				return errorz.New(response.SYSTEM_ERROR, err.Error())
			} else {
				logx.Infof("代付API提单 %s，錢包出款成功", merchantBalanceRecord.OrderNo)
				txOrder.BeforeBalance = merchantBalanceRecord.BeforeBalance // 商戶錢包異動紀錄
				txOrder.Balance = merchantBalanceRecord.AfterBalance
			}

			// 更新订单
			if txOrder != nil {
				if errUpdate := l.svcCtx.MyDB.Table("tx_orders").Updates(txOrder).Error; errUpdate != nil {
					logx.Error("代付订单更新状态错误: ", errUpdate.Error())
				}
			}

			return nil
		})
	} else {
		logx.WithContext(l.ctx).Errorf("商户钱包处理中，Err:%s。 %s", redisErr.Error(), redisKey)
		return &transactionclient.ProxyOrderTestResponse{
			Code:    response.BALANCE_PROCESSING,
			Message: i18n.Sprintf(response.BALANCE_PROCESSING),
		}, nil
	}

	// 更新訂單歷程 (不抱錯)
	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     txOrder.OrderNo,
			Action:      "TRANSFER_NORMAL",
			UserAccount: txOrder.MerchantCode,
			Comment:     "",
		},
	}).Error; err4 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err4.Error())
	}

	return &transactionclient.ProxyOrderTestResponse{}, nil
}
