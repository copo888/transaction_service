package merchantbalanceservice

import (
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantPtBalanceService"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// FrozenManually (手動調整凍結金額) FrozenAmount需正負(凍結正/解凍負), BalanceType:餘額類型 (DFB=代付餘額 XFB=下發餘額)
func FrozenManually(db *gorm.DB, frozenManually types.FrozenManually, merchantPtBalanceId int64) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {
	var beforeBalance float64
	var afterBalance float64
	var transactionType string



	// 1. 取得 商戶餘額表
	var merchantBalance types.MerchantBalance
	if err = db.Table("mc_merchant_balances").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("merchant_code = ? AND currency_code = ? AND balance_type = ?", frozenManually.MerchantCode, frozenManually.CurrencyCode, frozenManually.BalanceType).
		Take(&merchantBalance).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 2. 計算
	beforeFrozen := merchantBalance.FrozenAmount
	afterFrozen := utils.FloatAdd(beforeFrozen, frozenManually.FrozenAmount)
	merchantBalance.FrozenAmount = afterFrozen

	// (依照 BalanceType 決定異動哪種餘額)
	selectBalance := "balance"
	beforeBalance = merchantBalance.Balance
	afterBalance = utils.FloatSub(beforeBalance, frozenManually.FrozenAmount)
	merchantBalance.Balance = afterBalance

	if afterFrozen < 0 {
		return  merchantBalanceRecord, errorz.New(response.INSUFFICIENT_IN_AMOUNT, "總錢包余额不足")
	}

	// 3. 變更 商戶餘額&凍結金額
	if err = db.Table("mc_merchant_balances").Select(selectBalance, "frozen_amount").
		Updates(types.MerchantBalanceX{
			MerchantBalance: merchantBalance,
		}).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 3. 新增 餘額紀錄
	merchantBalanceRecord = types.MerchantBalanceRecord{
		MerchantBalanceId: merchantBalance.ID,
		MerchantCode:      merchantBalance.MerchantCode,
		CurrencyCode:      merchantBalance.CurrencyCode,
		OrderNo:           frozenManually.OrderNo,
		OrderType:         frozenManually.OrderType,
		TransactionType:   frozenManually.TransactionType,
		BalanceType:       frozenManually.BalanceType,
		BeforeBalance:     beforeBalance,
		TransferAmount:    -frozenManually.FrozenAmount,
		AfterBalance:      afterBalance,
		Comment:           frozenManually.Comment,
		CreatedBy:         frozenManually.CreatedBy,
	}

	if err = db.Table("mc_merchant_balance_records").Create(&types.MerchantBalanceRecordX{
		MerchantBalanceRecord: merchantBalanceRecord,
	}).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 4. 新增 凍結紀錄

	merchantFrozenRecord := types.MerchantFrozenRecord{
		MerchantBalanceId: merchantBalance.ID,
		MerchantCode:      merchantBalance.MerchantCode,
		CurrencyCode:      merchantBalance.CurrencyCode,
		OrderNo:           frozenManually.OrderNo,
		OrderType:         frozenManually.OrderType,
		TransactionType:   transactionType,
		BeforeFrozen:      beforeFrozen,
		FrozenAmount:      frozenManually.FrozenAmount,
		AfterFrozen:       afterFrozen,
		Comment:           frozenManually.Comment,
		CreatedBy:         frozenManually.CreatedBy,
	}

	if err = db.Table("mc_merchant_balance_frozen_records").Create(&types.MerchantFrozenRecordX{
		MerchantFrozenRecord: merchantFrozenRecord,
	}).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 若有啟用顯示子錢包
	if merchantPtBalanceId != 0 {
		// 變更 商戶子錢包餘額
		_, err = merchantPtBalanceService.UpdateFrozenAmount(db, types.UpdateFrozenAmount{
			MerchantCode:    frozenManually.MerchantCode,
			CurrencyCode:    frozenManually.CurrencyCode,
			OrderNo:         frozenManually.OrderNo,
			MerchantOrderNo: "",
			OrderType:       frozenManually.OrderType,
			ChannelCode:     "",
			PayTypeCode:     "",
			TransactionType: frozenManually.TransactionType,
			BalanceType:     frozenManually.BalanceType,
			FrozenAmount:    frozenManually.FrozenAmount,
			Comment:         frozenManually.Comment,
			CreatedBy:       frozenManually.CreatedBy,
		}, merchantPtBalanceId)
		if err != nil {
			return
		}
	}

	return
}
