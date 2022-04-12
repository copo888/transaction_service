package gormx

import (
	"gorm.io/gorm"
	"reflect"
)

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
