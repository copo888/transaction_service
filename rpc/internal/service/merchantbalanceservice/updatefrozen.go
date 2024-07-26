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

// UpdateFrozenAmount FrozenAmount需正負(凍結正/解凍負), BalanceType:餘額類型 (DFB=代付餘額 XFB=下發餘額)
func UpdateFrozenAmount(db *gorm.DB, updateFrozenAmount types.UpdateFrozenAmount) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {
	var beforeBalance float64
	var afterBalance float64

	// 1. 取得 商戶餘額表
	var merchantBalance types.MerchantBalance
	if err = db.Table("mc_merchant_balances").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("merchant_code = ? AND currency_code = ? AND balance_type = ?", updateFrozenAmount.MerchantCode, updateFrozenAmount.CurrencyCode, updateFrozenAmount.BalanceType).
		Take(&merchantBalance).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 2. 計算
	beforeFrozen := merchantBalance.FrozenAmount
	afterFrozen := utils.FloatAddC(beforeFrozen, updateFrozenAmount.FrozenAmount, updateFrozenAmount.CurrencyCode)
	merchantBalance.FrozenAmount = afterFrozen

	// (依照 BalanceType 決定異動哪種餘額)
	selectBalance := "balance"
	beforeBalance = merchantBalance.Balance
	afterBalance = utils.FloatSubC(beforeBalance, updateFrozenAmount.FrozenAmount, updateFrozenAmount.CurrencyCode)
	merchantBalance.Balance = afterBalance

	if afterFrozen < 0 {
		return merchantBalanceRecord, errorz.New(response.INSUFFICIENT_IN_AMOUNT, err.Error())
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
		OrderNo:           updateFrozenAmount.OrderNo,
		MerchantOrderNo:   updateFrozenAmount.MerchantOrderNo,
		OrderType:         updateFrozenAmount.OrderType,
		ChannelCode:       updateFrozenAmount.ChannelCode,
		PayTypeCode:       updateFrozenAmount.PayTypeCode,
		TransactionType:   updateFrozenAmount.TransactionType,
		BalanceType:       updateFrozenAmount.BalanceType,
		BeforeBalance:     beforeBalance,
		TransferAmount:    -updateFrozenAmount.FrozenAmount,
		AfterBalance:      afterBalance,
		Comment:           updateFrozenAmount.Comment,
		CreatedBy:         updateFrozenAmount.CreatedBy,
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
		OrderNo:           updateFrozenAmount.OrderNo,
		OrderType:         updateFrozenAmount.OrderType,
		TransactionType:   updateFrozenAmount.TransactionType,
		BeforeFrozen:      beforeFrozen,
		FrozenAmount:      updateFrozenAmount.FrozenAmount,
		AfterFrozen:       afterFrozen,
		Comment:           updateFrozenAmount.Comment,
		CreatedBy:         updateFrozenAmount.CreatedBy,
	}

	if err = db.Table("mc_merchant_balance_frozen_records").Create(&types.MerchantFrozenRecordX{
		MerchantFrozenRecord: merchantFrozenRecord,
	}).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	var merchantPtBalanceId int64
	err = db.Table("mc_merchant_channel_rate").
		Select("merchant_pt_balance_id").
		Where("merchant_code = ?", updateFrozenAmount.MerchantCode).
		Where("channel_code = ? AND pay_type_code = ?", updateFrozenAmount.ChannelCode, updateFrozenAmount.PayTypeCode).
		Find(&merchantPtBalanceId).Error

	// 若有啟用顯示子錢包
	if merchantPtBalanceId != 0 {
		// 變更 商戶子錢包餘額
		_, err = merchantPtBalanceService.UpdateFrozenAmount(db, updateFrozenAmount, merchantPtBalanceId)
		if err != nil {
			return
		}
	}

	return
}
