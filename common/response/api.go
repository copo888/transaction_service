package response

var (
	API_SUCCESS                 = "000" // "操作成功"
	NO_RECORD_DATA              = "001" // "无记录"
	API_GENERAL_ERROR           = "002" // "系统忙碌中,请稍后再试"
	GENERAL_EXCEPTION           = "003" // "系统錯誤,请稍后再试"
	API_INVALID_PARAMETER       = "004" // "无效参数"
	SERVICE_RESPONSE_ERROR      = "005" // "服务回傳失败"
	SERVICE_RESPONSE_DATA_ERROR = "006" // "服务回傳資料錯誤"
	API_IP_DENIED               = "007" // "此IP非法登錄，請設定白名單"
	ContentType_ERROR           = "008" // "内容类型错误，请使用 application/json"
	API_PARAMETER_TYPE_ERROE    = "009" // "JSON格式或参数类型错误"
	WAIT_LOCK_EXCEPTION         = "010" // "此交易目前正执行中，请稍后再试"
	CACHED_DATA_NOT_FOUND       = "011" // "此交易已执行完毕或交易参数错误"
	/**
	 * 参数错误讯息码
	 */
	INVALID_TIMESTAMP                  = "101" // "无效时间戳"
	INVALID_SIGN                       = "102" // "无效验签"
	INVALID_CURRENCY_CODE              = "103" // "无效货币编码"
	INVALID_ORDER_NO                   = "104" // "无效订单编号"
	REPEAT_ORDER_NO                    = "105" // "重复订单号"
	INVALID_START_DATE                 = "106" // "无效开始日期时间"
	INVALID_END_DATE                   = "107" // "无效结束日期时间"
	INVALID_DATE_RANGE                 = "108" // "无效日期范围"
	INVALID_DATE_TYPE                  = "109" // "无效日期筛选类型"
	INVALID_MERCHANT_CODE              = "110" // "无效商户号"
	MERCHANT_IS_DISABLE                = "111" // "商户已被禁用或结清"
	INVALID_AMOUNT                     = "112" // "无效金额"
	INVALID_LANGUAGE_CODE              = "113" // "无效语言编码"
	INVALID_BANK_ID                    = "114" // "无效开户行号"
	INVALID_BANK_NAME                  = "115" // "无效开户行名"
	INVALID_BANK_NO                    = "116" // "无效银行卡号"
	INVALID_DEFRAY_NAME                = "117" // "无效开户人姓名"
	INVALID_ACCESS_TYPE                = "118" // "无效接入类型"
	INVALID_NOTIFY_URL                 = "119" // "无效URL格式"
	SIGN_ERROR                         = "120" // "签名出错"
	NO_AVAILABLE_CHANNEL_ERROR         = "121" // "没有可用通道"
	CHANNEL_NOT_OPEN_OR_NOT_MEET_RULES = "122" // "指定通道没有开启或不符合指定通道规则"
	NO_CHANNEL_SET                     = "123" // "未指定通道或不符合指定通道规则"
	INVALID_MERCHANT_ID                = "124" // "无效商户号"
	INVALID_MERCHANT_AGENT             = "125" // "无效代理商户"
	INVALID_MERCHANT_ACCOUNT           = "126" // "无效商户帐号"
	ERROR_CHANGE_PASSWORD              = "127" // "商户密码变更错误"
	ERROR_CHANGE_MERCHANT_KEY          = "128" // "商户密钥变更错误"
	INVALID_USER_NAME                  = "129" // "开户人姓名无效或未输入"
	INVALID_USDT_TYPE                  = "130" // "无效协议"
	INVALID_USDT_WALLET_ADDRESS        = "131" // "无效钱包地址"
	INVALID_PAY_TYPE_SUB_NO            = "132" // "多指定模式，PayTypeSubNo為必填"

	// for channel test only
	INVALID_MERCHANT_OR_CHANNEL_PAYTYPE = "160" // "資料庫無此商户号或商户未配置渠道、支付方式等設定错误或关闭或维护"
	INVALID_CHANNEL_PAYTYPE_CURRENCY    = "161" // "商户配置之渠道支付方式與幣別有誤"
	/**
	 *  交易错误讯息码
	 */
	TRANSACTION_FAILURE             = "201" // "交易失败"
	INSUFFICIENT_IN_AMOUNT          = "202" // "余额不足"
	CURRENCY_INCONSISTENT           = "203" // "商户货币別不一致"
	IS_LESS_MINIMUN_AMOUNT          = "204" // "单笔小于最低交易金额"
	IS_GREATER_MXNIMUN_AMOUNT       = "205" // "单笔大于最高交易金额"
	MERCHANT_IS_NOT_SETTING_CHANNEL = "206" // "尚未配置渠道"
	BANK_CODE_EMPTY                 = "207" // "银行代码不得为空值"
	BANK_CODE_INVALID               = "208" // "银行代码错误"
	PAY_TYPE_INVALID                = "209" // "支付类型代码错误"
	CHANNEL_REPLY_ERROR             = "210" // "渠道返回错误"
	INVALID_STATUS_CODE             = "211" // "Http状态码错误"
	/**
	 * 内部错误
	 */
	INTERNAL_SIGN_ERROR = "301" // "系统验签错误"

	/**
	 * 系统层级错误
	 */
	SYSTEM_ERROR  = "400" // "系统错误"
	NETWORK_ERROR = "401" // "网路异常"

	/**
	 * 應用层级错误： 支付 500~599
	 */
	ORDER_NUMBER_EXIST                            = "500" // "商户订单号重复"
	ORDER_NUMBER_NOT_EXIST                        = "501" // "商户订单号不存在"
	MERCHANT_PAY_TYPE_INVALID_OR_CHANNEL_NOT_SET  = "502" // "商户代码[%s]或支付类型代码[%s]或幣別[%s]错误或指定渠道设定错误或关闭或维护"
	ORDER_AMOUNT_INVALID                          = "503" // "商户下单金额错误"
	ORDER_AMOUNT_LIMIT_MIN                        = "504" // "商户下单金额太小"
	ORDER_AMOUNT_LIMIT_MAX                        = "505" // "商户下单金额太大"
	WALLET_NOT_SET                                = "506" // "商户渠道錢包未设定"
	API_MERCHANT_CHANNEL_NOT_SET                  = "507" // "商户渠道未建立"
	MERCHANT_PAY_TYPE_INVALID_OR_CHANNEL_NOT_SET2 = "508" // "商户代码[%s]或支付类型代码[%s][%s]错误或指定渠道设定错误或关闭或维护"
	WALLET_UPDATE_ERROR                           = "509" // "商户錢包資料错误"
	ORDER_AMOUNT_ERROR                            = "510" // "商户下单金额和回調金額不符"
	ORDER_BANK_NO_LEN_ERROR                       = "511" // "银联行账(卡)号，长度必须13~22位"
)
