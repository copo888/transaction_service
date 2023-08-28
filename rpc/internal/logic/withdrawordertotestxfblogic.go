package logic

import (
	"context"
	"fmt"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"github.com/gioco-play/easy-i18n/i18n"
	"github.com/neccoys/go-zero-extension/redislock"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawOrderToTestXFBLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWithdrawOrderToTestXFBLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawOrderToTestXFBLogic {
	return &WithdrawOrderToTestXFBLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WithdrawOrderToTestXFBLogic) WithdrawOrderToTest_XFB(in *transactionclient.WithdrawOrderTestRequest) (*transactionclient.WithdrawOrderTestResponse, error) {
	merchantPtBalanceId := in.PtBalanceId
	txOrder := &types.OrderX{}
	var err error
	if txOrder, err = model.QueryOrderByOrderNo(l.svcCtx.MyDB, in.WithdrawOrderNo, ""); err != nil {
		return nil, errorz.New(response.DATABASE_FAILURE, err.Error())
	} else if txOrder == nil {
		return nil, errorz.New(response.ORDER_NUMBER_NOT_EXIST)
	}

	redisKey := fmt.Sprintf("%s-%s", txOrder.MerchantCode, txOrder.CurrencyCode)
	redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(5)
	if isOK, redisErr := redisLock.TryLockTimeout(5); isOK {
		defer redisLock.Release()

		//如果月結傭金"已結算/確認報表無誤按鈕" : 不扣款
		txOrder.IsTest = "1"
		txOrder.Memo = "下发订单转测试单" + in.Remark + " \n" + txOrder.Memo

		if err = l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

			merchantBalanceRecord := types.MerchantBalanceRecord{}

			// 新增收支记录，与更新商户余额(商户账户号是黑名单，把交易金额为设为 0)
			updateBalance := types.UpdateBalance{
				MerchantCode:    txOrder.MerchantCode,
				CurrencyCode:    txOrder.CurrencyCode,
				OrderNo:         txOrder.OrderNo,
				MerchantOrderNo: txOrder.MerchantOrderNo,
				OrderType:       txOrder.Type,
				PayTypeCode:     txOrder.PayTypeCode,
				TransferAmount:  txOrder.TransferAmount,
				TransactionType: constants.TRANSACTION_TYPE_REFUND, //異動類型 (1=收款 ; 2=解凍;  3=沖正 4=還款;  5=補單; 11=出款 ; 12=凍結 ; 13=追回; 20=調整)
				BalanceType:     constants.XF_BALANCE,
				Comment:         "下发转测试單",
				CreatedBy:       txOrder.MerchantCode,
				ChannelCode:     txOrder.ChannelCode,
				MerPtBalanceId:  merchantPtBalanceId,
			}

			//异动子钱包
			if merchantPtBalanceId > 0 {
				updateBalance.MerPtBalanceId = merchantPtBalanceId
				if _, err = merchantbalanceservice.UpdateXF_Pt_Balance_Deposit(l.ctx, db, &updateBalance); err != nil {
					txOrder.RepaymentStatus = constants.REPAYMENT_FAIL
					logx.WithContext(l.ctx).Errorf("商户:%s，更新子钱錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err.Error(), updateBalance)
					return err
				} else {
					txOrder.RepaymentStatus = constants.REPAYMENT_SUCCESS
					logx.WithContext(l.ctx).Infof("代付API提单失败 %s，代付錢包退款成功", merchantBalanceRecord.OrderNo)
				}
			}

			if merchantBalanceRecord, err = l.UpdateBalance(db, updateBalance); err != nil {
				logx.Errorf("商户:%s，更新錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err.Error(), updateBalance)
				return errorz.New(response.SYSTEM_ERROR, err.Error())
			} else {
				logx.Infof("下發API提单 %s，錢包還款成功", merchantBalanceRecord.OrderNo)
				txOrder.BeforeBalance = merchantBalanceRecord.BeforeBalance // 商戶錢包異動紀錄
				txOrder.Balance = merchantBalanceRecord.AfterBalance
			}

			// 更新订单
			if txOrder != nil {
				if errUpdate := l.svcCtx.MyDB.Table("tx_orders").Updates(txOrder).Error; errUpdate != nil {
					logx.Error("下發订单更新状态错误: ", errUpdate.Error())
				}
			}

			return nil
		}); err != nil {
			return &transactionclient.WithdrawOrderTestResponse{
				Code:    response.WALLET_UPDATE_ERROR,
				Message: err.Error(),
			}, nil
		}
	} else {
		logx.WithContext(l.ctx).Errorf("商户钱包处理中，Err:%s。 %s", redisErr.Error(), redisKey)
		return &transactionclient.WithdrawOrderTestResponse{
			Code:    response.BALANCE_PROCESSING,
			Message: i18n.Sprintf(response.BALANCE_PROCESSING),
		}, nil
	}

	// 新單新增訂單歷程 (不抱錯) TODO: 異步??
	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     txOrder.OrderNo,
			Action:      "TRANSFER_TEST",
			UserAccount: txOrder.MerchantCode,
			Comment:     "",
		},
	}).Error; err4 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err4.Error())
	}

	return &transactionclient.WithdrawOrderTestResponse{
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}, nil
}

func (l WithdrawOrderToTestXFBLogic) UpdateBalance(db *gorm.DB, updateBalance types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {

	//redisKey := fmt.Sprintf("%s-%s-%s", updateBalance.MerchantCode, updateBalance.CurrencyCode, updateBalance.BalanceType)
	//redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	//redisLock.SetExpire(5)
	//if isOk, _ := redisLock.TryLockTimeout(5); isOk {
	//	defer redisLock.Release()
	//	if merchantBalanceRecord, err = l.doUpdateBalance(db, updateBalance); err != nil {
	//		return
	//	}
	//} else {
	//	return merchantBalanceRecord, errorz.New(response.BALANCE_REDISLOCK_ERROR)
	//}
	merchantBalanceRecord, err = l.doUpdateBalance(db, updateBalance)

	return
}

func (l WithdrawOrderToTestXFBLogic) doUpdateBalance(db *gorm.DB, updateBalance types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {
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
