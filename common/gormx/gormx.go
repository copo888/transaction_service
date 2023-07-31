package gormx

import (
	"github.com/copo888/transaction_service/common/utils"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

func GetPartition(time string, tableName string, model interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		suffix := strings.ReplaceAll(time[:7], "-", "")

		if err := db.Table(tableName + "_" + suffix).AutoMigrate(model); err != nil {
			panic("Failed to auto migrate table: " + err.Error())
		}
		return db
	}
}

func QueryByPartition(req interface{}, tableName string, model interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 本月搜寻 七月 startAt "2023-05-30 16:00:00"，endAt "2023-07-31 15:59:59"
		duration := reflect.ValueOf(req)
		startAt := duration.FieldByName("StartAt").Interface().(string)
		endAt := utils.ParseTimeAddOneSecond(duration.FieldByName("EndAt").Interface().(string))

		timeStart, _ := time.Parse("2006-01-02 15:04:05", startAt)
		timeEnd, _ := time.Parse("2006-01-02 15:04:05", endAt)

		if timeStart.Year() == timeEnd.Year() && timeStart.Month() == timeEnd.Month() {
			suffix_start := strings.ReplaceAll(startAt[:7], "-", "")
			//suffix_end := strings.ReplaceAll(endAt[:7], "-", "")
			return db.Table(tableName+"_"+suffix_start).Where("`created_at` >= ? AND `created_at` < ?", startAt, endAt)

		} else {

			subQueryStart := db.Table(tableName+"_"+strings.ReplaceAll(startAt[:7], "-", "")).
				Session(&gorm.Session{DryRun: true}).
				Where("`created_at` >= ? AND `created_at` < ?", startAt, getLastDayTime(startAt)).
				Find(model).Statement

			subQueryEnd := db.Table(tableName+"_"+strings.ReplaceAll(endAt[:7], "-", "")).
				Session(&gorm.Session{DryRun: true}).
				Where("`created_at` >= ? AND `created_at` < ?", getFirstDayTime(endAt), endAt).
				Find(model).Statement

			query := "(" + subQueryStart.SQL.String() + ") UNION (" + subQueryEnd.SQL.String() + ")"
			args := append(subQueryStart.Vars, subQueryEnd.Vars...)

			iteratorTime := timeStart.AddDate(0, 1, 0)
			for iteratorTime.Year() != timeEnd.Year() || iteratorTime.Month() != timeEnd.Month() {

				subQueryInterVal := db.Table(tableName+"_"+strings.ReplaceAll(iteratorTime.Format("2006-01-02 15:04:05")[:7], "-", "")).
					Session(&gorm.Session{DryRun: true}).
					Where("`created_at` >= ? AND `created_at` < ?", getFirstDayTime(iteratorTime.Format("2006-01-02 15:04:05")), getLastDayTime(iteratorTime.Format("2006-01-02 15:04:05"))).
					Find(model).Statement

				query = query + " UNION (" + subQueryInterVal.SQL.String() + ")"
				args = append(args, subQueryInterVal.Vars...)

				iteratorTime = iteratorTime.AddDate(0, 1, 0)
			}

			return db.Raw(query, args...)
		}
	}
}

//取得當月最後一天秒
func getLastDayTime(dateTime string) string {

	dateTimeParse, _ := time.Parse("2006-01-02 15:04:05", dateTime)

	nextMonth := dateTimeParse.AddDate(0, 1, 0)
	firstDayOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())

	// 减去 1 秒，得到当月的最后一天的最后一秒
	lastDayToSec := firstDayOfNextMonth.Add(-time.Second)
	return lastDayToSec.Format("2006-01-02 15:04:05")

}

func getFirstDayTime(dateTime string) string {

	dateTimeParse, _ := time.Parse("2006-01-02 15:04:05", dateTime)

	// 获取这个月的第一天
	firstDayOfMonth := time.Date(dateTimeParse.Year(), dateTimeParse.Month(), 1, 0, 0, 0, 0, dateTimeParse.Location())

	return firstDayOfMonth.Format("2006-01-02 15:04:05")

}

//分页功能
func Paginate(page interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := reflect.ValueOf(page)
		pageNum := page.FieldByName("PageNum").Interface().(int) - 1
		pageSize := page.FieldByName("PageSize").Interface().(int)
		offset := pageNum * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

type Sortx struct {
	Column string `json:"column, optional" gorm:"-"`
	Asc    bool   `json:"asc, optional" gorm:"-"`
}

func Sort(sorts []Sortx) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if sorts == nil {
			return db
		}

		orderStr := ""
		for i := range sorts {
			orderStr += sorts[i].Column
			if sorts[i].Asc {
				orderStr += " asc"
			} else {
				orderStr += " desc"
			}
			if i+1 < len(sorts) {
				orderStr += ", "
			}
		}
		return db.Order(orderStr)
	}
}
