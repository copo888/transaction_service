package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/redislock"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"github.com/gioco-play/easy-i18n/i18n"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawReviewSuccessTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWithdrawReviewSuccessTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawReviewSuccessTransactionLogic {
	return &WithdrawReviewSuccessTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WithdrawReviewSuccessTransactionLogic) WithdrawReviewSuccessTransaction(in *transactionclient.WithdrawReviewSuccessRequest) (resp *transactionclient.WithdrawReviewSuccessResponse, err error) {
	var txOrder = &types.OrderX{}
	var merchantPtBalanceId int64
	var totalWithdrawAmount float64 = 0.0
	var totalChannelHandlingFee float64
	var orderChannels []types.OrderChannelsX
	channelWithdraws := in.ChannelWithdraw

	if err = l.svcCtx.MyDB.Table("tx_orders").Where("order_no = ?", in.OrderNo).Take(txOrder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &transactionclient.WithdrawReviewSuccessResponse{
				Code:    response.DATA_NOT_FOUND,
				Message: "下发订单查无资料，orderNo = " + in.OrderNo,
			}, nil
		}
		return &transactionclient.WithdrawReviewSuccessResponse{
			Code:    response.DATABASE_FAILURE,
			Message: "数据库错误，err : " + err.Error(),
		}, nil
	}

	if in.IsCharged == "1" {
		if err = l.svcCtx.MyDB.Table("mc_merchant_pt_balance_records").
			Select("merchant_pt_balance_id").Where("order_no = ?", txOrder.OrderNo).Take(&merchantPtBalanceId).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return &transactionclient.WithdrawReviewSuccessResponse{
				Code:    response.DATABASE_FAILURE,
				Message: "数据库错误，查询子钱包错误，err : " + err.Error(),
			}, nil
		}
	}

	redisKey := fmt.Sprintf("%s-%s", txOrder.MerchantCode, txOrder.CurrencyCode)
	redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(8)
	if isOK, redisErr := redisLock.TryLockTimeout(8); isOK {
		defer redisLock.Release()

		if err = l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

			//下發審核若低於金額，則需還收手續費(因維手續費是預先扣除)
			if in.IsCharged == "1" {
				txOrder.Fee = 0.0
				//txOrder.HandlingFee = 0.0
				txOrder.TransferHandlingFee = 0.0
				txOrder.TransferAmount = txOrder.OrderAmount

				merchantBalanceRecord := types.MerchantBalanceRecord{}

				// 新增收支记录，与更新商户余额(商户账户号是黑名单，把交易金额为设为 0)
				updateBalance := types.UpdateBalance{
					MerchantCode:    txOrder.MerchantCode,
					CurrencyCode:    txOrder.CurrencyCode,
					OrderNo:         txOrder.OrderNo,
					MerchantOrderNo: txOrder.MerchantOrderNo,
					OrderType:       txOrder.Type,
					PayTypeCode:     txOrder.PayTypeCode,
					TransferAmount:  txOrder.DefaultHandlingFee,
					TransactionType: constants.TRANSACTION_TYPE_REFUND, //異動類型 (1=收款 ; 2=解凍;  3=沖正 4=還款;  5=補單; 11=出款 ; 12=凍結 ; 13=追回; 20=調整)
					BalanceType:     constants.XF_BALANCE,
					Comment:         "還款下發手續費",
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

				if merchantBalanceRecord, err = l.doUpdateBalance(db, updateBalance); err != nil {
					logx.Errorf("商户:%s，更新錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err.Error(), updateBalance)
					return errorz.New(response.SYSTEM_ERROR, err.Error())
				} else {
					logx.Infof("API提单 %s，錢包出款成功", merchantBalanceRecord.OrderNo)
					txOrder.BeforeBalance = merchantBalanceRecord.BeforeBalance // 商戶錢包異動紀錄
					txOrder.Balance = merchantBalanceRecord.AfterBalance
				}
			}

			perChannelWithdrawHandlingFee := utils.FloatDivCWithTrancated(txOrder.HandlingFee, l.intToFloat64(len(channelWithdraws)), txOrder.CurrencyCode)
			for _, channelWithdraw := range channelWithdraws {
				if channelWithdraw.WithdrawAmount > 0 { // 下发金额不得为0
					// 取得渠道下發手續費
					var channelWithdrawHandlingFee float64

					if err1 := l.svcCtx.MyDB.Table("ch_channels").Select("channel_withdraw_charge").Where("code = ?", channelWithdraw.ChannelCode).
						Take(&channelWithdrawHandlingFee).Error; err1 != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							return errorz.New(response.DATA_NOT_FOUND)
						}
						return errorz.New(response.DATABASE_FAILURE, err1.Error())
					}

					// 记录下发记录
					orderChannel := types.OrderChannelsX{
						OrderChannels: types.OrderChannels{
							OrderNo:             txOrder.OrderNo,
							ChannelCode:         channelWithdraw.ChannelCode,
							HandlingFee:         channelWithdrawHandlingFee,
							OrderAmount:         channelWithdraw.WithdrawAmount,
							TransferHandlingFee: perChannelWithdrawHandlingFee,
						},
					}

					orderChannels = append(orderChannels, orderChannel)
					totalWithdrawAmount = utils.FloatAddCWithTrancated(totalWithdrawAmount, channelWithdraw.WithdrawAmount, txOrder.CurrencyCode)
					totalChannelHandlingFee = utils.FloatAddCWithTrancated(totalChannelHandlingFee, channelWithdrawHandlingFee, txOrder.CurrencyCode)
				}
			}
			// 判断渠道下发金额家总须等于订单的下发金额
			if totalWithdrawAmount != txOrder.OrderAmount {
				return errorz.New(response.MERCHANT_WITHDRAW_AUDIT_ERROR)
			}

			// 儲存下發明細記錄
			if err1 := db.Table("tx_order_channels").CreateInBatches(orderChannels, len(orderChannels)).Error; err1 != nil {
				return errorz.New(response.DATABASE_FAILURE, err1.Error())
			}
			// 更新交易時間
			txOrder.TransAt = types.JsonTime{}.New()
			txOrder.Status = constants.SUCCESS
			txOrder.ReviewedBy = in.UserAccount
			txOrder.Memo = in.Memo

			if err = db.Table("tx_orders").Where("id = ?", txOrder.ID).Updates(txOrder).Error; err != nil {
				return errorz.New(response.DATABASE_FAILURE, err.Error())
			}
			//下發回U還商戶手續費
			if in.IsCharged == "1" {
				if err = db.Table("tx_orders").Where("id = ?", txOrder.ID).
					Updates(map[string]interface{}{"transfer_handling_fee": 0.0, "fee": 0.0}).Error; err != nil {
					return errorz.New(response.DATABASE_FAILURE, err.Error())
				}
			}

			// 更新下发利润
			oldOrder := &types.Order{
				OrderNo:             txOrder.OrderNo,
				BalanceType:         txOrder.BalanceType,
				TransferHandlingFee: txOrder.TransferHandlingFee,
			}
			if err = l.CalculateSystemProfit(db, oldOrder, totalChannelHandlingFee); err != nil {
				logx.Errorf("审核通过，计算下发利润失败，商户号: %s, 订单号: %s, err : %s", txOrder.MerchantCode, txOrder.OrderNo, err.Error())
				return err
			}

			return nil

		}); err != nil {
			return &transactionclient.WithdrawReviewSuccessResponse{
				Code:    response.UPDATE_DATABASE_FAILURE,
				Message: "下发审核更新失敗，orderNo = " + in.OrderNo + "，err : " + err.Error(),
			}, nil
		}

	} else {
		logx.WithContext(l.ctx).Errorf("商户钱包处理中，Err:%s。 %s", redisErr.Error(), redisKey)
		return &transactionclient.WithdrawReviewSuccessResponse{
			Code:    response.BALANCE_PROCESSING,
			Message: i18n.Sprintf(response.BALANCE_PROCESSING),
		}, nil
	}

	// 新單新增訂單歷程 (不抱錯)
	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     txOrder.OrderNo,
			Action:      "REVIEW_SUCCESS",
			UserAccount: in.UserAccount,
			Comment:     txOrder.Memo,
		},
	}).Error; err4 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err4.Error())
	}

	resp = &transactionclient.WithdrawReviewSuccessResponse{
		OrderNo: txOrder.OrderNo,
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}

	return resp, nil
}

func (l *WithdrawReviewSuccessTransactionLogic) intToFloat64(i int) float64 {
	intStr := strconv.Itoa(i)
	res, _ := strconv.ParseFloat(intStr, 64)
	return res
}

func (l *WithdrawReviewSuccessTransactionLogic) CalculateSystemProfit(db *gorm.DB, order *types.Order, TransferHandlingFee float64) (err error) {

	systemFeeProfit := types.OrderFeeProfit{
		OrderNo:             order.OrderNo,
		MerchantCode:        "00000000",
		BalanceType:         order.BalanceType,
		Fee:                 0,
		HandlingFee:         TransferHandlingFee,
		TransferHandlingFee: TransferHandlingFee,
		// 商戶手續費 - 渠道總手續費 = 利潤 (有可能是負的)
		ProfitAmount: utils.FloatSubCWithTrancated(order.TransferHandlingFee, TransferHandlingFee, order.CurrencyCode),
	}

	// 保存系統利潤
	if err = db.Table("tx_orders_fee_profit").Create(&types.OrderFeeProfitX{
		OrderFeeProfit: systemFeeProfit,
	}).Error; err != nil {
		return errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	if err = l.updateOrderByIsCalculateProfit(db, order.OrderNo); err != nil {
		return errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	return
}

func (l *WithdrawReviewSuccessTransactionLogic) updateOrderByIsCalculateProfit(db *gorm.DB, orderNo string) error {
	return db.Table("tx_orders").
		Where("order_no = ?", orderNo).
		Updates(map[string]interface{}{"is_calculate_profit": constants.IS_CALCULATE_PROFIT_YES}).Error
}

func (l WithdrawReviewSuccessTransactionLogic) doUpdateBalance(db *gorm.DB, updateBalance types.UpdateBalance) (merchantBalanceRecord types.MerchantBalanceRecord, err error) {
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
	afterBalance = utils.FloatAddCWithTrancated(beforeBalance, updateBalance.TransferAmount, updateBalance.CurrencyCode)
	merchantBalance.Balance = afterBalance

	// 3. 變更 商戶餘額
	if err = db.Table("mc_merchant_balances").Select(selectBalance).Updates(types.MerchantBalanceX{
		MerchantBalance: merchantBalance,
	}).Error; err != nil {
		logx.WithContext(l.ctx).Errorf("mc_merchant_balances Err: %s", err.Error())
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
