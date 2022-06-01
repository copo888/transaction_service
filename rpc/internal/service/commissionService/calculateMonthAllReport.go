package commissionService

import (
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

// CalculateMonthAllReport 計算當月傭金報表
func CalculateMonthAllReport(db *gorm.DB, month string) error {
	monthArray := strings.Split(month, "-")

	// 檢查月份格式
	if len(monthArray) != 2 {
		// todo: 時間格是錯誤
	}
	y, err1 := strconv.Atoi(monthArray[0])
	m, err2 := strconv.Atoi(monthArray[1])
	if err1 != nil || err2 != nil {
		// todo: 時間格是錯誤
	}

	// 取得此月份起訖時間
	startAt := BeginningOfMonth(y, m).Format("2006-01-02 15:04:05")
	endAt := EndOfMonth(y, m).Format("2006-01-02 15:04:05")

	// 取得此月份所有要計算的代理商戶
	reports, err := getAllMonthReports(db, startAt, endAt)
	if err != nil {
		return errorz.New(response.DATABASE_FAILURE)
	}
	logx.Infof("開始計算 %s 月份傭金 總共有 %d 筆 Transaction start", month, len(reports))
	if errTx := db.Transaction(func(txdb *gorm.DB) (err error) {
		// 迴圈計算 單筆代理傭金報表
		for _, report := range reports {
			report.Month = month
			report.Status = "0"
			// 保存
			if err = txdb.Table("cm_commission_month_reports").Create(&report).Error; err != nil {
				logx.Errorf("建立傭金報表失敗: %#v, error: %s", report, err.Error())
				return errorz.New(response.DATABASE_FAILURE)
			}
			// 計算報表詳情
			if err:= CalculateMonthReport(txdb, report, startAt, endAt); err != nil {
				return err
			}
		}
		return
	}); errTx != nil {
		return errTx
	}
	logx.Infof("完成計算 %s 月份傭金 Transaction end", month)

	return nil
}


func createMonthReport(db *gorm.DB, month string, report *types.CommissionMonthReportX) error {
	return db.Table("cm_commission_month_reports").Create(report).Error
}

func getAllMonthReports(db *gorm.DB, startAt, endAt string) ([]types.CommissionMonthReportX, error) {
	var commissionMonthReports []types.CommissionMonthReportX

	selectX := "p.agent_layer_no, " +
		"p.merchant_code, " +
		"o.currency_code "

	err := db.Table("tx_orders_fee_profit m").
		Select(selectX).
		Joins("JOIN tx_orders_fee_profit p on p.merchant_code = m.agent_parent_code and p.order_no = m.order_no").
		Joins("JOIN tx_orders o on o.order_no = m.order_no").
		Where("o.trans_at >= ? and o.trans_at < ? ", startAt, endAt).
		Where("(o.status = 20 || o.status = 31) ").
		Where("o.is_test != 1 ").
		Distinct().Find(&commissionMonthReports).Error

	return commissionMonthReports, err
}

// BeginningOfMonth 取得月開始時間
func BeginningOfMonth(year, month int) time.Time {
	// UTC +8 要扣回
	return time.Date(year, time.Month(month), 1, -8, 0, 0, 0, time.UTC)
}

// EndOfMonth 取得月結束時間 (下個月1號)
func EndOfMonth(year, month int) time.Time {
	return BeginningOfMonth(year, month).AddDate(0, 1, 0)
}
