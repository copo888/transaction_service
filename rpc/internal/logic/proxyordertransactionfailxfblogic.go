package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/gormx"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"gorm.io/gorm"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProxyOrderTransactionFailXFBLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProxyOrderTransactionFailXFBLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyOrderTransactionFailXFBLogic {
	return &ProxyOrderTransactionFailXFBLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProxyOrderTransactionFailXFBLogic) ProxyOrderTransactionFail_XFB(in *transactionclient.ProxyPayFailRequest) (resp *transactionclient.ProxyPayFailResponse, err error) {
	merchantBalanceRecord := types.MerchantBalanceRecord{}
	var txOrder = &types.OrderX{}
	var errQuery error
	if txOrder, errQuery = model.QueryOrderByOrderNo(l.svcCtx.MyDB, in.OrderNo, ""); errQuery != nil {
		return &transactionclient.ProxyPayFailResponse{
			Code:    response.DATABASE_FAILURE,
			Message: "查詢訂單資料錯誤，orderNo : " + in.OrderNo,
		}, nil
	}

	//失败单
	txOrder.Status = constants.FAIL
	txOrder.TransAt = types.JsonTime{}.New()

	updateBalance := &types.UpdateBalance{
		MerchantCode:    txOrder.MerchantCode,
		CurrencyCode:    txOrder.CurrencyCode,
		OrderNo:         txOrder.OrderNo,
		MerchantOrderNo: txOrder.MerchantOrderNo,
		OrderType:       txOrder.Type,
		PayTypeCode:     txOrder.PayTypeCode,
		TransferAmount:  txOrder.TransferAmount,
		TransactionType: "4", //異動類型 (1=收款; 2=解凍; 3=沖正;4=出款退回,11=出款 ; 12=凍結)
		BalanceType:     constants.XF_BALANCE,
		CreatedBy:       txOrder.MerchantCode,
		ChannelCode:     txOrder.ChannelCode,
	}

	//调整异动钱包，并更新订单
	if err = l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

		var merchantPtBalanceId int64
		if err = db.Table("mc_merchant_channel_rate").
			Select("merchant_pt_balance_id").
			Where("merchant_code = ? AND channel_pay_types_code = ?", txOrder.MerchantCode, txOrder.ChannelPayTypesCode).
			Find(&merchantPtBalanceId).Error; err != nil {
			logx.WithContext(l.ctx).Errorf("捞取子钱錢包錯誤，商户号:%s，ChannelPayTypesCode:%s，err:%s", txOrder.MerchantCode, txOrder.ChannelPayTypesCode, err.Error())
			return err
		}

		//异动子钱包
		if merchantPtBalanceId > 0 {
			updateBalance.MerPtBalanceId = merchantPtBalanceId
			if _, err = merchantbalanceservice.UpdateXF_Pt_Balance_Deposit(l.ctx, db, updateBalance); err != nil {
				txOrder.RepaymentStatus = constants.REPAYMENT_FAIL
				logx.WithContext(l.ctx).Errorf("商户:%s，更新子钱錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err.Error(), updateBalance)
				return err
			} else {
				txOrder.RepaymentStatus = constants.REPAYMENT_SUCCESS
				logx.WithContext(l.ctx).Infof("下发API提单失败 %s，子錢包退款成功", merchantBalanceRecord.OrderNo)
			}
		}

		if merchantBalanceRecord, err = merchantbalanceservice.UpdateXFBalance_Deposit(db, *updateBalance); err != nil {
			txOrder.RepaymentStatus = constants.REPAYMENT_FAIL
			logx.WithContext(l.ctx).Errorf("商户:%s，更新錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err.Error(), updateBalance)
		} else {
			txOrder.RepaymentStatus = constants.REPAYMENT_SUCCESS
			logx.WithContext(l.ctx).Infof("代付API提单失败 %s，下發錢包退款成功", merchantBalanceRecord.OrderNo)
		}

		if err = db.Scopes(gormx.GetPartition(txOrder.CreatedAt.Format("2006-01-02 15:04:05"), "tx_orders", &types.OrderX{})).Updates(txOrder).Error; err != nil {
			logx.WithContext(l.ctx).Errorf("代付出款失敗(下发餘額)_ %s 更新订单失败: %s", txOrder.OrderNo, err.Error())
			return
		}

		return
	}); err != nil {
		return &transactionclient.ProxyPayFailResponse{
			Code:    response.UPDATE_DATABASE_FAILURE,
			Message: "異動錢包失敗，orderNo : " + in.OrderNo,
		}, nil
	}

	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     txOrder.OrderNo,
			Action:      constants.ACTION_DF_REFUND,
			UserAccount: in.MerchantCode,
			Comment:     "",
		},
	}).Error; err4 != nil {
		logx.WithContext(l.ctx).Errorf("紀錄訂單歷程出錯:%s", err4.Error())
	}

	return &transactionclient.ProxyPayFailResponse{
		Code:         response.API_SUCCESS,
		Message:      "操作成功",
		ProxyOrderNo: txOrder.OrderNo,
	}, nil
}
