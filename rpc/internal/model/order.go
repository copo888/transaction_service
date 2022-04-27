package model

import (
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/random"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/types"
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

/*
	@param orderNo    : copo訂單號
    @param merOrderNo : 商戶訂單號
*/
func QueryOrderByOrderNo(db *gorm.DB, orderNo string, merOrderNo string) (*types.OrderX, error) {
	txOrder := &types.OrderX{}
	if orderNo != "" || len(orderNo) > 0 {
		if err := db.Table("tx_orders").Where("order_no = ?", orderNo).Find(&txOrder).Error; err != nil {
			return nil, errorz.New(response.DATABASE_FAILURE, err.Error())
		}
	} else if merOrderNo != "" || len(merOrderNo) > 0 {
		if err := db.Table("tx_orders").Where("merchant_order_no = ?", orderNo, merOrderNo).Find(&txOrder).Error; err != nil {
			return nil, errorz.New(response.DATABASE_FAILURE, err.Error())
		}
	} else {
		return nil, errorz.New(response.DATABASE_FAILURE)
	}

	return txOrder, nil
}
