package constants

const (
	UI  = "1" /*平台提单*/
	API = "2" /*API提单*/

	WAIT_PROCESS = "0"  /*订单状态:待处理*/
	PROCESSING   = "1"  /*订单状态:处理中*/
	TRANSACTION  = "2"  /*订单状态:交易中*/
	SUCCESS      = "20" /*订单状态:成功*/
	FAIL         = "30" /*订单状态:失败*/
	FROZEN       = "31" /*订单状态:凍結*/

	CALL_BACK_STATUS_PROCESSING = "0" /*渠道回調狀態: 處理中*/
	CALL_BACK_STATUS_SUCCESS    = "1" /*渠道回調狀態: 成功*/
	CALL_BACK_STATUS_FAIL       = "2" /*渠道回調狀態: 失敗*/

	ORDER_TYPE_NC = "NC" /*订单类型:内充*/
	ORDER_TYPE_DF = "DF" /*订单类型:代付*/
	ORDER_TYPE_ZF = "ZF" /*订单类型:支付*/
	ORDER_TYPE_XF = "XF" /*订单类型:下发*/

	IS_LOCK_NO  = "0" /*是否锁定状态: 否*/
	IS_LOCK_YES = "1" /*是否锁定状态: 是*/

	IS_TEST_NO  = "0" /*是否锁定状态: 否*/
	IS_TEST_YES = "1" /*是否锁定状态: 是*/

	IS_MERCHANT_CALLBACK_YES      = "1" /*是否已經回調商戶: 是*/
	IS_MERCHANT_CALLBACK_NO       = "0" /*是否已經回調商戶: 否*/
	IS_MERCHANT_CALLBACK_NOT_NEED = "2" /*是否已經回調商戶: 不需*/

	PERSON_PROCESS_STATUS_WAIT_PROCESS = "0"  /*人工处理状态: 待處理*/
	PERSON_PROCESS_STATUS_PROCESSING   = "1"  /*人工处理状态: 處理中*/
	PERSON_PROCESS_STATUS_SUCCESS      = "2"  /*人工处理状态: 成功*/
	PERSON_PROCESS_STATUS_FAIL         = "3"  /*人工处理状态: 失敗*/
	PERSON_PROCESS_STATUS_NO_ROCESSING = "10" /*人工处理状态: 不需处理*/

	DF_BALANCE = "DFB"
	XF_BALANCE = "XFB"
	YJ_BALANCE = "YJB"

	IS_CALCULATE_PROFIT_YES = "1" /*已記算傭金利潤: 是*/
	IS_CALCULATE_PROFIT_NO  = "0" /*已記算傭金利潤: 否*/

	ACTION_FROZEN             = "FROZEN"        //冻结
	ACTION_MAKE_UP_ORDER      = "MAKE_UP_ORDER" //补单
	ACTION_PLACE_ORDER        = "PLACE_ORDER"   //创建订单
	ACTION_REVIEW_FAIL        = "REVIEW_FAIL"   //审核失败
	ACTION_MAKE_UP_LOCK_ORDER = "MAKE_UP_LOCK_ORDER"
	ACTION_PERSON_PROCESSING  = "PERSON_PROCESSING" //人工处理中
	ACTION_PROCESS_SUCCESS    = "PROCESS_SUCCESS"   //人工处理通过
	ACTION_UNFROZEN           = "UNFROZEN"          //解冻
	ACTION_REVERSAL           = "REVERSAL"          //冲正
	ACTION_REVIEW_SUCCESS     = "REVIEW_SUCCESS"    //审核成功
	ACTION_SUCCESS            = "SUCCESS"           //成功
	ACTION_FAILURE            = "FAILURE"           //失敗
	ACTION_DF_REFUND          = "DF_REFUND"         //代付还款
	ACTION_TRANSFER_TEST      = "TRANSFER_TEST"     //轉測試單
	ACTION_TRANSFER_NORMAL    = "TRANSFER_NORMAL"   //轉正式單

	ORDER_REASON_TYPE_UPDATE_AMOUNT = "1"  // 修改金額
	ORDER_REASON_TYPE_REPAYMENT     = "2"  // 重複支付
	ORDER_REASON_TYPE_OTHER         = "3"  // 其它
	ORDER_REASON_TYPE_RECOVER       = "11" // 追回
)
