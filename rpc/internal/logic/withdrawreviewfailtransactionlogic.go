package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"github.com/gioco-play/easy-i18n/i18n"
	"github.com/neccoys/go-zero-extension/redislock"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawReviewFailTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWithdrawReviewFailTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawReviewFailTransactionLogic {
	return &WithdrawReviewFailTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WithdrawReviewFailTransactionLogic) WithdrawReviewFailTransaction(in *transactionclient.WithdrawReviewFailRequest) (resp *transactionclient.WithdrawReviewFailResponse, err error) {
	var txOrder types.OrderX
	var merchantBalanceRecord types.MerchantBalanceRecord

	if err = l.svcCtx.MyDB.Table("tx_orders").Where("order_no = ?", in.OrderNo).Take(&txOrder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &transactionclient.WithdrawReviewFailResponse{
				Code:    response.DATA_NOT_FOUND,
				Message: "找不到资料，orderNo = " + in.OrderNo,
			}, nil
		}
		return &transactionclient.WithdrawReviewFailResponse{
			Code:    response.DATABASE_FAILURE,
			Message: "查询下发订单失败，orderNo = " + in.OrderNo,
		}, nil
	}

	redisKey := fmt.Sprintf("%s-%s", txOrder.MerchantCode, txOrder.CurrencyCode)
	redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(8)
	if isOK, redisErr := redisLock.TryLockTimeout(8); isOK {
		defer redisLock.Release()

		if err = l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

			//  異動錢包 更新商户子钱包馀额
			if in.PtBalanceId > 0 {
				if err = l.UpdatePtBalance(db, l.ctx, types.UpdateBalance{
					MerchantCode:    txOrder.MerchantCode,
					CurrencyCode:    txOrder.CurrencyCode,
					OrderNo:         txOrder.OrderNo,
					MerchantOrderNo: txOrder.MerchantOrderNo,
					OrderType:       txOrder.Type,
					TransactionType: "4",
					BalanceType:     txOrder.BalanceType,
					TransferAmount:  txOrder.TransferAmount,
					Comment:         txOrder.Memo,
					CreatedBy:       in.UserAccount,
					MerPtBalanceId:  in.PtBalanceId,
				}); err != nil {
					return err
				}
			}

			// 異動錢包 更新商户总馀额
			if merchantBalanceRecord, err = l.UpdateBalance(db, l.ctx, types.UpdateBalance{
				MerchantCode:    txOrder.MerchantCode,
				CurrencyCode:    txOrder.CurrencyCode,
				OrderNo:         txOrder.OrderNo,
				MerchantOrderNo: txOrder.MerchantOrderNo,
				OrderType:       txOrder.Type,
				TransactionType: "4",
				BalanceType:     txOrder.BalanceType,
				TransferAmount:  txOrder.TransferAmount,
				Comment:         txOrder.Memo,
				CreatedBy:       in.UserAccount,
			}); err != nil {
				return err
			}

			txOrder.BeforeBalance = merchantBalanceRecord.BeforeBalance
			txOrder.Balance = merchantBalanceRecord.AfterBalance
			txOrder.TransAt = types.JsonTime{}.New()

			txOrder.Status = constants.FAIL
			txOrder.ReviewedBy = in.UserAccount
			txOrder.Memo = in.Memo

			// 編輯訂單
			if err = db.Table("tx_orders").Updates(&txOrder).Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			return &transactionclient.WithdrawReviewFailResponse{
				Code:    response.SYSTEM_ERROR,
				Message: "钱包异动失败，orderNo = " + in.OrderNo,
			}, nil
		}
	} else {
		logx.WithContext(l.ctx).Errorf("商户钱包处理中，Err:%s。 %s", redisErr.Error(), redisKey)
		return &transactionclient.WithdrawReviewFailResponse{
			Code:    response.BALANCE_PROCESSING,
			Message: i18n.Sprintf(response.BALANCE_PROCESSING),
		}, nil
	}

	// 新單新增訂單歷程 (不抱錯)
	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     txOrder.OrderNo,
			Action:      "REVIEW_FAIL",
			UserAccount: in.UserAccount,
			Comment:     in.Memo,
		},
	}).Error; err4 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err4.Error())
	}

	resp = &transactionclient.WithdrawReviewFailResponse{
		OrderNo: txOrder.OrderNo,
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}

	return resp, nil
}

func (l WithdrawReviewFailTransactionLogic) UpdateBalance(db *gorm.DB, ctx context.Context, updateBalance types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {
	//redisKey := fmt.Sprintf("%s-%s-%s", updateBalance.MerchantCode, updateBalance.CurrencyCode, updateBalance.BalanceType)
	//redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	//redisLock.SetExpire(8)
	//if isOk, _ := redisLock.TryLockTimeout(8); isOk {
	//	defer redisLock.Release()
	//	if merchantBalanceRecord, err = l.doUpdateBalance(db, ctx, updateBalance); err != nil {
	//		return
	//	}
	//} else {
	//	return merchantBalanceRecord, errorz.New(response.BALANCE_REDISLOCK_ERROR)
	//}

	merchantBalanceRecord, err = l.doUpdateBalance(db, ctx, updateBalance)

	return
}

func (l WithdrawReviewFailTransactionLogic) UpdatePtBalance(db *gorm.DB, ctx context.Context, updateBalance types.UpdateBalance) error {
	//redisKey := fmt.Sprintf("%s-%s-%s", updateBalance.MerchantCode, updateBalance.CurrencyCode, updateBalance.BalanceType)
	//redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	//redisLock.SetExpire(8)
	//if isOk, _ := redisLock.TryLockTimeout(8); isOk {
	//	defer redisLock.Release()
	//	if err := l.doUpdatePtBalance(db, ctx, updateBalance); err != nil {
	//		return err
	//	}
	//} else {
	//	return errorz.New(response.BALANCE_REDISLOCK_ERROR)
	//}

	if err := l.doUpdatePtBalance(db, ctx, updateBalance); err != nil {
		return err
	}

	return nil
}

func (l WithdrawReviewFailTransactionLogic) doUpdateBalance(db *gorm.DB, ctx context.Context, updateBalance types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {
	var beforeBalance float64
	var afterBalance float64

	// 1. 取得 商戶餘額表
	var merchantBalance types.MerchantBalance
	if err = db.Table("mc_merchant_balances").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("merchant_code = ? AND currency_code = ? AND balance_type = ?", updateBalance.MerchantCode, updateBalance.CurrencyCode, updateBalance.BalanceType).
		Take(&merchantBalance).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 2. 計算
	var selectBalance string
	if utils.FloatAdd(merchantBalance.Balance, updateBalance.TransferAmount) < 0 {
		logx.WithContext(ctx).Errorf("商户:%s，余额类型:%s，余额:%s，交易金额:%s", merchantBalance.MerchantCode, merchantBalance.BalanceType, fmt.Sprintf("%f", merchantBalance.Balance), fmt.Sprintf("%f", updateBalance.TransferAmount))
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
		logx.WithContext(ctx).Error(err.Error())
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

func (l WithdrawReviewFailTransactionLogic) doUpdatePtBalance(db *gorm.DB, ctx context.Context, updateBalance types.UpdateBalance) error {
	var beforeBalance float64
	var afterBalance float64

	// 1. 取得 商戶餘額表
	var merchantPtBalance types.MerchantPtBalance
	if err := db.Table("mc_merchant_pt_balances").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", updateBalance.MerPtBalanceId).
		Take(&merchantPtBalance).Error; err != nil {
		return errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 2. 計算
	var selectBalance string
	if utils.FloatAdd(merchantPtBalance.Balance, updateBalance.TransferAmount) < 0 {
		logx.WithContext(ctx).Errorf("商户:%s，幣別: %s, 子錢包类型:%s ，余额:%s，交易金额:%s", merchantPtBalance.MerchantCode, merchantPtBalance.CurrencyCode, merchantPtBalance.PayTypeCode, fmt.Sprintf("%f", merchantPtBalance.Balance), fmt.Sprintf("%f", updateBalance.TransferAmount))
		return errorz.New(response.MERCHANT_INSUFFICIENT_DF_BALANCE)
	}
	selectBalance = "balance"
	beforeBalance = merchantPtBalance.Balance
	afterBalance = utils.FloatAdd(beforeBalance, updateBalance.TransferAmount)
	merchantPtBalance.Balance = afterBalance

	// 3. 變更 商戶餘額
	if err := db.Table("mc_merchant_pt_balances").Select(selectBalance).Updates(types.MerchantPtBalanceX{
		MerchantPtBalance: merchantPtBalance,
	}).Error; err != nil {
		logx.WithContext(ctx).Error(err.Error())
		return errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 4. 新增 餘額紀錄
	merchantPtBalanceRecord := types.MerchantPtBalanceRecord{
		MerchantPtBalanceId: merchantPtBalance.ID,
		MerchantCode:        merchantPtBalance.MerchantCode,
		CurrencyCode:        merchantPtBalance.CurrencyCode,
		OrderNo:             updateBalance.OrderNo,
		OrderType:           updateBalance.OrderType,
		ChannelCode:         updateBalance.ChannelCode,
		PayTypeCode:         updateBalance.PayTypeCode,
		TransactionType:     updateBalance.TransactionType,
		BeforeBalance:       beforeBalance,
		TransferAmount:      updateBalance.TransferAmount,
		AfterBalance:        afterBalance,
		Comment:             updateBalance.Comment,
		CreatedBy:           updateBalance.CreatedBy,
	}

	if err := db.Table("mc_merchant_pt_balance_records").Create(&types.MerchantPtBalanceRecordX{
		MerchantPtBalanceRecord: merchantPtBalanceRecord,
	}).Error; err != nil {
		return errorz.New(response.DATABASE_FAILURE, err.Error())
	}
	return nil
}
