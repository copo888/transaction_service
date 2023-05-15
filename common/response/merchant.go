package response

var (
	/**
	 * 通用系统讯息码
	 */
	//SUCCESS = "1050001" //"操作成功"

	/**
	 * 前端操作讯息码 10520XX
	 */

	FAIL_REASON_IS_NULL = "1052008" //"未通过原因，不得为空值"
	DELETE_FAIL         = "1052009" //"删除资料失败"
	UPDATE_FAIL         = "1052010" //"更新资料失败"

	/**
	 * 系统类型讯息码 10521XX
	 */
	SEND_EMAIL_FAILURE                                     = "1052112" //"发送邮件错误"
	CREATE_MERCHANT_ACCOUNT_FAILURE                        = "1052113" //"建立商户帐号失败"
	DELETE_MERCHANT_TEMP_FAILURE                           = "1052114" //"删除注册的暂存资料失败"
	INTERNAL_CHARGE_ORDER_REVIEWED_STATUS                  = "1052116" //"审核参数错误"
	MODIFY_MERCHANT_RATE_COUNT_NOT_EQUAL_BEFORE_RATE_COUNT = "1052117" //"修改商户费率数量和修改前费率数量不一致"
	UPPER_LAYER_STATUS_NOT_OPEN                            = "1052118" //"上层相关代理角色或状态尚未启用，请确认"
	CHANNEL_USE_TYPE_IS_POLLING                            = "1052119" //"商户渠道使用类型已经为：轮询模式"
	CHANNEL_USE_TYPE_IS_ASSIGN                             = "1052120" //"商户渠道使用类型已经为：非轮询模式=State{Code:指定)"
	PLEASE_DISABLE_LAYER_STATUS                            = "1052121" //"请先禁用代理身份，再进行启用状态操作"
	WITHDRAW_PASSWORD_NOT_SETTING                          = "1052122" // "下发密码未设定"
	/**
	 * 使用者操作错误 10522XX
	 */
	MERCHANT_INSUFFICIENT_XF_BALANCE = "1052200" //此商户下发余额不足
	USER_HAS_REGISTERED              = "1052201" //"该用户名已存在"
	MAILBOX_HAS_REGISTERED           = "1052202" //"该邮箱已存在"
	IP_DENIED                        = "1052205" //"此IP非法登錄，請設定白名單"

	MERCHANT_ACCOUNT_NOT_FOUND       = "1052206" //"无此商户帐号"
	MERCHANT_CARD_BINDED             = "1052207" //"此商户已經綁定此銀行卡"
	MERCHANT_INSUFFICIENT_DF_BALANCE = "1052208" //"此商户代付余额不足"
	MERCHANT_WALLET_RECORD_ERROR     = "1052209" //"钱包存取记录错误"
	MERCHANT_WALLET_ERROR            = "1052210" //"此商户无此钱包"
	MERCHANT_WALLET_ROLLBACK_ERROR   = "1052211" //"下发还款失败，下发扣款资料错误"
	MERCHANT_ORDER_NO_NOT_FOUND      = "1052213" //"商户无此单号"
	MERCHANT_CALLBACK_ERROR          = "1052214" //"商户回调失败"

	MERCHANT_AMOUNT_EXCEED      = "1052215" //"钱包金额超出限额"
	SETTING_MERCHANT_RATE_ERROR = "1052216" //"配置商户费率不可低于渠道成本费率"
	PAY_ORDER_ID_NOT_FOUND      = "1052217" //"无此支付订单号"
	MERCHANT_CHANNEL_NOT_SET    = "1052218" //"商户渠道未建立"

	MERCHANT_WALLET_TRANSFER_ERROR = "1052219" //"支转代钱包余额不足"
	MERCHANT_WALLET_UPDATE_ERROR   = "1052220" //"商户钱包更新错误"

	SETTING_MERCHANT_RATE_MIN_CHARGE_ERROR              = "1052221" //"配置商户费率低消不可低于渠道成本费率低消"
	SETTING_MERCHANT_CHARGE_ERROR                       = "1052222" //"配置商户手续费不可低于渠道成本手续费"
	SETTING_MERCHANT_RATE_LOWER_PARENT_ERROR            = "1052223" //"配置商户费率不可小于父层费率"
	SETTING_MERCHANT_RATE_MIN_CHARGE_LOWER_PARENT_ERROR = "1052224" //"配置商户费率低消不可小于父层费率低消"
	SETTING_MERCHANT_CHAREG_LOWER_PARENT_ERROR          = "1052225" //"配置商户手续费不可小于父层手续费"
	SETTING_MERCHANT_RATE_MIN_CHARGE_TO_NULL_ERROR      = "1052226" //"配置商户费率低消不可为空值"
	SETTING_MERCHANT_RATE_MIN_CHARGE_TO_NUM_ERROR       = "1052227" //"配置商户费率低消不可为数值"

	BANK_ACCOUNT_IN_BLACK_LIST = "1052228" //"银行账户或户名为黑名单账户"

	SETTING_NO_PARENT_MERCHANT_RATE_ERROR = "1052229" //"代理商户父层费率未设置，请先设置父层商户费率"

	SETTING_MERCHANT_PAY_CHAREG_LOWER_PARENT_ERROR             = "1052230" //"配置商户支付手续费不可小于父层手续费"
	SETTING_MERCHANT_PAY_CHARGE_ERROR                          = "1052231" //"配置商户支付手续费不可低于渠道支付手续费"
	SETTING_MERCHANT_RATE_OVER_LOWER_MERCHANT_ERROR            = "1052232" //"配置商户费率不得高于下层商户费率"
	SETTING_MERCHANT_CHARGE_OVER_LOWER_MERCHANT_ERROR          = "1052233" //"配置商户手续费不得高于下层商户手续费"
	SETTING_MERCHANT_RATE_MIN_CHARGE_OVER_LOWER_MERCHANT_ERROR = "1052234" //"配置商户费率低消不得高于下层商户费率低消"
	SETTING_MERCHANT_PAY_CHARGE_OVER_LOWER_MERCHANT_ERROR      = "1052235" //"配置商户支付手续费不得高于下层商户支付手续费"
	SETTING_MERCHANT_PAY_TYPE_SUB_CODING_ERROR                 = "1052236" //"支付方式指定代码格式错误，请输入数字1-6"
	SETTING_MERCHANT_PAY_TYPE_SUB_CODING_NULL_ERROR            = "1052237" //"支付方式指定代码不可为空值"

	SUB_WALLET_ENABLED_THEREFORE_OPERATION_PROHIBITED     = "1052238" //"已启用子钱包,所以禁止操作"
	SUB_WALLET_NOT_ENABLED_THEREFORE_OPERATION_PROHIBITED = "1052239" //"未启用子钱包,所以禁止操作"

	MERCHANT_INSUFFICIENT_PT_BALANCE = "1052240" // 子錢包餘額不足

	PAY_ORDER_CALLBACK                        = "1052250" //"此支付号已付款成功，不可再执行一次"
	PAY_CONFIRM_AMOUNT_REQUIRE                = "1052251" //"此支付号为异常单，必须输入金额以确认"
	SETTING_MERCHANT_PAY_CHARGE_TO_NULL_ERROR = "1052252" //"配置商户手续费不可为空值"
	PAY_ORDER_NOT_ERROR                       = "1052253" //"此支付号非异常，不可设定未付款"

	SMS_ORDER_NOT_FOUND = "1052254" //"简讯记录资料不存在"

	/**
	 * 錢包凍結及調整相關錯誤10523XX
	 */
	ORDER_NOT_SUCCESS_NO_FROZEN      = "1052300" //"此支付号已尚未付款，不可冻结"
	MERCHANT_FROZEN_VALUE_ERROR      = "1052301" //"輸入冻结金额，不得小于或等于零"
	MERCHANT_FROZEN_ERROR            = "1052302" //"商户冻结金额，不得小于零"
	MERCHANT_ADJUEST_ERROR           = "1052303" //"商户调整金额错误，余额不足"
	MERCHANT_FROZEN_ORDER_TYPE_ERROR = "1052304" //"商户冻结号型別错误"
	MERCHANT_FROZEN_ORDER_ALREADY    = "1052305" //"此订单已冻结"
	MERCHANT_ORDER_NO_FROZEN         = "1052306" //"此订单未冻结"
	MERCHANT_CURRENCY_NOT_SET        = "1052307" //"未指定币别"
	MERCHANT_CURRENCY_WRONG          = "1052308" //"币别不正确"
	MERCHANT_FROZEN_NOT_ENOUGH       = "1052309" //"账户余额不足以冻结"
	MERCHANT_UNFROZEN_ERROR          = "1052310" //"商户解冻金额不得大于冻结金额"
	MERCHANT_UNFROZEN_ORDER_ERROR    = "1052311" //"商户解冻後金额不得小于冻结訂單的金额"
	BALANCE_PROCESSING               = "1052312" //"商户钱包处理中"

	/**
	 * 商戶代理错误 10530XX
	 */
	MERCHANT_CHARGE_INFO_NOT_FOUND  = "1053000" //"代理商戶費率資訊未建立"
	MERCHANT_CHARGE_INFO_ERROR      = "1053001" //"代理商戶費率資訊層級错误"
	MERCHANT_CHARGE_INFO_ERROR_RATE = "1053002" //"代理商戶費率資訊错误，上層費率大於下層"
	MERCHANT_CHARGE_INFO_DUP_ERROR  = "1053003" //"代理商戶代理利潤資訊記錄已存在"
	MERCHANT_AGENT_NOT_FOUND        = "1053004" //"商戶代理層級編號未建立"

	/**
	 * 代理佣金錯誤10540XX
	 */
	MERCHANT_PROFIT_ORDER_DUPLICATE            = "1054001" //"利润提单重覆"
	MERCHANT_COMMISSION_RECORD_DUPLICATE       = "1054002" //"佣金纪录重覆"
	MERCHANT_COMMISSION_CALCULATE_BUY          = "1054003" //"佣金计算忙录，请重新请求"
	MERCHANT_COMMISSION_AUDIT                  = "1054004" //"佣金资料已审核，不能异动资料"
	MERCHANT_COMMISSION_WALLET_BUY             = "1054005" //"佣金钱包忙錄，请重新请求"
	MERCHANT_COMMISSION_TIME_ERROR             = "1054006" //"佣金结算时间错误"
	MERCHANT_COMMISSION_ORDER_ERROR            = "1054007" //"提单金额大於佣金钱包"
	MERCHANT_COMMISSION_ORDER_BUY              = "1054008" //"佣金提单忙錄，请重新请求"
	MERCHANT_COMMISSION_ORDER_REVIEWED         = "1054009" //"佣金提单已审核"
	MERCHANT_COMMISSION_ORDER_WITHDRAW_ERROR   = "1054010" //"非佣金提单日，不得提单"
	MERCHANT_COMMISSION_WALLET_NO_ENOUGH       = "1054011" //"佣金錢包餘額不足"
	SETTING_MERCAHNT_INCOME_OVER_CHN_MIN_LIMIT = "1054012" //"配置费率和手续费计算高于渠道单笔限额最小值"
	MERCHANT_COMMISSION_MER_AGENTLAYER_ERROR   = "1054013" //"层级编码搜寻条件错误"

	/**
	 * 通知訊息
	 */
	INVALID_NOTIFICATION_TYPE = "1055001" //"通知訊息未設定"

	/**
	 * 錯誤10560XX
	 */
	//INVALID_USDT_CHANNEL_CODING  = "1056001" //"虛擬幣渠道未設定"
	INVALID_USDT_RATE            = "1056002" //"虛擬幣匯率不正確"
	INVALID_USDT_EXCHANGE_LIMIT  = "1056003" //"虛擬幣換匯失敗，限額不正確"
	EXCHANGE_AMOUNT_NOT_ENOUGH   = "1056004" //"虛擬幣換匯失敗，金額不足"
	EXCHANGE_NOT_ENABLED         = "1056005" //"虛擬幣換匯功能未啟用"
	EXCHANGE_ORDER_NO_NOT_FOUND  = "1056006" //"无此換匯記錄单号"
	EXCHANGE_ORDER_STATUS_ERROR  = "1056007" //"換匯記錄狀態不合法"
	EXCHANGE_ORDER_CONTENT_ERROR = "1056008" //"換匯記錄返還內容不合法"
	EXCHANGE_ORDER_ADJUEST_ERROR = "1056009" //"換匯記錄返還內容與調整金額不符合"
)
