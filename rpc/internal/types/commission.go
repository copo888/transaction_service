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
	PayCommissionTotalAmount  float64 `json:"payCommissionTotalAmount"`
	InternalChargeTotalAmount float64 `json:"internalChargeTotalAmount"`
	InternalChargeCommission  float64 `json:"internalChargeCommission"`
	ProxyPayTotalAmount       float64 `json:"proxyPayTotalAmount"`
	ProxyPayTotalNumber       float64 `json:"proxyPayTotalNumber"`
	ProxyPayCommission        float64 `json:"proxyPayCommission"`
	ProxyCommissionTotalAmount float64 `json:"proxyCommissionTotalAmount"`
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
	CommissionTotalAmount   float64 `json:"commissionTotalAmount"`
	TotalNumber             float64 `json:"totalNumber"`
	TotalCommission         float64 `json:"totalCommission"`
}

type UpdateCommissionAmount struct {
	MerchantCode            string
	CurrencyCode            string
	CommissionMonthReportId int64
	OrderNo                 string
	TransactionType         string
	TransferAmount          float64
	Comment                 string
	CreatedBy               string
}

type CommissionWithdrawOrderX struct {
	CommissionWithdrawOrder
	CreatedAt JsonTime `json:"createdAt, optional"`
	UpdatedAt JsonTime `json:"createdAt, optional"`
}

type CommissionWithdrawOrder struct {
	ID                   int64   `json:"id, optional"`
	OrderNo              string  `json:"orderNo, optional"`
	MerchantCode         string  `json:"merchantCode, optional"`
	WithdrawCurrencyCode string  `json:"withdrawCurrencyCode, optional"`
	PayCurrencyCode      string  `json:"payCurrencyCode, optional"`
	WithdrawAmount       float64 `json:"withdrawAmount, optional"`
	ExchangeRate         float64 `json:"exchangeRate, optional"`
	PayAmount            float64 `json:"payAmount, optional"`
	AfterCommission      float64 `json:"afterCommission, optional"`
	Remark               string  `json:"remark, optional"`
	AttachmentPath       string  `json:"attachmentPath, optional"`
	CreatedBy            string  `json:"createdBy, optional"`
	PayAt                string  `json:"payAt, optional"`
}