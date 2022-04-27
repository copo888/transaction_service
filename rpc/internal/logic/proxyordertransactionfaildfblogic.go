package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"gorm.io/gorm"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProxyOrderTransactionFailDFBLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProxyOrderTransactionFailDFBLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyOrderTransactionFailDFBLogic {
	return &ProxyOrderTransactionFailDFBLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

/*
	代付出款失敗(代付餘額)_還款
	更新餘額、餘額異動紀錄、更新訂單(失敗單)、異動訂單紀錄
*/
func (l *ProxyOrderTransactionFailDFBLogic) ProxyOrderTransactionFail_DFB(in *transactionclient.ProxyPayFailRequest) (resp *transactionclient.ProxyPayFailResponse, err error) {

	merchantBalanceRecord := types.MerchantBalanceRecord{}
	var txOrder = &types.OrderX{}
	var errQuery error
	if txOrder, errQuery = model.QueryOrderByOrderNo(l.svcCtx.MyDB, in.OrderNo, ""); err != nil {
		return nil, errQuery
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
		PayTypeCodeNum:  txOrder.PayTypeCodeNum,
		TransferAmount:  txOrder.TransferAmount,
		TransactionType: "4", //異動類型 (1=收款; 2=解凍; 3=沖正;4=出款退回,11=出款 ; 12=凍結)
		BalanceType:     constants.DF_BALANCE,
		CreatedBy:       txOrder.MerchantCode,
	}

	//调整异动钱包，并更新订单
	if err = l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

		if merchantBalanceRecord, err = merchantbalanceservice.UpdateDFBalance_Deposit(db, *updateBalance); err != nil {
			txOrder.RepaymentStatus = constants.REPAYMENT_FAIL
			logx.Errorf("商户:%s，更新錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err.Error(), updateBalance)
			return
		} else {
			logx.Infof("代付API提单失败 %s，代付錢包退款成功", merchantBalanceRecord.OrderNo)
			txOrder.RepaymentStatus = constants.REPAYMENT_SUCCESS
		}

		if err = db.Table("tx_orders").Updates(txOrder).Error; err != nil {
			return
		}
		return
	}); err != nil {
		return
	}

	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     txOrder.OrderNo,
			Action:      constants.ACTION_DF_REFUND,
			UserAccount: in.MerchantCode,
			Comment:     "",
		},
	}).Error; err4 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err4.Error())
	}

	failResp := &transactionclient.ProxyPayFailResponse{
		ProxyOrderNo: txOrder.OrderNo,
	}

	return failResp, nil
}
