package merchantbalanceservice

import (
	"context"
	"fmt"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantPtBalanceService"

	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/go-redis/redis/v8"
	"github.com/neccoys/go-zero-extension/redislock"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func DoUpdateDFBalance_Debit(ctx context.Context, svcCtx *svc.ServiceContext, db *gorm.DB, updateBalance *types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {
	var resp *types.MerchantBalanceRecord
	redisKey := fmt.Sprintf("%s-%s-%s", updateBalance.MerchantCode, updateBalance.CurrencyCode, updateBalance.BalanceType)
	redisLock := redislock.New(svcCtx.RedisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(5)
	if isOK, _ := redisLock.TryLockTimeout(5); isOK {
		defer redisLock.Release()
		if resp, err = UpdateDFBalance_Debit(ctx, db, updateBalance); err != nil {
			return types.MerchantBalanceRecord{}, err
		}
	} else {
		return types.MerchantBalanceRecord{}, errorz.New(response.BALANCE_PROCESSING)
	}
	return *resp, nil
}

func DoUpdateXFBalance_Debit(ctx context.Context, svcCtx *svc.ServiceContext, db *gorm.DB, updateBalance *types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {
	var resp *types.MerchantBalanceRecord
	redisKey := fmt.Sprintf("%s-%s-%s", updateBalance.MerchantCode, updateBalance.CurrencyCode, updateBalance.BalanceType)
	redisLock := redislock.New(svcCtx.RedisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(5)
	if isOK, _ := redisLock.TryLockTimeout(5); isOK {
		defer redisLock.Release()
		if resp, err = UpdateXFBalance_Debit(ctx, db, updateBalance); err != nil {
			return types.MerchantBalanceRecord{}, err
		}
	} else {
		return types.MerchantBalanceRecord{}, errorz.New(response.BALANCE_PROCESSING)
	}
	return *resp, nil
}

func DoUpdateDF_Pt_Balance_Debit(ctx context.Context, svcCtx *svc.ServiceContext, db *gorm.DB, updateBalance *types.UpdateBalance) (merchantBalanceRecord types.MerchantPtBalanceRecord, err error) {
	var resp types.MerchantPtBalanceRecord
	redisKey := fmt.Sprintf("%s-%s-%s", updateBalance.MerchantCode, updateBalance.CurrencyCode, updateBalance.PayTypeCode)
	redisLock := redislock.New(svcCtx.RedisClient, redisKey, "merchant-pt-balance:")
	redisLock.SetExpire(5)
	if isOK, _ := redisLock.TryLockTimeout(5); isOK {
		defer redisLock.Release()
		if resp, err = UpdateDF_Pt_Balance_Debit(ctx, db, updateBalance); err != nil {
			return types.MerchantPtBalanceRecord{}, err
		}
	} else {
		return types.MerchantPtBalanceRecord{}, errorz.New(response.BALANCE_PROCESSING)
	}
	return resp, nil
}

func DoUpdateXF_Pt_Balance_Debit(ctx context.Context, svcCtx *svc.ServiceContext, db *gorm.DB, updateBalance *types.UpdateBalance) (merchantBalanceRecord types.MerchantPtBalanceRecord, err error) {
	var resp *types.MerchantPtBalanceRecord
	redisKey := fmt.Sprintf("%s-%s-%s", updateBalance.MerchantCode, updateBalance.CurrencyCode, updateBalance.PayTypeCode)
	redisLock := redislock.New(svcCtx.RedisClient, redisKey, "merchant-pt-balance:")
	redisLock.SetExpire(5)
	if isOK, _ := redisLock.TryLockTimeout(5); isOK {
		defer redisLock.Release()
		if resp, err = UpdateXF_Pt_Balance_Debit(ctx, db, updateBalance); err != nil {
			return types.MerchantPtBalanceRecord{}, err
		}
	} else {
		return types.MerchantPtBalanceRecord{}, errorz.New(response.BALANCE_PROCESSING)
	}
	return *resp, nil
}

/*
	更新代付餘額_扣款(代付提單扣款)
*/
func UpdateDFBalance_Debit(ctx context.Context, db *gorm.DB, updateBalance *types.UpdateBalance) (merchantBalanceRecord *types.MerchantBalanceRecord, err error) {
	var beforeBalance float64
	var afterBalance float64

	// 1. 取得 商戶餘額表
	var merchantBalance *types.MerchantBalance
	if err = db.Table("mc_merchant_balances").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("merchant_code = ? AND currency_code = ? AND balance_type = ?", updateBalance.MerchantCode, updateBalance.CurrencyCode, updateBalance.BalanceType).
		Take(&merchantBalance).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 2. 計算 (依照 BalanceType 決定異動哪種餘額)
	var selectBalance string
	if utils.FloatAdd(merchantBalance.Balance, -updateBalance.TransferAmount) < 0 { //判斷餘額是否不足
		logx.WithContext(ctx).Errorf("商户:%s，余额类型:%s，余额:%s，交易金额:%s", merchantBalance.MerchantCode, merchantBalance.BalanceType, fmt.Sprintf("%f", merchantBalance.Balance), fmt.Sprintf("%f", updateBalance.TransferAmount))
		return merchantBalanceRecord, errorz.New(response.MERCHANT_INSUFFICIENT_DF_BALANCE)
	}
	selectBalance = "balance"
	beforeBalance = merchantBalance.Balance
	afterBalance = utils.FloatAdd(beforeBalance, -updateBalance.TransferAmount) //代付出款固定TransferAmount代負號
	merchantBalance.Balance = afterBalance

	// 3. 變更 商戶餘額
	if err = db.Table("mc_merchant_balances").Select(selectBalance).Updates(types.MerchantBalanceX{
		MerchantBalance: *merchantBalance,
	}).Error; err != nil {
		logx.WithContext(ctx).Error(err.Error())
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 4. 新增 餘額紀錄
	merchantBalanceRecord = &types.MerchantBalanceRecord{
		MerchantBalanceId: merchantBalance.ID,
		MerchantCode:      merchantBalance.MerchantCode,
		CurrencyCode:      merchantBalance.CurrencyCode,
		OrderNo:           updateBalance.OrderNo,
		MerchantOrderNo:   updateBalance.MerchantOrderNo,
		OrderType:         updateBalance.OrderType,
		ChannelCode:       updateBalance.ChannelCode,
		PayTypeCode:       updateBalance.PayTypeCode,
		TransactionType:   updateBalance.TransactionType,
		BalanceType:       updateBalance.BalanceType,
		BeforeBalance:     beforeBalance,
		TransferAmount:    -updateBalance.TransferAmount,
		AfterBalance:      afterBalance,
		Comment:           updateBalance.Comment,
		CreatedBy:         updateBalance.CreatedBy,
	}

	if err = db.Table("mc_merchant_balance_records").Create(&types.MerchantBalanceRecordX{
		MerchantBalanceRecord: *merchantBalanceRecord,
	}).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	return
}

//更新子錢包餘額_扣款(代付提單扣款)
func UpdateDF_Pt_Balance_Debit(ctx context.Context, db *gorm.DB, updateBalance *types.UpdateBalance) (merchantPtBalanceRecord types.MerchantPtBalanceRecord, err error) {
	var beforeBalance float64
	var afterBalance float64

	// 1. 取得 商戶餘額表
	var merchantPtBalance types.MerchantPtBalance
	if err = db.Table("mc_merchant_pt_balances").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", updateBalance.MerPtBalanceId).
		Take(&merchantPtBalance).Error; err != nil {
		return types.MerchantPtBalanceRecord{}, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	//判斷餘額是否不足 2. 計算 (依照 BalanceType 決定異動哪種餘額)
	if utils.FloatAdd(merchantPtBalance.Balance, -updateBalance.TransferAmount) < 0 { //判斷餘額是否不足
		logx.WithContext(ctx).Errorf("商户:%s，幣別:%s, 子錢包类型:%s ，余额:%s，交易金额:%s", merchantPtBalance.MerchantCode, merchantPtBalance.CurrencyCode, merchantPtBalance.Name, fmt.Sprintf("%f", merchantPtBalance.Balance), fmt.Sprintf("%f", updateBalance.TransferAmount))
		return merchantPtBalanceRecord, errorz.New(fmt.Sprintf("商户子錢包餘額不足:%s，幣別:%s, 子錢包类型:%s ，余额:%s，交易金额:%s", merchantPtBalance.MerchantCode, merchantPtBalance.CurrencyCode, merchantPtBalance.Name, fmt.Sprintf("%f", merchantPtBalance.Balance), fmt.Sprintf("%f", updateBalance.TransferAmount)))
	}

	beforeBalance = merchantPtBalance.Balance
	afterBalance = utils.FloatAdd(beforeBalance, -updateBalance.TransferAmount)
	merchantPtBalance.Balance = afterBalance

	// 3. 變更 子錢包餘額
	if err = db.Table("mc_merchant_pt_balances").Select("balance").Updates(&types.MerchantPtBalanceX{
		MerchantPtBalance: merchantPtBalance,
	}).Error; err != nil {
		logx.WithContext(ctx).Error("更新子錢包餘額錯誤: %s", err.Error())
		logx.WithContext(ctx).Error(err.Error())
		return merchantPtBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 4. 新增 餘額紀錄
	merchantPtBalanceRecord = types.MerchantPtBalanceRecord{
		MerchantPtBalanceId: merchantPtBalance.ID,
		MerchantCode:        merchantPtBalance.MerchantCode,
		CurrencyCode:        merchantPtBalance.CurrencyCode,
		OrderNo:             updateBalance.OrderNo,
		MerchantOrderNo:     updateBalance.MerchantOrderNo,
		ChannelCode:         updateBalance.ChannelCode,
		PayTypeCode:         updateBalance.PayTypeCode,
		TransactionType:     updateBalance.TransactionType,
		BeforeBalance:       beforeBalance,
		TransferAmount:      updateBalance.TransferAmount,
		AfterBalance:        afterBalance,
		Comment:             updateBalance.Comment,
		CreatedBy:           updateBalance.CreatedBy,
		OrderType:           "DF",
	}

	if err = db.Table("mc_merchant_pt_balance_records").Create(&types.MerchantPtBalanceRecordX{
		MerchantPtBalanceRecord: merchantPtBalanceRecord,
	}).Error; err != nil {
		logx.WithContext(ctx).Error("更新子錢包紀錄錯誤: %s", err.Error())
		return merchantPtBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	return
}

//更新子錢包餘額_扣款(下發提單扣款)
func UpdateXF_Pt_Balance_Debit(ctx context.Context, db *gorm.DB, updateBalance *types.UpdateBalance) (merchantPtBalanceRecord *types.MerchantPtBalanceRecord, err error) {
	var beforeBalance float64
	var afterBalance float64

	// 1. 取得 商戶餘額表
	var merchantPtBalance *types.MerchantPtBalance
	if err = db.Table("mc_merchant_pt_balances").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", updateBalance.MerPtBalanceId).
		Take(&merchantPtBalance).Error; err != nil {
		return nil, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	//判斷餘額是否不足 2. 計算 (依照 BalanceType 決定異動哪種餘額)
	if utils.FloatAdd(merchantPtBalance.Balance, -updateBalance.TransferAmount) < 0 { //判斷餘額是否不足
		logx.WithContext(ctx).Errorf("商户:%s，幣別: %s, 子錢包类型:%s ，余额:%s，交易金额:%s", merchantPtBalance.MerchantCode, merchantPtBalance.CurrencyCode, merchantPtBalance.PayTypeCode, fmt.Sprintf("%f", merchantPtBalance.Balance), fmt.Sprintf("%f", updateBalance.TransferAmount))
		return merchantPtBalanceRecord, errorz.New(response.MERCHANT_INSUFFICIENT_DF_BALANCE)
	}

	beforeBalance = merchantPtBalance.Balance
	afterBalance = utils.FloatAdd(beforeBalance, -updateBalance.TransferAmount)
	merchantPtBalance.Balance = afterBalance

	// 3. 變更 子錢包餘額
	if err = db.Table("mc_merchant_pt_balances").Select("balance").Updates(types.MerchantPtBalanceX{
		MerchantPtBalance: *merchantPtBalance,
	}).Error; err != nil {
		logx.WithContext(ctx).Error(err.Error())
		return merchantPtBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 4. 新增 餘額紀錄
	merchantPtBalanceRecord = &types.MerchantPtBalanceRecord{
		MerchantPtBalanceId: merchantPtBalance.ID,
		MerchantCode:        merchantPtBalance.MerchantCode,
		CurrencyCode:        merchantPtBalance.CurrencyCode,
		OrderNo:             updateBalance.OrderNo,
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
		MerchantPtBalanceRecord: *merchantPtBalanceRecord,
	}).Error; err != nil {
		return merchantPtBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	return
}

/*
	更新子錢包餘額_代付失败退回
*/
func UpdateDF_Pt_Balance_Deposit(ctx context.Context, db *gorm.DB, updateBalance *types.UpdateBalance) (merchantPtBalanceRecord *types.MerchantPtBalanceRecord, err error) {
	var beforeBalance float64
	var afterBalance float64

	// 1. 取得 商戶餘額表
	var merchantPtBalance types.MerchantPtBalance
	if err = db.Table("mc_merchant_pt_balances").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", updateBalance.MerPtBalanceId).
		Take(&merchantPtBalance).Error; err != nil {
		return merchantPtBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	beforeBalance = merchantPtBalance.Balance
	afterBalance = utils.FloatAdd(beforeBalance, updateBalance.TransferAmount)
	merchantPtBalance.Balance = afterBalance

	// 3. 變更 子錢包餘額
	if err = db.Table("mc_merchant_pt_balances").Select("balance").Updates(&types.MerchantPtBalanceX{
		MerchantPtBalance: merchantPtBalance,
	}).Error; err != nil {
		logx.WithContext(ctx).Error("更新子錢包餘額錯誤: %s", err.Error())
		return merchantPtBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 4. 新增 餘額紀錄
	merchantPtBalanceRecord = &types.MerchantPtBalanceRecord{
		MerchantPtBalanceId: merchantPtBalance.ID,
		MerchantCode:        merchantPtBalance.MerchantCode,
		CurrencyCode:        merchantPtBalance.CurrencyCode,
		OrderNo:             updateBalance.OrderNo,
		MerchantOrderNo:     updateBalance.MerchantOrderNo,
		ChannelCode:         updateBalance.ChannelCode,
		PayTypeCode:         updateBalance.PayTypeCode,
		TransactionType:     updateBalance.TransactionType,
		BeforeBalance:       beforeBalance,
		TransferAmount:      updateBalance.TransferAmount,
		AfterBalance:        afterBalance,
		Comment:             updateBalance.Comment,
		CreatedBy:           updateBalance.CreatedBy,
		OrderType:           "DF",
	}

	if err = db.Table("mc_merchant_pt_balance_records").Create(&types.MerchantPtBalanceRecordX{
		MerchantPtBalanceRecord: *merchantPtBalanceRecord,
	}).Error; err != nil {
		logx.WithContext(ctx).Error("更新子錢包紀錄錯誤: %s", err.Error())
		return merchantPtBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	return
}

func UpdateXF_Pt_Balance_Deposit(ctx context.Context, db *gorm.DB, updateBalance *types.UpdateBalance) (merchantPtBalanceRecord *types.MerchantPtBalanceRecord, err error) {
	var beforeBalance float64
	var afterBalance float64

	// 1. 取得 商戶餘額表
	var merchantPtBalance types.MerchantPtBalance
	if err = db.Table("mc_merchant_pt_balances").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", updateBalance.MerPtBalanceId).
		Take(&merchantPtBalance).Error; err != nil {
		return merchantPtBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	beforeBalance = merchantPtBalance.Balance
	afterBalance = utils.FloatAdd(beforeBalance, updateBalance.TransferAmount)
	merchantPtBalance.Balance = afterBalance

	// 3. 變更 子錢包餘額
	if err = db.Table("mc_merchant_pt_balances").Select("balance").Updates(&types.MerchantPtBalanceX{
		MerchantPtBalance: merchantPtBalance,
	}).Error; err != nil {
		logx.WithContext(ctx).Error("更新子錢包餘額錯誤: %s", err.Error())
		return merchantPtBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 4. 新增 餘額紀錄
	merchantPtBalanceRecord = &types.MerchantPtBalanceRecord{
		MerchantPtBalanceId: merchantPtBalance.ID,
		MerchantCode:        merchantPtBalance.MerchantCode,
		CurrencyCode:        merchantPtBalance.CurrencyCode,
		OrderNo:             updateBalance.OrderNo,
		MerchantOrderNo:     updateBalance.MerchantOrderNo,
		ChannelCode:         updateBalance.ChannelCode,
		PayTypeCode:         updateBalance.PayTypeCode,
		TransactionType:     updateBalance.TransactionType,
		BeforeBalance:       beforeBalance,
		TransferAmount:      updateBalance.TransferAmount,
		AfterBalance:        afterBalance,
		Comment:             updateBalance.Comment,
		CreatedBy:           updateBalance.CreatedBy,
		OrderType:           "XF",
	}

	if err = db.Table("mc_merchant_pt_balance_records").Create(&types.MerchantPtBalanceRecordX{
		MerchantPtBalanceRecord: *merchantPtBalanceRecord,
	}).Error; err != nil {
		logx.WithContext(ctx).Error("更新子錢包紀錄錯誤: %s", err.Error())
		return merchantPtBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	return
}

/*
	更新下發餘額(支轉代)_扣款(代付提單扣款)
*/
func UpdateXFBalance_Debit(ctx context.Context, db *gorm.DB, updateBalance *types.UpdateBalance) (merchantBalanceRecord *types.MerchantBalanceRecord, err error) {
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

	// 2. 計算 (依照 BalanceType 決定異動哪種餘額)
	var selectBalance string
	if utils.FloatAdd(merchantBalance.Balance, -updateBalance.TransferAmount) < 0 { //判斷餘額是否不足
		logx.WithContext(ctx).Errorf("商户:%s，余额类型:%s，余额:%s，交易金额:%s", merchantBalance.MerchantCode, merchantBalance.BalanceType, fmt.Sprintf("%f", merchantBalance.Balance), fmt.Sprintf("%f", updateBalance.TransferAmount))
		return merchantBalanceRecord, errorz.New(response.MERCHANT_INSUFFICIENT_DF_BALANCE)
	}
	selectBalance = "balance"
	beforeBalance = merchantBalance.Balance
	afterBalance = utils.FloatAdd(beforeBalance, -updateBalance.TransferAmount)
	merchantBalance.Balance = afterBalance

	// 3. 變更 商戶餘額
	if err = db.Table("mc_merchant_balances").Select(selectBalance).Updates(types.MerchantBalanceX{
		MerchantBalance: merchantBalance,
	}).Error; err != nil {
		logx.WithContext(ctx).Error(err.Error())
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 4. 新增 餘額紀錄
	merchantBalanceRecord = &types.MerchantBalanceRecord{
		MerchantBalanceId: merchantBalance.ID,
		MerchantCode:      merchantBalance.MerchantCode,
		CurrencyCode:      merchantBalance.CurrencyCode,
		OrderNo:           updateBalance.OrderNo,
		MerchantOrderNo:   updateBalance.MerchantOrderNo,
		OrderType:         updateBalance.OrderType,
		ChannelCode:       updateBalance.ChannelCode,
		PayTypeCode:       updateBalance.PayTypeCode,
		TransactionType:   updateBalance.TransactionType,
		BalanceType:       updateBalance.BalanceType,
		BeforeBalance:     beforeBalance,
		TransferAmount:    -updateBalance.TransferAmount,
		AfterBalance:      afterBalance,
		Comment:           updateBalance.Comment,
		CreatedBy:         updateBalance.CreatedBy,
	}

	if err = db.Table("mc_merchant_balance_records").Create(&types.MerchantBalanceRecordX{
		MerchantBalanceRecord: *merchantBalanceRecord,
	}).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	return
}

/*
	更新代付余额_下發余额(代付失败退回)

*/
func UpdateXFBalance_Deposit(db *gorm.DB, updateBalance types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {
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
	selectBalance := "balance"
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
		MerchantOrderNo:   updateBalance.MerchantOrderNo,
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

/*
	更新代付余额_代付余额(代付失败退回)

*/
func UpdateDFBalance_Deposit(db *gorm.DB, updateBalance *types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {
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
	selectBalance := "balance"
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
		MerchantOrderNo:   updateBalance.MerchantOrderNo,
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

/*
	UpdateBalanceForZF 支付異動錢包
*/
func UpdateBalanceForZF(db *gorm.DB, ctx context.Context, redisClient *redis.Client, updateBalance types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {

	redisKey := fmt.Sprintf("%s-%s-%s", merchantBalanceRecord.MerchantCode, merchantBalanceRecord.CurrencyCode, merchantBalanceRecord.BalanceType)
	redisLock := redislock.New(redisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(5)

	if isOK, _ := redisLock.TryLockTimeout(5); isOK {
		defer redisLock.Release()
		if merchantBalanceRecord, err = doUpdateBalanceForZF(db, ctx, redisClient, updateBalance); err != nil {
			return
		}
	} else {
		return merchantBalanceRecord, errorz.New(response.BALANCE_REDISLOCK_ERROR)
	}

	return
}

func doUpdateBalanceForZF(db *gorm.DB, ctx context.Context, redisClient *redis.Client, updateBalance types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {

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
	beforeBalance = merchantBalance.Balance
	afterBalance = utils.FloatAdd(beforeBalance, updateBalance.TransferAmount)
	merchantBalance.Balance = afterBalance

	// 3. 變更 商戶餘額
	if err = db.Table("mc_merchant_balances").Select("balance").Updates(types.MerchantBalanceX{
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
		MerchantOrderNo:   updateBalance.MerchantOrderNo,
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

	var merchantPtBalanceId int64
	err = db.Table("mc_merchant_channel_rate").
		Select("merchant_pt_balance_id").
		Where("merchant_code = ?", updateBalance.MerchantCode).
		Where("channel_code = ? AND pay_type_code = ?", updateBalance.ChannelCode, updateBalance.PayTypeCode).
		Find(&merchantPtBalanceId).Error

	// 若有啟用顯示子錢包
	if merchantPtBalanceId != 0 {
		// 變更 商戶子錢包餘額
		_, err = merchantPtBalanceService.UpdatePtBalanceForZF(db, redisClient, updateBalance, merchantPtBalanceId)
		if err != nil {
			return
		}
	}

	return
}

/*
	更新子錢包餘額_下發失败退回
*/
// UpdateBalance TransferAmount需正負(收款传正值/扣款传負值), BalanceType:餘額類型 (DFB=代付餘額 XFB=下發餘額)
func UpdateBalance(db *gorm.DB, updateBalance types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {

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

	// 2. 計算 (依照 BalanceType 決定異動哪種餘額)
	var selectBalance string

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
		MerchantOrderNo:   updateBalance.MerchantOrderNo,
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
