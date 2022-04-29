package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"gorm.io/gorm"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProxyOrderToTestDFBLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProxyOrderToTestDFBLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyOrderToTestDFBLogic {
	return &ProxyOrderToTestDFBLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProxyOrderToTestDFBLogic) ProxyOrderToTest_DFB(in *transactionclient.ProxyOrderTestRequest) (*transactionclient.ProxyOrderTestResponse, error) {
	txOrder := &types.OrderX{}
	var err error
	if txOrder, err = model.QueryOrderByOrderNo(l.svcCtx.MyDB, in.ProxyOrderNo, ""); err != nil {
		return nil, errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	//如果月結傭金"已結算/確認報表無誤按鈕" : 不扣款
	txOrder.IsTest = "1"

	// 更新订单
	if txOrder != nil {
		if errUpdate := l.svcCtx.MyDB.Table("tx_orders").Updates(txOrder).Error; errUpdate != nil {
			logx.Error("代付订单更新状态错误: ", errUpdate.Error())
		}
	}

	l.svcCtx.MyDB.Transaction(func(db *gorm.DB) (err error) {

		merchantBalanceRecord := types.MerchantBalanceRecord{}

		// 新增收支记录，与更新商户余额(商户账户号是黑名单，把交易金额为设为 0)
		updateBalance := &types.UpdateBalance{
			MerchantCode:    txOrder.MerchantCode,
			CurrencyCode:    txOrder.CurrencyCode,
			OrderNo:         txOrder.OrderNo,
			MerchantOrderNo: txOrder.MerchantOrderNo,
			OrderType:       txOrder.Type,
			PayTypeCode:     txOrder.PayTypeCode,
			PayTypeCodeNum:  txOrder.PayTypeCodeNum,
			TransferAmount:  txOrder.TransferAmount,
			TransactionType: "4", //異動類型 (1=收款 ; 2=解凍;  3=沖正 4=還款;  5=補單; 11=出款 ; 12=凍結 ; 13=追回; 20=調整)
			BalanceType:     constants.DF_BALANCE,
			CreatedBy:       txOrder.MerchantCode,
		}

		if merchantBalanceRecord, err = merchantbalanceservice.UpdateDFBalance_Deposit(db, *updateBalance); err != nil {
			logx.Errorf("商户:%s，更新錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err.Error(), updateBalance)
			return errorz.New(response.SYSTEM_ERROR, err.Error())
		} else {
			logx.Infof("代付API提单 %s，錢包還款成功", merchantBalanceRecord.OrderNo)
			txOrder.BeforeBalance = merchantBalanceRecord.BeforeBalance // 商戶錢包異動紀錄
			txOrder.Balance = merchantBalanceRecord.AfterBalance
		}

		return nil
	})

	return &transactionclient.ProxyOrderTestResponse{}, nil
}
