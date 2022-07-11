package logic

import (
	"context"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/service/commissionService"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"gorm.io/gorm"
	"strconv"
	"strings"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CalculateMonthProfitReportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCalculateMonthProfitReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CalculateMonthProfitReportLogic {
	return &CalculateMonthProfitReportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CalculateMonthProfitReportLogic) CalculateMonthProfitReport(in *transactionclient.CalculateMonthProfitReportRequest) (*transactionclient.CalculateMonthProfitReportResponse, error) {
	monthArray := strings.Split(in.Month, "-")

	// 檢查月份格式
	if len(monthArray) != 2 {
		return nil,errorz.New(response.DATABASE_FAILURE)
	}
	y, err1 := strconv.Atoi(monthArray[0])
	m, err2 := strconv.Atoi(monthArray[1])
	if err1 != nil || err2 != nil {
		return nil, errorz.New(response.DATABASE_FAILURE)
	}

	// 上月
	if m == 1 {
		m = 12
		y -= 1
	}else {
		m -= 1
	}
	y2 := strconv.Itoa(y)
	m2 := strconv.Itoa(m)
	preMonth := y2 +"-"+ m2

	// 取得此月份起訖時間
	startAt := commissionService.BeginningOfMonth(y, m).Format("2006-01-02 15:04:05")
	endAt := commissionService.EndOfMonth(y, m).Format("2006-01-02 15:04:05")

	reports, err := getAllMonthReports(l.svcCtx.MyDB, startAt, endAt)
	if err != nil {
		return nil, errorz.New(response.DATABASE_FAILURE)
	}
	if errTx := l.svcCtx.MyDB.Transaction(func(tx *gorm.DB) error {
		for _, report := range reports {
			if err := l.calculateMonthProfitReport(tx, report, startAt, endAt, in.Month, preMonth); err != nil {
				return err
			}
		}
		return nil
	}); errTx != nil {
		return &transactionclient.CalculateMonthProfitReportResponse{
			Code: response.SYSTEM_ERROR,
			Message: errTx.Error(),
		}, nil
	}

	return &transactionclient.CalculateMonthProfitReportResponse{
		Code: response.API_SUCCESS,
		Message: "操作成功",
	}, nil
}

func (l *CalculateMonthProfitReportLogic) calculateMonthProfitReport(db *gorm.DB, report types.CaculateMonthProfitReport, startAt, endAt, month, preMonth string) error {

	// 计算支付资料
	zfDetail, err := l.calculateMonthProfitReportDetails(db, "ZF", startAt, endAt, report.CurrencyCode)
	if err != nil {
		logx.Errorf("計算收益報表 支付資料 失敗: %#v, error: %s", zfDetail, err.Error())
		return err
	}

	// 计算内充资料
	ncDetail, err := l.calculateMonthProfitReportDetails(db, "NC", startAt, endAt, report.CurrencyCode)
	if err != nil {
		logx.Errorf("計算收益報表 內充資料 失敗: %#v, error: %s", ncDetail, err.Error())
		return err
	}

	// 计算代付资料
	dfDetail, err := l.calculateMonthProfitReportDetails(db, "DF", startAt, endAt, report.CurrencyCode)
	if err != nil {
		logx.Errorf("計算收益報表 代付資料 失敗: %#v, error: %s", dfDetail, err.Error())
		return err
	}

	// 计算下发资料
	wfDetail, err := l.calculateMonthProfitReportDetails(db, "DF", startAt, endAt, report.CurrencyCode)
	if err != nil {
		logx.Errorf("計算收益報表 下發資料 失敗: %#v, error: %s", wfDetail, err.Error())
		return err
	}

	receivedTotalNetProfit := 0.0
	remitTotalNetProfit := 0.0
	totalNetProfit := 0.0
	profitGrowthRate := 0.0

	receivedTotalNetProfit = utils.FloatAdd(zfDetail.TotalProfit, ncDetail.TotalProfit)
	remitTotalNetProfit = utils.FloatAdd(dfDetail.TotalProfit, wfDetail.TotalProfit)

	// 計算佣金資料
	commissionTotalAmount, err := l.calculateCommissionMonthData(db, month, report.CurrencyCode)
	if err != nil {
		logx.Errorf("計算傭金總額 失敗: error: %s", err.Error())
		return err
	}

	// 取得上個月收益資料，計算成長率
	var oldIncomReport *types.IncomReport
	if err := db.Table("re_incom_repot").
		Where("month = ? AND currency_code = ?",preMonth, report.CurrencyCode).
		Find(oldIncomReport).Error; err != nil {
		logx.Errorf("查詢上月收益報表失敗: error: %s", err.Error())
		return errorz.New(response.DATABASE_FAILURE)
	}
	if oldIncomReport != nil {
		// 盈利成長率 = (當月總盈利-上月總盈利)/上月總盈利*100
		profitGrowthRate = utils.FloatMul(utils.FloatDiv(utils.FloatSub(totalNetProfit, oldIncomReport.TotalNetProfit),oldIncomReport.TotalNetProfit), 100)
	}


	var incomReportX types.IncomReportX
	incomReportX.Month = month
	incomReportX.CurrencyCode = report.CurrencyCode
	incomReportX.PayTotalAmount = zfDetail.TotalAmount
	incomReportX.PayNetProfit = zfDetail.TotalProfit
	incomReportX.InternalChargeTotalAmount = ncDetail.TotalAmount
	incomReportX.InternalChargeNetProfit = ncDetail.TotalProfit
	incomReportX.WithdrawTotalAmount = wfDetail.TotalAmount
	incomReportX.WithdrawNetProfit = wfDetail.TotalProfit
	incomReportX.ProxyPayTotalAmount = dfDetail.TotalAmount
	incomReportX.ProxyPayNetProfit = dfDetail.TotalProfit
	incomReportX.ReceivedTotalNetProfit = receivedTotalNetProfit
	incomReportX.RemitTotalNetProfit = remitTotalNetProfit
	incomReportX.TotalNetProfit = totalNetProfit
	incomReportX.CommissionTotalAmount = commissionTotalAmount
	incomReportX.ProfitGrowthRate = profitGrowthRate

	if err := db.Table("rp_incom_report").Create(&incomReportX).Error; err != nil {
		logx.Errorf("新增收益報表失敗: %#v, error: %s", incomReportX, err.Error())
		return errorz.New(response.DATABASE_FAILURE)
	}
	return nil
}

func (l *CalculateMonthProfitReportLogic) calculateMonthProfitReportDetails(db *gorm.DB, orderType, startAt, endAt, currencyCode string) ( *types.CaculateMonthProfitReport, error) {
 	var caculateMonthProfitReport types.CaculateMonthProfitReport
	selectX := "m.merchant_code, " +
		"o.curency_code, " +
		"sum( m.profit_amount ) AS total_profit, "

	if orderType == "NC" {
		// 內充要以 orderType 替代 pay_type_code
		selectX += " 'NC' as pay_type_code,"
	} else {
		selectX += "o.pay_type_code as pay_type_code,"
	}

	if orderType == "ZF" {
		// 支付 使用實際付款金額
		selectX += "sum(o.actual_amount) as total_amount"
	} else {
		// 內充 代付 使用訂單金額
		selectX += "sum(o.order_amount) as total_amount"
	}

	err := db.Table("tx_orders_fee_profit m").
		Select(selectX).
		Joins("JOIN tx_orders o ON o.order_no = m.order_no").
		Where("o.trans_at >= ? and o.trans_at < ?",startAt, endAt).
		Where("m.merchant_code = '00000000'").
		Where("o.currency_code = ?", currencyCode).
		Where("o.type = ?", orderType).
		Where("(o.status = 20)").
		Where("o.is_test != 1").
		Group("merchant_code, currency_code, pay_type_code").
		Find(&caculateMonthProfitReport).Error

	return &caculateMonthProfitReport, err
}

func getAllMonthReports(db *gorm.DB, startAt, endAt string) ([]types.CaculateMonthProfitReport, error) {
	var resp []types.CaculateMonthProfitReport
	selectX := "m.merchant_code, " +
		"o.currency_code "

	err := db.Table("tx_orders_fee_profit m").
		Select(selectX).
		Joins("JOIN tx_orders o on o.order_no = m.order_no").
		Where("o.trans_at >= ? and o.trans_at < ? ", startAt, endAt).
		Where("(o.status = 20) ").
		Where("o.is_test != 1 ").
		Where("m.merchant_code = '00000000'").
		Group("currency_code").
		Distinct().Find(&resp).Error

		return resp, err
}

func (l *CalculateMonthProfitReportLogic) calculateCommissionMonthData (db *gorm.DB, month, currencyCode string) (float64, error){
	var commissionMonthReports []types.CommissionMonthReport
	if err := db.Table("cm_commission_month_reports").
		Where("month = ? AND currency_code = ? And ", month, currencyCode).Find(&commissionMonthReports).Error; err != nil {
		return 0.0, errorz.New(response.DATABASE_FAILURE)
	}

	resp := 0.0
	for _, report := range commissionMonthReports {
		if report.ChangeCommission != 0 {
			resp = utils.FloatAdd(resp, report.ChangeCommission)
		}else{
			resp = utils.FloatAdd(resp, report.TotalCommission)
		}
	}

	return resp, nil
}