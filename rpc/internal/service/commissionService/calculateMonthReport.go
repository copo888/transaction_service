package commissionService

import (
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// CalculateMonthReport 計算當月傭金報表(單筆代理報表)
func CalculateMonthReport(db *gorm.DB, report types.CommissionMonthReportX, startAt, endAt string) error {

	payTotalAmount := 0.0
	payCommission := 0.0
	internalChargeTotalAmount := 0.0
	internalChargeCommission := 0.0
	proxyPayTotalNumber := 0.0
	proxyPayCommission := 0.0
	totalCommission := 0.0

	// 計算報表 支付 資料
	zfDetails, err := calculateMonthReportDetails(db, report, "ZF", startAt, endAt)
	if err != nil {
		logx.Errorf("計算傭金報表 支付詳情 失敗: %#v, error: %s", report, err.Error())
		return err
	}

	// 計算報表 內充 資料
	ncDetails, err := calculateMonthReportDetails(db, report, "NC", startAt, endAt)
	if err != nil {
		logx.Errorf("計算傭金報表 內充詳情 失敗: %#v, error: %s", report, err.Error())
		return err
	}

	// 計算報表 代付 資料
	dfDetails, err := calculateMonthReportDetails(db, report, "DF", startAt, endAt)
	if err != nil {
		logx.Errorf("計算傭金報表 代付詳情 失敗: %#v, error: %s", report, err.Error())
		return err
	}

	// 保存 並 計算支付總額
	for _, detail := range zfDetails {
		payTotalAmount += detail.TotalAmount
		payCommission += detail.TotalCommission
		totalCommission += detail.TotalCommission
		detail.CommissionMonthReportId = report.ID
		detail.OrderType = "ZF"
		if err = db.Table("cm_commission_month_report_details").Create(&detail).Error; err != nil {
			logx.Errorf("保存傭金報表支付詳情 失敗: %#v, error: %s", detail, err.Error())
			return errorz.New(response.DATABASE_FAILURE)
		}
	}

	// 保存 並 計算內充總額
	for _, detail := range ncDetails {
		internalChargeTotalAmount += detail.TotalAmount
		internalChargeCommission += detail.TotalCommission
		totalCommission += detail.TotalCommission
		detail.CommissionMonthReportId = report.ID
		detail.OrderType = "NC"
		if err = db.Table("cm_commission_month_report_details").Create(&detail).Error; err != nil {
			logx.Errorf("保存傭金報表內充詳情 失敗: %#v, error: %s", detail, err.Error())
			return errorz.New(response.DATABASE_FAILURE)
		}
	}

	// 保存 並 計算代付總額
	for _, detail := range dfDetails {
		proxyPayTotalNumber += detail.TotalNumber
		proxyPayCommission += detail.TotalCommission
		totalCommission += detail.TotalCommission
		detail.CommissionMonthReportId = report.ID
		detail.OrderType = "DF"
		if err = db.Table("cm_commission_month_report_details").Create(&detail).Error; err != nil {
			logx.Errorf("保存傭金報表代付詳情 失敗: %#v, error: %s", detail, err.Error())
			return errorz.New(response.DATABASE_FAILURE)
		}
	}

	// 編輯主表
	report.PayTotalAmount = payTotalAmount
	report.PayCommission = payCommission
	report.InternalChargeTotalAmount = internalChargeTotalAmount
	report.InternalChargeCommission = internalChargeCommission
	report.ProxyPayTotalNumber = proxyPayTotalNumber
	report.ProxyPayCommission = proxyPayCommission
	report.TotalCommission = totalCommission
	if err = db.Table("cm_commission_month_reports").Updates(&report).Error; err != nil {
		logx.Errorf("編輯傭金報表失敗: %#v, error: %s", report, err.Error())
		return errorz.New(response.DATABASE_FAILURE)
	}

	return nil
}

func calculateMonthReportDetails(db *gorm.DB, report types.CommissionMonthReportX, orderType, startAt, endAt string) ([]types.CommissionMonthReportDetailX, error) {
	var reportDetails []types.CommissionMonthReportDetailX

	selectX := "o.merchant_code, " +
		"o.currency_code," +
		"m.fee as merchant_fee," +
		"p.fee as agent_fee," +
		"m.fee - p.fee as diff_fee," +
		"m.handling_fee as merchant_handling_fee," +
		"p.handling_fee as agent_handling_fee," +
		"m.handling_fee - p.handling_fee as diff_handling_fee," +
		"count(o.order_no) as total_number," +
		"sum(p.profit_amount) as total_commission,"

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

	err := db.Table("tx_orders_fee_profit m"). //上層商戶
							Select(selectX).
							Joins("JOIN tx_orders_fee_profit p on p.merchant_code = m.agent_parent_code and p.order_no = m.order_no"). //代理商戶
							Joins("JOIN tx_orders o on o.order_no = m.order_no").                                                      // 訂單
							Where("o.trans_at >= ? and o.trans_at < ? ", startAt, endAt).
							Where("p.merchant_code = ? ", report.MerchantCode).
							Where("o.currency_code = ? ", report.CurrencyCode).
							Where("o.type = ? ", orderType).
							Where("p.profit_amount != 0 ").
							Where("(o.status = 20 || o.status = 31) ").
							Group("merchant_code, currency_code, pay_type_code, merchant_fee, agent_fee").
							Find(&reportDetails).Error

	return reportDetails, err
}
