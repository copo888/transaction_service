package response

var (
	/**
	 * 系统类型讯息码
	 */
	ILLEGAL_REQUEST         = "1062101" //"非法请求"
	ILLEGAL_PARAMETER       = "1062102" //"非法参数"
	UN_ROLE_AUTHORIZE       = "1062103" //"无此应用服务使用授权"
	DATA_NOT_FOUND          = "1062104" //"资料无法取得"
	GENERAL_ERROR           = "1062105" //"通用错误"
	CONNECT_SERVICE_FAILURE = "1062106" //"服务连线失败"
	UPDATE_DATABASE_FAILURE = "1062107" //"数据更新失败"
	UPDATE_DATABASE_REPEAT  = "1062108" //"数据重复更新"
	INSERT_DATABASE_FAILURE = "1062109" //"数据新增失败"
	DELETE_DATABASE_FAILURE = "1062110" //"数据删除失败"
	DATABASE_FAILURE        = "1062111" //"数据库错误"

	CHANNEL_CLOSED_OR_DEFEND                      = "1062112" //"渠道关闭或维护中"
	RATE_NOT_CONFIGURED                           = "1062113" //"未配置商户渠道费率"
	REPLY_MESSAGE_MALFORMED                       = "1062114" //"返回资讯格式错误"
	PROXY_BAL_MIN_LIMIT_NOT_CONFIGURED            = "1062115" //"渠道代付馀额下限值未设定"
	CHN_BALANCE_NOT_ENOUGH                        = "1062116" //"渠道余额不足扣款金额"
	SINGLE_LIMIT_SETTING_ERROR                    = "1062117" //"单笔限额设定错误，最小值必需小于最大值"
	INVALID_USDT_CHANNEL_CODING                   = "1062118" //"无效的USDT渠道编码"
	USDT_CHANNEL_NAME_DIFFERENT                   = "1062119" //"与对应的USDT渠道名称不一致"
	USDT_CHANNEL_REPEAT_DESIGNATION               = "1062120" //"对应的USDT渠道已被指定"
	RESET_DESIGNATION_MER_RATE_ERROR              = "1062121" //"不得重置已指定的商户费率"
	RESETRADIS_FROM_DB_FAILURE                    = "1062122" //"更新Redis資料錯誤"
	INVALID_CHANNEL_INFO                          = "1062123" //"無對應相關渠道資料"
	RATE_SETTING_ERROR                            = "1062124" //"费率设定异常"
	RATE_NOT_CONFIGURED_OR_CHANNEL_NOT_CONFIGURED = "1062125" //"未配置商户渠道费率或渠道配置错误"
	PROXY_TRANSACTION_FAILURE                     = "1062126" //"代付交易錯誤"
)
