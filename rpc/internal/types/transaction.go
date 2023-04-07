package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/copo888/transaction_service/common/gormx"
	"mime/multipart"
	"strings"
	"time"
	"unsafe"
)

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	var res string
	if !time.Time(j).IsZero() {
		var stamp = fmt.Sprintf("%s", time.Time(j).Format("2006-01-02 15:04:05"))
		str := strings.Split(stamp, " +")
		res = str[0]
		return json.Marshal(res)
	}
	return json.Marshal("")
}

func (j JsonTime) Time() time.Time {
	return time.Time(j)
}

func (j JsonTime) Value() (driver.Value, error) {
	return time.Time(j), nil
}

func (j JsonTime) Parse(s string, zone ...string) (JsonTime, error) {

	var (
		loc *time.Location
		err error
	)
	if len(zone) > 0 {
		loc, err = time.LoadLocation(zone[0])
	} else {
		loc, err = time.LoadLocation("")
	}

	if err != nil {
		return j, err
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", s, loc)
	if err != nil {
		return j, err
	}
	jt := (*JsonTime)(unsafe.Pointer(&t))

	return *jt, nil
}

func (j JsonTime) New(ts ...time.Time) JsonTime {
	var t time.Time

	if len(ts) > 0 {
		t = ts[0]
	} else {
		t = time.Now().UTC()
	}

	jt := (*JsonTime)(unsafe.Pointer(&t))
	return *jt
}

func (OrderFeeProfit) TableName() string {
	return "tx_orders_fee_profit"
}

func (OrderChannels) TableName() string {
	return "tx_order_channels"
}

type OrderX struct {
	Order
	TransAt   JsonTime `json:"transAt, optional"`
	FrozenAt  JsonTime `json:"frozenAt, optional"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderInternalCreate struct {
	OrderX
	FormData map[string][]*multipart.FileHeader `gorm:"-"`
}

type OrderInternalUpdate struct {
	OrderX
}

type OrderWithdrawUpdate struct {
	OrderX
}

type UploadImageRequestX struct {
	UploadImageRequest
	UploadFile   multipart.File
	UploadHeader *multipart.FileHeader
	Files        map[string][]*multipart.FileHeader
}

type MerchantRateListViewX struct {
	MerchantRateListView
	Balance float64 `json:"balance"`
}

type MerchantOrderRateListViewX struct {
	MerchantOrderRateListView
	Balance float64 `json:"balance"`
}

type OrderQueryMerchantCurrencyAndBanksResponseX struct {
	MerchantOrderRateListViewX *MerchantOrderRateListViewX `json:"merchantOrderRateListViewX"`
	ChannelBanks               []ChannelBankX
}

type ReceiptRecordQueryAllRequestX struct {
	ReceiptRecordQueryAllRequest
	Orders []gormx.Sortx `json:"orders, optional" gorm:"-"`
}

type FrozenRecordQueryAllRequestX struct {
	FrozenRecordQueryAllRequest
	Orders []gormx.Sortx `json:"orders, optional" gorm:"-"`
}

type DeductRecordQueryAllRequestX struct {
	DeductRecordQueryAllRequest
	Orders []gormx.Sortx `json:"orders, optional" gorm:"-"`
}

type WithdrawOrderUpdateRequestX struct {
	List []ChannelWithdraw `json:"list"`
	OrderX
}

type OrderActionX struct {
	OrderAction
	CreatedAt time.Time
}

type OrderChannelsX struct {
	OrderChannels
	CreatedAt time.Time
	UpdatedAt time.Time
}
type OrderFeeProfitX struct {
	OrderFeeProfit
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ReceiptRecordX struct {
	ReceiptRecord
	TransAt JsonTime `json:"transAt, optional"`
}

type FrozenRecordX struct {
	FrozenRecord
	TransAt JsonTime `json:"transAt, optional"`
}

type DeductRecordX struct {
	DeductRecord
	TransAt JsonTime `json:"trans_at, optional"`
}

type ReceiptRecordQueryAllResponseX struct {
	List     []ReceiptRecordX `json:"list"`
	PageNum  int              `json:"pageNum" gorm:"-"`
	PageSize int              `json:"pageSize" gorm:"-"`
	RowCount int64            `json:"rowCount"`
}

type FrozenRecordQueryAllResponseX struct {
	List     []FrozenRecordX `json:"list"`
	PageNum  int             `json:"pageNum" gorm:"-"`
	PageSize int             `json:"pageSize" gorm:"-"`
	RowCount int64           `json:"rowCount"`
}

type DeductRecordQueryAllResponseX struct {
	List     []DeductRecordX `json:"list"`
	PageNum  int             `json:"pageNum" gorm:"-"`
	PageSize int             `json:"pageSize" gorm:"-"`
	RowCount int64           `json:"rowCount"`
}

type OrderActionQueryAllRequestX struct {
	OrderActionQueryAllRequest
	Orders []gormx.Sortx `json:"orders, optional" gorm:"-"`
}

type PersonalRepaymentRequestX struct {
	PersonalRepaymentRequest
	Orders []gormx.Sortx `json:"orders, optional" gorm:"-"`
}

type PersonalRepaymentResponseX struct {
	List     []PersonalRepayment `json:"list"`
	PageNum  int                 `json:"pageNum" gorm:"-"`
	PageSize int                 `json:"pageSize" gorm:"-"`
	RowCount int64               `json:"rowCount"`
}

type PersonalRepaymentX struct {
	PersonalRepayment
	TransAt JsonTime `json:"transAt, optional"`
}

type PersonalStatusUpdateResponseX struct {
	PersonalStatusUpdateResponse
	ChannelTransAt JsonTime `json:"channelTransAt, optional"`
}

type OrderFeeProfitQueryAllRequestX struct {
	OrderFeeProfitQueryAllRequest
	Orders []gormx.Sortx `json:"orders, optional" gorm:"-"`
}

type CalculateProfit struct {
	MerchantCode        string
	OrderNo             string
	Type                string
	CurrencyCode        string
	BalanceType         string
	ChannelCode         string
	ChannelPayTypesCode string
	OrderAmount         float64
	IsRate              string //代付是否用費率計算
}

type CalculateSubOrderProfit struct {
	OldOrderNo            string
	NewOrderNo            string
	OrderAmount           float64
	IsCalculateCommission bool
}

type PayOrderRequestX struct {
	PayOrderRequest
	MyIp string `json:"my_ip, optional"`
}

type PayQueryRequestX struct {
	PayQueryRequest
	MyIp string `json:"my_ip, optional"`
}

type PayQueryBalanceRequestX struct {
	PayQueryBalanceRequest
	MyIp string `json:"my_ip, optional"`
}

type ProxyPayRequestX struct {
	ProxyPayOrderRequest
	Ip string `json:"ip, optional"`
}

type WithdrawApiOrderRequestX struct {
	WithdrawApiOrderRequest
	MyIp string `json:"my_ip, optional"`
}

type MultipleOrderWithdrawCreateRequestX struct {
	List []OrderWithdrawCreateRequestX `json:"list"`
}

type OrderWithdrawCreateRequestX struct {
	OrderWithdrawCreateRequest
	MerchantCode    string `json:"merchantCode, optional"`
	MerchantOrderNo string `json:"merchant_order_no, optional"`
	UserAccount     string `json:"userAccount, optional"`
	NotifyUrl       string `json:"notify_url, optional"`
	PageUrl         string `json:"page_url, optional"`
	Source          string `json:"source, optional"`
	Type            string `json:"type, optional"`
}

type OrderWithdrawCreateResponse struct {
	OrderX
	Index []string `json:"index"`
	Errs  []string `json:"errs"`
}

type WithdrawApiQueryRequestX struct {
	WithdrawApiQueryRequest
	MyIp string `json:"my_ip, optional"`
}

type Rate struct {
	AgentLayerCode string `json:"agentLayerCode"`
	Rate float64 `json:"rate"`
}