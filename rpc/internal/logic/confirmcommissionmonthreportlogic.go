package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmCommissionMonthReportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmCommissionMonthReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmCommissionMonthReportLogic {
	return &ConfirmCommissionMonthReportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConfirmCommissionMonthReportLogic) ConfirmCommissionMonthReport(in *transactionclient.ConfirmCommissionMonthReportRequest) (*transactionclient.ConfirmCommissionMonthReportResponse, error) {
	var report types.CommissionMonthReportX

	/****     交易開始      ****/
	txDB := l.svcCtx.MyDB.Begin()
	// 取得報表
	if err := txDB.Table("cm_commission_month_reports").Where("id = ?", in.ID).Find(&report).Error; err != nil {
		txDB.Rollback()
		return &transactionclient.ConfirmCommissionMonthReportResponse{
			Code:    response.DATABASE_FAILURE,
			Message: err.Error(),
		}, nil
	}

	if report.Status == "1" {
		// 已審核報表不可再異動
		txDB.Rollback()
		return &transactionclient.ConfirmCommissionMonthReportResponse{
			Code:    response.MERCHANT_COMMISSION_AUDIT,
			Message: "佣金资料已审核，不能异动资料",
		}, nil
	}

	// 將報表改為已審核
	report.ConfirmBy = in.ConfirmBy
	report.ConfirmAt = types.JsonTime{}.New()
	report.Status = "1"

	if err := txDB.Table("cm_commission_month_reports").Updates(&report).Error; err != nil {
		txDB.Rollback()
		logx.Errorf("確認佣金報表update失败，AugentLayerNo: %s, Month: %s, err : %s", report.AgentLayerNo, report.Month, err.Error())
		return &transactionclient.ConfirmCommissionMonthReportResponse{
			Code:    response.MERCHANT_COMMISSION_AUDIT,
			Message: err.Error(),
		}, nil
	}

	transferAmount := 0.0
	if report.ChangeCommission > 0 {
		transferAmount = report.ChangeCommission
	} else {
		transferAmount = report.TotalCommission
	}
	if _, err := merchantbalanceservice.UpdateCommissionAmount(txDB, types.UpdateCommissionAmount{
		MerchantCode:            report.MerchantCode,
		CurrencyCode:            report.CurrencyCode,
		CommissionMonthReportId: report.ID,
		TransactionType:         constants.COMMISSION_TRANSACTION_TYPE_MONTHLY,
		TransferAmount:          transferAmount,
		Comment:                 "佣金月結:" + report.Month,
		CreatedBy:               in.ConfirmBy,
	}); err != nil {
		txDB.Rollback()
		return &transactionclient.ConfirmCommissionMonthReportResponse{
			Code:    response.SYSTEM_ERROR,
			Message: "更新錢包失敗",
		}, nil
	}

	if err := txDB.Commit().Error; err != nil {
		txDB.Rollback()
		logx.Errorf("確認佣金報表Commit失败，AgentLayerNo: %s, Month: %s, err : %s", report.AgentLayerNo, report.Month, err.Error())
		return &transactionclient.ConfirmCommissionMonthReportResponse{
			Code:    response.DATABASE_FAILURE,
			Message: "资料库错误 Commit失败",
		}, nil
	}
	/****     交易結束      ****/

	return &transactionclient.ConfirmCommissionMonthReportResponse{
		Code:    response.API_SUCCESS,
		Message: "操作成功",
	}, nil
}
