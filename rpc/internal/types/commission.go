package types

type CommissionMonthReportX struct {
	CommissionMonthReport
	ConfirmAt JsonTime `json:"confirmAt, optional"`
	CreatedAt JsonTime `json:"createdAt, optional"`
	UpdatedAt JsonTime `json:"createdAt, optional"`
}

type CommissionMonthReport struct {
	ID                        int64   `json:"id"`
	MerchantCode              string  `json:"merchantCode"`
	AgentLayerNo              string  `json:"agentLayerNo"`
	Month                     string  `json:"month"`
	CurrencyCode              string  `json:"currencyCode"`
	Status                    string  `json:"status"`
	PayTotalAmount            float64 `json:"payTotalAmount"`
	PayCommission             float64 `json:"payCommission"`
	InternalChargeTotalAmount float64 `json:"internalChargeTotalAmount"`
	InternalChargeCommission  float64 `json:"internalChargeCommission"`
	ProxyPayTotalNumber       float64 `json:"proxyPayTotalNumber"`
	ProxyPayCommission        float64 `json:"proxyPayCommission"`
	TotalCommission           float64 `json:"totalCommission"`
	ChangeCommission          float64 `json:"changeCommission"`
	Comment                   string  `json:"comment"`
	ConfirmBy                 string  `json:"confirmBy"`
}

type CommissionMonthReportDetailX struct {
	CommissionMonthReportDetail
	CreatedAt JsonTime `json:"createdAt, optional"`
	UpdatedAt JsonTime `json:"createdAt, optional"`
}

type CommissionMonthReportDetail struct {
	CommissionMonthReportId int64   `json:"commission_month_report_id"`
	MerchantCode            string  `json:"merchantCode"`
	PayTypeCode             string  `json:"payTypeCode"`
	OrderType               string  `json:"orderType"`
	MerchantFee             float64 `json:"merchantFee"`
	AgentFee                float64 `json:"agentFee"`
	DiffFee                 float64 `json:"diffFee"`
	MerchantHandlingFee     float64 `json:"merchantHandlingFee"`
	AgentHandlingFee        float64 `json:"agentHandlingFee"`
	DiffHandlingFee         float64 `json:"diffHandlingFee"`
	TotalAmount             float64 `json:"totalAmount"`
	TotalNumber             float64 `json:"totalNumber"`
	TotalCommission         float64 `json:"totalCommission"`
}
