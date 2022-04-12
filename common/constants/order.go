package constants

const (
	UI  = "1" /*平台提单*/
	API = "2" /*API提单*/

	WAIT_PROCESS = "0"  /*订单状态:待处理*/
	PROCESSING   = "1"  /*订单状态:处理中*/
	SUCCESS      = "20" /*订单状态:成功*/
	FAIL         = "30" /*订单状态:失败*/

	CALL_BACK_STATUS_PROCESSING = "0" /*渠道回調狀態: 處理中*/
	CALL_BACK_STATUS_SUCCESS    = "1" /*渠道回調狀態: 成功*/
	CALL_BACK_STATUS_FAIL       = "2" /*渠道回調狀態: 失敗*/

	ORDER_TYPE_NC = "NC" /*订单类型:内充*/
	ORDER_TYPE_DF = "DF" /*订单类型:代付*/
	ORDER_TYPE_ZF = "ZF" /*订单类型:支付*/
	ORDER_TYPE_XF = "XF" /*订单类型:下发*/

	IS_LOCK_NO  = "0" /*是否锁定状态: 否*/
	IS_LOCK_YES = "1" /*是否锁定状态: 是*/

	IS_MERCHANT_CALLBACK_YES = "1" /*是否已經回調商戶: 是*/
	IS_MERCHANT_CALLBACK_NO  = "0" /*是否已經回調商戶: 否*/

	PERSON_PROCESS_STATUS_WAIT_PROCESS = "0"  /*人工处理状态: 待處理*/
	PERSON_PROCESS_STATUS_PROCESSING   = "1"  /*人工处理状态: 處理中*/
	PERSON_PROCESS_STATUS_SUCCESS      = "2"  /*人工处理状态: 成功*/
	PERSON_PROCESS_STATUS_FAIL         = "3"  /*人工处理状态: 失敗*/
	PERSON_PROCESS_STATUS_NO_ROCESSING = "10" /*人工处理状态: 不需处理*/

//(1=收款; 2=解凍; 3=沖正; 11=出款 ; 12=凍結)

)
