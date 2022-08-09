package merchantbalanceservice

import (
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UpdateCommissionAmount 異動佣金並增加異動紀錄 Amount需正負
func UpdateCommissionAmount(db *gorm.DB, updateCommissionAmount types.UpdateCommissionAmount) (merchantCommissionRecord types.MerchantCommissionRecord, err error) {
	var beforeCommission float64
	var afterCommission float64

	// 1. 取得 商戶佣金餘額
	var merchantBalance types.MerchantBalance
	if err = db.Table("mc_merchant_balances").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("merchant_code = ? AND currency_code = ? AND balance_type = ?", updateCommissionAmount.MerchantCode, updateCommissionAmount.CurrencyCode, constants.YJ_BALANCE).
		Take(&merchantBalance).Error; err != nil {
		return merchantCommissionRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	beforeCommission = merchantBalance.Balance
	afterCommission = utils.FloatAdd(beforeCommission, updateCommissionAmount.TransferAmount)
	merchantBalance.Balance = afterCommission

	// 3. 變更 商戶佣金
	if err = db.Table("mc_merchant_balances").Select("balance").
		Updates(types.MerchantBalanceX{
			MerchantBalance: merchantBalance,
		}).Error; err != nil {
		return merchantCommissionRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 3. 新增 佣金紀錄
	merchantCommissionRecord = types.MerchantCommissionRecord{
		MerchantBalanceId:       merchantBalance.ID,
		MerchantCode:            merchantBalance.MerchantCode,
		CurrencyCode:            merchantBalance.CurrencyCode,
		CommissionMonthReportId: updateCommissionAmount.CommissionMonthReportId,
		OrderNo:                 updateCommissionAmount.OrderNo,
		TransactionType:         updateCommissionAmount.TransactionType,
		BeforeCommission:        beforeCommission,
		TransferAmount:          updateCommissionAmount.TransferAmount,
		AfterCommission:         afterCommission,
		Comment:                 updateCommissionAmount.Comment,
		CreatedBy:               updateCommissionAmount.CreatedBy,
	}

	if err = db.Table("mc_merchant_balance_commission_records").Create(&types.MerchantCommissionRecordX{
		MerchantCommissionRecord: merchantCommissionRecord,
	}).Error; err != nil {
		return merchantCommissionRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	return
}
