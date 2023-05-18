package merchantPtBalanceService

import (
	"fmt"

	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/go-redis/redis/v8"
	"github.com/neccoys/go-zero-extension/redislock"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
	UpdatePtBalanceForZF 支付異動子錢包
*/
func UpdatePtBalanceForZF(db *gorm.DB, redisClient *redis.Client, updateBalance types.UpdateBalance, merchantPtBalanceId int64) (merchantPtBalanceRecord types.MerchantPtBalanceRecord, err error) {

	redisKey := fmt.Sprintf("%s-%s-%s-%s", merchantPtBalanceRecord.MerchantCode, merchantPtBalanceRecord.CurrencyCode, merchantPtBalanceRecord.ChannelCode, merchantPtBalanceRecord.PayTypeCode)
	redisLock := redislock.New(redisClient, redisKey, "merchant-pt-balance:")
	redisLock.SetExpire(5)

	if isOK, _ := redisLock.TryLockTimeout(5); isOK {
		defer redisLock.Release()
		if merchantPtBalanceRecord, err = doUpdatePtBalanceForZF(db, updateBalance, merchantPtBalanceId); err != nil {
			return
		}
	} else {
		return merchantPtBalanceRecord, errorz.New(response.BALANCE_REDISLOCK_ERROR)
	}

	return
}

func doUpdatePtBalanceForZF(db *gorm.DB, updateBalance types.UpdateBalance, merchantPtBalanceId int64) (merchantPtBalanceRecord types.MerchantPtBalanceRecord, err error) {

	var beforeBalance float64
	var afterBalance float64

	// 1. 取得 商戶餘額表
	var merchantBalance types.MerchantPtBalance
	if err = db.Table("mc_merchant_pt_balances").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", merchantPtBalanceId).
		Take(&merchantBalance).Error; err != nil {
		return merchantPtBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 2. 計算
	beforeBalance = merchantBalance.Balance
	afterBalance = utils.FloatAdd(beforeBalance, updateBalance.TransferAmount)
	merchantBalance.Balance = afterBalance

	// 3. 變更 子錢包餘額
	if err = db.Table("mc_merchant_pt_balances").Select("balance").Updates(types.MerchantPtBalanceX{
		MerchantPtBalance: merchantBalance,
	}).Error; err != nil {
		logx.Error(err.Error())
		return merchantPtBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 4. 新增 餘額紀錄
	merchantPtBalanceRecord = types.MerchantPtBalanceRecord{
		MerchantPtBalanceId: merchantBalance.ID,
		MerchantCode:        merchantBalance.MerchantCode,
		CurrencyCode:        merchantBalance.CurrencyCode,
		OrderNo:             updateBalance.OrderNo,
		OrderType:       	 updateBalance.OrderType,
		MerchantOrderNo:     updateBalance.MerchantOrderNo,
		ChannelCode:         updateBalance.ChannelCode,
		PayTypeCode:         updateBalance.PayTypeCode,
		TransactionType:     updateBalance.TransactionType,
		BeforeBalance:       beforeBalance,
		TransferAmount:      updateBalance.TransferAmount,
		AfterBalance:        afterBalance,
		Comment:             updateBalance.Comment,
		CreatedBy:           updateBalance.CreatedBy,
	}

	if err = db.Table("mc_merchant_pt_balance_records").Create(&types.MerchantPtBalanceRecordX{
		MerchantPtBalanceRecord: merchantPtBalanceRecord,
	}).Error; err != nil {
		return merchantPtBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	return
}
