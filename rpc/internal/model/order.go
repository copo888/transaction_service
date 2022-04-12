package model

import (
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/random"
	"github.com/copo888/transaction_service/common/response"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	MyDB  *gorm.DB
	Table string
}

func NewOrder(mydb *gorm.DB, t ...string) *Order {
	table := "tx_orders"
	if len(t) > 0 {
		table = t[0]
	}
	return &Order{
		MyDB:  mydb,
		Table: table,
	}
}

func (m *Order) IsExistByMerchantOrderNo(merchantCode, merchantOrderNo string) (isExist bool, err error) {
	if err = m.MyDB.Table(m.Table).
		Select("count(*) > 0").
		Where("merchant_code = ? AND merchant_order_no = ?", merchantCode, merchantOrderNo).
		Find(&isExist).Error; err != nil {
		err = errorz.New(response.DATABASE_FAILURE, err.Error())
	}
	return
}

// 生成訂單號代付 DF 支付 ZF 下發 XF 內充 NC
func GenerateOrderNo(orderType string) string {
	var result string
	t := time.Now().Format("20060102150405")
	randomStr := random.GetRandomString(5, random.ALL, random.MIX)
	result = orderType + t + randomStr
	return result
}
