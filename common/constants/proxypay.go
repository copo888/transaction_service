package constants

const (
	/*訂單狀態(0:待處理 1:處理中 20:成功 30:失敗 31:凍結)*/
	CHN_PAY_TYPE_PROXY_PAY = "DF"
	/** 0.待处理 **/
	PROXY_PAY_WAIT = "0"
	/** 1.处理中(排程捞出来后，未进到异步处理前) **/
	PROXY_PAY_PROCESSING = "1"
	/** 2.成功 **/
	PROXY_PAY_SUCCESS = "20"
	/** 3.失败 **/
	PROXY_PAY_FAIL = "30"
	/** 4.交易中 **/
	PROXY_PAY_FROZEN = "40"

	//**************** 代付订单错误类型 *******************
	/** 0.其他 **/
	ERROR0_OTHER_FAIL = "0"
	/** 1.未配置渠道 **/
	ERROR1_NO_CHANNEL = "1"
	/** 2.渠道余额不足 **/
	ERROR2_CHANNEL_INSUFFICIENT_BALANCE = "2"
	/** 3.代付帐户错误 **/
	ERROR3_PROXY_PAY_ACCOUNT_FAIL = "3"
	/** 4.商户余额不足 **/
	ERROR4_MERCHANT_INSUFFICIENT_BALANCE = "4"
	/** 5.扣款失败 **/
	ERROR5_DEDUCTION_FAIL = "5"
	/** 6.交易账户为黑名单 **/
	ERROR6_BANK_ACCOUNT_IS_BLACK = "6"

	//**************** 代付渠道送出状态 *******************
	/** 0.未上传 **/
	CHANNEL_NOT_UPLOAD = "0"
	/** 1.已上传 **/
	CHANNEL_UPLOAD_OK = "1"
	/** 2.上传失败 **/
	CHANNEL_UPLOAD_FAIL = "2"

	//**************** 代付手续费计算方式 *******************
	/** 1.费率 **/
	SUMARY_TYPE_RATE = "1"
	/** 2.单笔低消 **/
	SUMARY_TYPE_MINIMUM = "2"

	//**************** 订单来源类型 *******************
	/** 1.平台订单 **/
	ORDER_SOURCE_BY_PLATFORM = "1"
	/** 2.API订单 **/
	ORDER_SOURCE_BY_API = "2"

	//**************** 搜寻区间的类型 *******************
	/** 1.提单时间查询 **/
	SEARCH_BY_ORDER_DATE = "1"
	/** 2.审核时间查询 **/
	SEARCH_BY_REVIEW_DATE = "2"

	//**************** 渠道回复处里状态 *******************
	/** 0.待处理 **/
	CHANNEL_PROCESS_WAIT = "0"
	/** 1.处理中 **/
	CHANNEL_PROCESS_ING = "1"
	/** 2.成功 **/
	CHANNEL_PROCESS_SUCCESS = "2"
	/** 3.失败 **/
	CHANNEL_PROCESS_FAIL = "3"

	//**************** 还款处理状态 ***********************
	/** 0.不需还款 **/
	REPAYMENT_NOT = "0"
	/** 1.待还款 **/
	REPAYMENT_WAIT = "1"
	/** 2.还款成功 **/
	REPAYMENT_SUCCESS = "2"
	/** 3.还款失败 **/
	REPAYMENT_FAIL = "3"

	//**************** 代付提单错误显示 ***********************
	/** 未配置渠道 **/
	PROXYPAY_ERROR_NO_CHANNEL = "商户尚未配置渠道"
	/** 余额不足 **/
	PROXYPAY_ERROR_NO_MONEY = "余额不足"
	/** 扣款失败 **/
	PROXYPAY_ERROR_DEDUCTION_FAIL = "扣款失败"
	/** 渠道返还处理失败 **/
	PROXYPAY_ERROR_CHANNEL_PROCESS_FAIL = "渠道返还处理失败"
	/** 交易账户为黑名单 **/
	BANK_ACCOUNT_IS_BLACK = "交易银联卡(号)为黑名单"

	//**************** 结算类型 ***********************
	/** 扣款记录 **/
	PROXYPAY_SUM_TYPE_DEDUCTION = "1"
	/** 还款记录 **/
	PROXYPAY_SUM_TYPE_REPAYMENT = "2"

	//**************** 查询状态 ***********************
	/** 不需查询 **/
	ORDER_DONT_CHECK = "0"
	/** 等待查询 **/
	ORDER_NEED_CHECK = "1"
	/** 已查询 **/
	ORDER_IS_CHECK = "2"

	//**************** 人工处里状态 ***********************
	/** 不需人工 **/
	PERSON_PROCESS_NONE = "0"
	/** 等待人工处里(无查询API) **/
	PERSON_PROCESS_WAIT_NO_API = "1"
	/** 等待人工处里(查询异常) **/
	PERSON_PROCESS_WAIT_CHECK_FAIL = "2"
	/** 等待人工处里(人工转单) **/
	PERSON_PROCESS_WAIT_CHANGE = "3"
	/** 人工还款失败 **/
	PERSON_PROCESS_REPAYMENT_FAIL = "4"
	/** 已处里 **/
	PERSON_PROCESS_END = "9"

	//**************** 是否已回调商户 ***********************
	/** 未回调 **/
	MERCHANT_CALL_BACK_NO = "0"
	/** 已回调 **/
	MERCHANT_CALL_BACK_YES = "1"
	/** 不需回调 **/
	MERCHANT_CALL_BACK_DONT_USE = "2"

	//**************** 接入类型 ***********************
	/** 1.代付 **/
	ACCESS_TYPE_PROXY = "1"

	//**************** 语系 ***********************
	/** 简体中文 **/
	LANGUAGE_ZH_CN = "zh-CN"
	/** 繁体中文 **/
	LANGUAGE_ZH_TW = "zh-TW"

	//**************** 货币编码 ***********************
	/** 人民币 **/
	CURRENCY_CNY = "CNY"
	/** 台币 **/
	CURRENCY_TWD = "TWD"

	//**************** 代付订单结算 ***********************
	/** 扣款记录 **/
	PROXY_PAY_SUM_TYPE_WITHHOLD = "1"
	/** 还款记录 **/
	PROXY_PAY_SUM_TYPE_REPAYMENT = "2"

	//**************** 人工处里结单类型 *******************
	/** 成功单 **/
	PERSON_PROCESS_SUCCESS_ORDER = "1"
	/** 失败单 **/
	PERSON_PROCESS_FAIL_ORDER = "2"

	//**************** 提单计算佣金类型 *******************
	/** 支付单 **/
	COMMISSION_ZF_ORDER = "1"
	/** 内充单 **/
	COMMISSION_NC_ORDER = "2"
	/** 代付单 **/
	COMMISSION_DF_ORDER = "3"

	//
	TRANS_TYPE_PROXY_PAY               = "proxy-pay"
	TRANS_TYPE_QUERY_ORDER             = "query-order"
	TRANS_TYPE_QUERY_BALANCE           = "query-balance"
	TRANS_TYPE_PROXY_PAY_QUERY         = "proxy-pay-query"
	TRANS_TYPE_PROXY_PAY_QUERY_BALANCE = "proxy-pay-query-balance"
)
