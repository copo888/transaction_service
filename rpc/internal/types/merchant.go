package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/copo888/transaction_service/common/gormx"
	"time"
)

func (Merchant) TableName() string {
	return "mc_merchants"
}

func (MerchantCurrency) TableName() string {
	return "mc_merchant_currencies"
}

func (MerchantBalance) TableName() string {
	return "mc_merchant_balances"
}

func (o MerchantContact) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (o *MerchantContact) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), o)
}

func (o MerchantBizInfo) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (o *MerchantBizInfo) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), o)
}

type MerchantX struct {
	Merchant
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MerchantQueryListViewRequestX struct {
	MerchantQueryListViewRequest
	Orders []gormx.Sortx `json:"orders, optional" gorm:"-"`
}

type MerchantConfigureRateListRequestX struct {
	MerchantConfigureRateListRequest
	Orders []gormx.Sortx `json:"orders, optional" gorm:"-"`
}

type MerchantQueryRateListViewRequestX struct {
	MerchantQueryRateListViewRequest
	Orders []gormx.Sortx `json:"orders, optional" gorm:"-"`
}

type MerchantQueryAllRequestX struct {
	MerchantQueryAllRequest
	Orders []gormx.Sortx `json:"orders, optional" gorm:"-"`
}

type MerchantCurrencyQueryAllRequestX struct {
	MerchantCurrencyQueryAllRequest
	Orders []gormx.Sortx `json:"orders, optional" gorm:"-"`
}

type MerchantBalanceRecordQueryAllRequestX struct {
	MerchantBalanceRecordQueryAllRequest
	Orders []gormx.Sortx `json:"orders, optional" gorm:"-"`
}

type MerchantCurrencyCreate struct {
	MerchantCurrencyCreateRequest
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MerchantCurrencyUpdate struct {
	MerchantCurrencyUpdateRequest
	CreatedAt time.Time `json:"createdAt, optional"`
	UpdatedAt time.Time `json:"updatedAt, optional"`
}

type MerchantUpdateCurrenciesRequestX struct {
	MerchantUpdateCurrenciesRequest
	Currencies []MerchantCurrencyUpdate `json:"currencies"`
}

type MerchantUpdate2 struct {
	Merchant
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MerchantBalanceCreate struct {
	MerchantBalanceCreateRequest
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MerchantBalanceUpdate struct {
	MerchantBalanceUpdateRequest
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MerchantBalanceX struct {
	MerchantBalance
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MerchantPtBalanceX struct {
	MerchantPtBalance
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MerchantChannelRateConfigure struct {
	MerchantChannelRateConfigureRequest
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MerchantBindBankCreate struct {
	MerchantBindBankCreateRequest
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MerchantBindBankUpdate struct {
	MerchantBindBankUpdateRequest
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MerchantBalanceRecordX struct {
	MerchantBalanceRecord
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MerchantPtBalanceRecordX struct {
	MerchantPtBalanceRecord
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MerchantCommissionRecordX struct {
	MerchantCommissionRecord
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MerchantFrozenRecordX struct {
	MerchantFrozenRecord
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UpdateBalance struct {
	MerchantCode    string
	CurrencyCode    string
	OrderNo         string
	MerchantOrderNo string
	OrderType       string
	ChannelCode     string
	PayTypeCode     string
	PayTypeCodeNum  string
	TransactionType string
	BalanceType     string
	TransferAmount  float64
	Comment         string
	CreatedBy       string
}

type UpdateFrozenAmount struct {
	MerchantCode    string
	CurrencyCode    string
	OrderNo         string
	MerchantOrderNo string
	OrderType       string
	ChannelCode     string
	PayTypeCode     string
	PayTypeCodeNum  string
	TransactionType string
	BalanceType     string
	FrozenAmount    float64
	Comment         string
	CreatedBy       string
}

type FrozenManually struct {
	MerchantCode    string
	CurrencyCode    string
	OrderNo         string
	OrderType       string
	TransactionType string
	BalanceType     string
	FrozenAmount    float64
	Comment         string
	CreatedBy       string
}

type CorrespondMerChnRate struct {
	MerchantCode        string  `json:"merchantCode"`
	ChannelPayTypesCode string  `json:"channelPayTypesCode"`
	ChannelCode         string  `json:"channelCode"`
	PayTypeCode         string  `json:"payTypeCode"`
	PayTypeCodeNum      string  `json:"payTypeCodeNum"`
	Designation         string  `json:"designation"`
	DesignationNo       string  `json:"designationNo"`
	Fee                 float64 `json:"fee"`
	HandlingFee         float64 `json:"handlingFee"`
	ChFee               float64 `json:"chFee"`
	ChHandlingFee       float64 `json:"chHandlingFee"`
	SingleMinCharge     float64 `json:"singleMinCharge"`
	SingleMaxCharge     float64 `json:"singleMaxCharge"`
	CurrencyCode        string  `json:"currencyCode"`
	ApiUrl              string  `json:"apiUrl"`
}
