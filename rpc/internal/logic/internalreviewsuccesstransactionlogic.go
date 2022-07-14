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
	"gorm.io/gorm"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type InternalReviewSuccessTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInternalReviewSuccessTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalReviewSuccessTransactionLogic {
	return &InternalReviewSuccessTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InternalReviewSuccessTransactionLogic) InternalReviewSuccessTransaction(in *transactionclient.InternalReviewSuccessRequest) (resp *transactionclient.InternalReviewSuccessResponse, err error) {
	var txOrder types.OrderX
	var merchantBalanceRecord types.MerchantBalanceRecord

	if err = l.svcCtx.MyDB.Table("tx_orders").Where("order_no = ?", in.OrderNo).Take(&txOrder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &transactionclient.InternalReviewSuccessResponse{
				Code:    response.DATA_NOT_FOUND,
				Message: "找不到资料，orderNo = " + in.OrderNo,
			}, nil
		}
		return
		return &transactionclient.InternalReviewSuccessResponse{
			Code:    response.DATABASE_FAILURE,
			Message: "找不到资料，orderNo = " + in.OrderNo,
		}, nil
	}
	if err = l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {
		// 異動錢包
		if merchantBalanceRecord, err = l.UpdateBalance(db, types.UpdateBalance{
			MerchantCode:    txOrder.MerchantCode,
			CurrencyCode:    txOrder.CurrencyCode,
			OrderNo:         txOrder.OrderNo,
			MerchantOrderNo: txOrder.MerchantOrderNo,
			OrderType:       txOrder.Type,
			ChannelCode:     txOrder.ChannelCode,
			PayTypeCode:     txOrder.PayTypeCode,
			TransactionType: "1",
			BalanceType:     txOrder.BalanceType,
			TransferAmount:  txOrder.TransferAmount,
			Comment:         txOrder.Memo,
			CreatedBy:       in.UserAccount,
		}); err != nil {
			return
		}

		txOrder.BeforeBalance = merchantBalanceRecord.BeforeBalance
		txOrder.Balance = merchantBalanceRecord.AfterBalance
		txOrder.TransAt = types.JsonTime{}.New()
		txOrder.ActualAmount = txOrder.OrderAmount

		txOrder.Status = constants.SUCCESS

		// 編輯訂單
		if err = db.Table("tx_orders").Updates(&txOrder).Error; err != nil {
			return
		}

		return
	}); err != nil {
		return &transactionclient.InternalReviewSuccessResponse{
			Code:    response.UPDATE_DATABASE_FAILURE,
			Message: "数据库错误 tx_orders Update Internal Charge review， err : " + err.Error(),
			OrderNo: in.OrderNo,
		}, nil
	}

	// 新單新增訂單歷程 (不抱錯)
	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     txOrder.OrderNo,
			Action:      "SUCCESS",
			UserAccount: in.UserAccount,
			Comment:     "",
		},
	}).Error; err4 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err4.Error())
	}

	resp = &transactionclient.InternalReviewSuccessResponse{
		OrderNo: txOrder.OrderNo,
		Code: response.API_SUCCESS,
		Message: "操作成功",
	}

	return resp, nil
}

// updateBalance
func (l *InternalReviewSuccessTransactionLogic) UpdateBalance(db *gorm.DB, updateBalance types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {

	var beforeBalance float64
	var afterBalance float64

	// 1. 取得 商戶餘額表
	var merchantBalance types.MerchantBalance
	if err = db.Table("mc_merchant_balances").
		Where("merchant_code = ? AND currency_code = ? AND balance_type = ?", updateBalance.MerchantCode, updateBalance.CurrencyCode, updateBalance.BalanceType).
		Take(&merchantBalance).Error; err != nil {
		return merchantBalanceRecord, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 2. 計算
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
