package merchantbalanceservice

import (
	"fmt"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

/*
	更新代付餘額_扣款(代付提單扣款)
*/
func UpdateDFBalance_Debit(db *gorm.DB, updateBalance *types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {
	var beforeBalance float64
	var afterBalance float64

	// 1. 取得 商戶餘額表
	var merchantBalance types.MerchantBalance
	if err = db.Table("mc_merchant_balances").
		Where("merchant_code = ? AND currency_code = ? AND balance_type = ?", updateBalance.MerchantCode, updateBalance.CurrencyCode, updateBalance.BalanceType).
		Take(&merchantBalance).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 2. 計算 (依照 BalanceType 決定異動哪種餘額)
	var selectBalance string
	if utils.FloatAdd(merchantBalance.Balance, -updateBalance.TransferAmount) < 0 { //判斷餘額是否不足
		logx.Errorf("商户:%s，余额类型:%s，余额:%s，交易金额:%s", merchantBalance.MerchantCode, merchantBalance.BalanceType, fmt.Sprintf("%f", merchantBalance.Balance), fmt.Sprintf("%f", updateBalance.TransferAmount))
		return merchantBalanceRecord, errorz.New(response.MERCHANT_INSUFFICIENT_DF_BALANCE)
	}
	selectBalance = "balance"
	beforeBalance = merchantBalance.Balance
	afterBalance = utils.FloatAdd(beforeBalance, -updateBalance.TransferAmount) //代付出款固定TransferAmount代負號
	merchantBalance.Balance = afterBalance

	// 3. 變更 商戶餘額
	if err = db.Table("mc_merchant_balances").Select(selectBalance).Updates(types.MerchantBalanceX{
		MerchantBalance: merchantBalance,
	}).Error; err != nil {
		logx.Error(err.Error())
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
		PayTypeCodeNum:    updateBalance.PayTypeCodeNum,
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

/*
	更新下發餘額(支轉代)_扣款(代付提單扣款)
*/
func UpdateXFBalance_Debit(db *gorm.DB, updateBalance *types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {
	var beforeBalance float64
	var afterBalance float64

	// 1. 取得 商戶餘額表
	var merchantBalance types.MerchantBalance
	if err = db.Table("mc_merchant_balances").
		Where("merchant_code = ? AND currency_code = ? AND balance_type = ?", updateBalance.MerchantCode, updateBalance.CurrencyCode, updateBalance.BalanceType).
		Take(&merchantBalance).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 2. 計算 (依照 BalanceType 決定異動哪種餘額)
	var selectBalance string
	if utils.FloatAdd(merchantBalance.Balance, -updateBalance.TransferAmount) < 0 { //判斷餘額是否不足
		logx.Errorf("商户:%s，余额类型:%s，余额:%s，交易金额:%s", merchantBalance.MerchantCode, merchantBalance.BalanceType, fmt.Sprintf("%f", merchantBalance.Balance), fmt.Sprintf("%f", updateBalance.TransferAmount))
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
		logx.Error(err.Error())
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
		PayTypeCodeNum:    updateBalance.PayTypeCodeNum,
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

// TODO 最後要刪除此方法
// UpdateBalance TransferAmount需正負(收款传正值/扣款传負值), BalanceType:餘額類型 (DFB=代付餘額 XFB=下發餘額)
func UpdateBalance(db *gorm.DB, updateBalance types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {

	var beforeBalance float64
	var afterBalance float64

	// 1. 取得 商戶餘額表
	var merchantBalance types.MerchantBalance
	if err = db.Table("mc_merchant_balances").
		Where("merchant_code = ? AND currency_code = ? AND balance_type = ?", updateBalance.MerchantCode, updateBalance.CurrencyCode, updateBalance.BalanceType).
		Take(&merchantBalance).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 2. 計算 (依照 BalanceType 決定異動哪種餘額)
	var selectBalance string
	if utils.FloatAdd(merchantBalance.Balance, updateBalance.TransferAmount) < 0 {
		logx.Errorf("商户:%s，余额类型:%s，余额:%s，交易金额:%s", merchantBalance.MerchantCode, merchantBalance.BalanceType, fmt.Sprintf("%f", merchantBalance.Balance), fmt.Sprintf("%f", updateBalance.TransferAmount))
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
		logx.Error(err.Error())
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
		PayTypeCodeNum:    updateBalance.PayTypeCodeNum,
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
