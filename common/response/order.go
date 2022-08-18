package response

var (
	WITHDRAW_AMT_NOT_REACH_MIN_LIMIT         = "1091001" //下发金额未达下限
	WITHDRAW_AMT_EXCEED_MAX_LIMIT            = "1091002" //下发金额超过上限
	MER_WITHDRAW_CHARGE_NOT_SET              = "1091003" //未设定商户下发手续费
	MER_WITHDRAW_MIN_LIMIT_NOT_SET           = "1091004" //未设定商户下发金额下限
	MER_WITHDRAW_MAX_LIMIT_NOT_SET           = "1091005" //未设定商户下发金额上限
	AVAILABLE_AMT_NOT_ENOUGH                 = "1091006" //可下发金额不足
	ORDER_STATUS_WRONG_CANNOT_REVERSAL       = "1091034" //订单状态错误，不可冲正
	ORDER_STATUS_WRONG_CANNOT_REPAYMENT      = "1091035" //订单状态错误，不可人工还款
	ORDER_STATUS_WRONG_CANNOT_PROCESSING     = "1091036" //訂單狀態錯誤，不可改為處理中
	INVALID_WITHDRAW_ORDER_NO                = "1092004" //下发订单号错误
	COMPLETED_ORDER_REVIVEW_REPEAT           = "1092007" //订单已结单,不可重复审核
	MERCHANT_WITHDRAW_AUDIT_ERROR            = "1092008" //下发金额人员输入不符
	MERCHANT_WITHDRAW_RECORD_ERROR           = "1092009" //下发渠道明细错误
	MERCHANT_REVERSAL_AUDIT_ERROR            = "1093001" //沖正金額人員輸入不符
	CURRENCY_NOT_THE_SAME                    = "1093002" //订单币别不相同
	ORDER_STATUS_WRONG_CANNOT_UNFROZEN       = "1094000" //訂單狀態錯誤,不可解凍訂單
	ORDER_STATUS_WRONG_CANNOT_FROZEN         = "1094001" //訂單狀態錯誤,不可凍結訂單
	FROZEN_AMOUNT_NOT_LESS_THAN_ORDER_AMOUNT = "1094002" //冻结金额不可小於交易金额
	ORDER_STATUS_WRONG_CANNOT_MAKE_UP        = "1094006" //订单状态错误,不可补单
	AMOUNT_MUST_BE_GREATER_THAN_ZERO         = "1094007" //金额需大於零
	ORDER_IS_MAKE_UP_ORDER_DONT_MAKE_UP      = "1094008" //此订单为补单, 不可再补单
	ONLY_SUCCESSFUL_ORDER_CAN_CALL_BACK      = "1094009" //只有成功單可以回調
	ORDER_DOES_NOT_NEED_CALL_BACK            = "1094010" //此订单不需回調
	ORDER_IS_STATUS_IS_LOCK                  = "1094011" //此订单狀態已鎖定
	ORDER_STATUS_WRONG                       = "1094012" //订单状态錯誤
	ORIGINAL_ORDER_NOT_CALCULATED_COMMISSION = "1094013" //原訂單還未計算傭金,請稍後
	ORDER_TYPE_IS_WRONG                      = "1094014" //订单类型错误
	ORDER_ALREADY_IS_TEST                    = "1094015" //已為測試單
	ORDER_NOT_CALCULATED_PROFIT_PLEASE_WAIT  = "1094016" //有订单还未计算佣金利润,请稍后
	ORDER_STATUS_IS_NOT_FOR_TEST             = "1094017" //只有成功單能改測試單
	BALANCE_REDISLOCK_ERROR                  = "1094018" //钱包LOCK错误
)
