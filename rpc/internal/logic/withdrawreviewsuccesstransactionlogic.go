package logic

import (
	"context"
	"errors"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"gorm.io/gorm"
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
	var txOrder types.OrderX
	var totalWithdrawAmount float64 = 0.0
	var totalChannelHandlingFee float64
	var orderChannels []types.OrderChannelsX
	channelWithdraws := in.ChannelWithdraw

	if err = l.svcCtx.MyDB.Table("tx_orders").Where("order_no = ?", in.OrderNo).Take(&txOrder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorz.New(response.DATA_NOT_FOUND)
		}
		return nil, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	if err = l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

		perChannelWithdrawHandlingFee := utils.FloatDiv(txOrder.HandlingFee, l.intToFloat64(len(channelWithdraws)))
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
				totalWithdrawAmount = utils.FloatAdd(totalWithdrawAmount, channelWithdraw.WithdrawAmount)
				totalChannelHandlingFee = utils.FloatAdd(totalChannelHandlingFee, channelWithdrawHandlingFee)
			}
		}
		// 判断渠道下发金额家总须等于订单的下发金额
		if totalWithdrawAmount != txOrder.OrderAmount {
			return  errorz.New(response.MERCHANT_WITHDRAW_AUDIT_ERROR)
		}

		// 更新下发利润
		oldOrder := &types.Order{
			OrderNo:             txOrder.OrderNo,
			BalanceType:         txOrder.BalanceType,
			TransferHandlingFee: txOrder.TransferHandlingFee,
		}
		l.CalculateSystemProfit(db, oldOrder, totalChannelHandlingFee)

		// 儲存下發明細記錄
		if err1 := db.Table("tx_order_channels").CreateInBatches(orderChannels, len(orderChannels)).Error; err1 != nil {
			return errorz.New(response.DATABASE_FAILURE, err1.Error())
		}
		// 更新交易時間
		txOrder.TransAt = types.JsonTime{}.New()
		txOrder.Status = constants.SUCCESS
		txOrder.IsMerchantCallback = constants.IS_MERCHANT_CALLBACK_YES
		// 更新审核通过
		if err2 := db.Table("tx_orders").Updates(txOrder).Error; err2 != nil {
			return errorz.New(response.DATABASE_FAILURE, err2.Error())
		}

		return nil

	}); err != nil {
		return
	}

	// 新單新增訂單歷程 (不抱錯)
	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     txOrder.OrderNo,
			Action:      "REVIEW_SUCCESS",
			UserAccount: in.UserAccount,
			Comment:     "",
		},
	}).Error; err4 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err4.Error())
	}

	resp = &transactionclient.WithdrawReviewSuccessResponse{
		OrderNo: txOrder.OrderNo,
	}

	return resp, nil
}

func (l *WithdrawReviewSuccessTransactionLogic) intToFloat64(i int) float64 {
	intStr := strconv.Itoa(i)
	res, _ := strconv.ParseFloat(intStr, 64)
	return res
}

func (l *WithdrawReviewSuccessTransactionLogic) CalculateSystemProfit(db *gorm.DB, order *types.Order, TransferHandlingFee float64) (err error){

	systemFeeProfit := types.OrderFeeProfit{
		OrderNo:             order.OrderNo,
		MerchantCode:        "00000000",
		BalanceType:         order.BalanceType,
		Fee:                 0,
		HandlingFee:         TransferHandlingFee,
		TransferHandlingFee: TransferHandlingFee,
		// 商戶手續費 - 渠道總手續費 = 利潤 (有可能是負的)
		ProfitAmount: utils.FloatSub(order.TransferHandlingFee, TransferHandlingFee),
	}

	// 保存系統利潤
	if err = db.Table("tx_orders_fee_profit").Create(&types.OrderFeeProfitX{
		OrderFeeProfit: systemFeeProfit,
	}).Error; err != nil {
		return errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	return
}