package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/copo888/transaction_service/common/gormx"
	"mime/multipart"
	"time"
)

func (ChannelPayType) TableName() string {
	return "ch_channel_pay_types"
}

func (ChannelData) TableName() string {
	return "ch_channels"
}

func (PayType) TableName() string {
	return "ch_pay_types"
}

type ChannelDataCreate struct {
	ChannelDataCreateRequest
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChannelDataUpdate struct {
	ChannelDataUpdateRequest
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChannelDataQueryAllRequestX struct {
	ChannelDataQueryAllRequest
	Orders []gormx.Sortx `json:"orders, optional" gorm:"-"`
}

type JSON json.RawMessage

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

func (o PayTypeMap) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (o *PayTypeMap) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), o)
}

//func (o *Banks) Value() (driver.Value, error) {
//	b, err := json.Marshal(o)
//	return string(b), err
//}
//
//func (o *Banks) Scan(input interface{}) error {
//	return json.Unmarshal(input.([]byte), o)
//}
type PayTypeQueryAllRequestX struct {
	PayTypeQueryAllRequest
	Orders []gormx.Sortx `json:"orders, optional" gorm:"-"`
}

func (o BankCodeMap) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (o *BankCodeMap) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), o)
}

type PayTypeCreate struct {
	PayTypeCreateRequest
	UploadFile   multipart.File        `json:"uploadFile, optional" gorm:"-"`
	UploadHeader *multipart.FileHeader `json:"uploadHeader, optional" gorm:"-"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type PayTypeUpdate struct {
	PayTypeUpdateRequest
	UploadFile   multipart.File        `json:"uploadFile, optional" gorm:"-"`
	UploadHeader *multipart.FileHeader `json:"uploadHeader, optional" gorm:"-"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ChannelPayTypeCreate struct {
	ChannelPayTypeCreateRequest
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChannelPayTypeUpdate struct {
	ChannelPayTypeUpdateRequest
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChannelPayTypeQueryAllRequestX struct {
	ChannelPayTypeQueryAllRequest
	Orders []gormx.Sortx `json:"orders, optional" gorm:"-"`
}

type ChannelBankX struct {
	ChannelCode string `json:"channel_code"`
	BankNo      string `json:"bank_no"`
	BankName    string `json:"bank_name"`
}
