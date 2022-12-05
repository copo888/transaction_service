package types

type IncomReport struct {
	ID int64 `json:"id"`
	CurrencyCode string `json:"currency_code"`
	Month string `json:"month"`
	PayTotalAmount float64 `json:"pay_total_amount"`
	PayNetProfit float64 `json:"pay_net_profit"`
	InternalChargeTotalAmount float64 `json:"internal_charge_total_amount"`
	InternalChargeNetProfit float64 `json:"internal_charge_net_profit"`
	WithdrawTotalAmount float64 `json:"withdraw_total_amount"`
	WithdrawNetProfit float64 `json:"withdraw_net_profit"`
	ProxyPayTotalAmount float64 `json:"proxy_pay_total_amount"`
	ProxyPayNetProfit float64 `json:"proxy_pay_net_profit"`
	ReceivedTotalNetProfit float64 `json:"received_total_net_profit"`
	RemitTotalNetProfit float64 `json:"remit_total_net_profit"`
	TotalNetProfit float64 `json:"total_net_profit"`
	CommissionTotalAmount float64 `json:"commission_total_amount"`
	ProfitGrowthRate float64 `json:"profit_growth_rate"`
	TotalAllocHandlingFee float64 `json:"total_alloc_handling_fee"`
}

type IncomReportX struct {
	IncomReport
	CreatedAt JsonTime
	UpdatedAt JsonTime
}

type CaculateMonthProfitReport struct {
	MerchantCode string `json:"merchant_code"`
	CurrencyCode string `json:"currency_code"`
	TotalProfit float64 `json:"total_profit"`
	PayTypeCode string `json:"pay_type_code"`
	TotalAmount float64 `json:"total_amount"`
}