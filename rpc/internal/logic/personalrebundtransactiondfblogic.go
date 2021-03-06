package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"gorm.io/gorm"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type PersonalRebundTransactionDFBLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPersonalRebundTransactionDFBLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PersonalRebundTransactionDFBLogic {
	return &PersonalRebundTransactionDFBLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PersonalRebundTransactionDFBLogic) PersonalRebundTransaction_DFB(in *transactionclient.PersonalRebundRequest) (resp *transactionclient.PersonalRebundResponse, err error) {

	merchantBalanceRecord := types.MerchantBalanceRecord{}
	var txOrder = types.Order{}
	if err = l.svcCtx.MyDB.Table("tx_orders").Where("order_no = ?", in.OrderNo).Take(&txOrder).Error; err != nil {
		return &transactionclient.PersonalRebundResponse{
			Code: response.DATA_NOT_FOUND,
			Message: "查无单号资料，orderNo = "+ in.OrderNo,
		}, nil
	}
	//失败单
	txOrder.Status = constants.FAIL
	transactionType := "4"
	if in.Action == "REVERSAL" {
		transactionType = "3"
	}
	updateBalance := &types.UpdateBalance{
		MerchantCode:    txOrder.MerchantCode,
		CurrencyCode:    txOrder.CurrencyCode,
		OrderNo:         txOrder.OrderNo,
		MerchantOrderNo: txOrder.MerchantOrderNo,
		OrderType:       txOrder.Type,
		PayTypeCode:     txOrder.PayTypeCode,
		TransferAmount:  txOrder.TransferAmount,
		TransactionType: transactionType, //異動類型 (1=收款; 2=解凍; 3=沖正;4=出款退回,11=出款 ; 12=凍結)
		BalanceType:     constants.DF_BALANCE,
		CreatedBy:       in.UserAccount,
		Comment:         in.Memo,
		ChannelCode:     txOrder.ChannelCode,
	}
	var jTime types.JsonTime
	//调整异动钱包，并更新订单
	if err = l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

		if merchantBalanceRecord, err = merchantbalanceservice.UpdateDFBalance_Deposit(db, updateBalance); err != nil {
			logx.Errorf("商户:%s，更新錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err.Error(), updateBalance)
			return err
		} else {
			logx.Infof("代付API提单失败 %s，代付錢包退款成功", merchantBalanceRecord.OrderNo)
		}

		if err = db.Table("tx_orders").Updates(&types.OrderX{
			Order:   txOrder,
			TransAt: jTime.New(),
		}).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return &transactionclient.PersonalRebundResponse{
			Code: response.UPDATE_DATABASE_FAILURE,
			Message: "人工还款失败，err : "+ err.Error(),
		}, nil
	}

	if err4 := l.svcCtx.MyDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     txOrder.OrderNo,
			Action:      in.Action,
			UserAccount: in.UserAccount,
			Comment:     in.Memo,
		},
	}).Error; err4 != nil {
		logx.Error("紀錄訂單歷程出錯:%s", err4.Error())
	}

	failResp := &transactionclient.PersonalRebundResponse{
		OrderNo: txOrder.OrderNo,
		Code: response.API_SUCCESS,
		Message: "操作成功",
	}

	return failResp, nil
}
