package constants

const (
	TRANSACTION_TYPE_RECEIPT  = "1"  /*錢包紀錄_異動類型: 收款*/
	TRANSACTION_TYPE_UNFREEZE = "2"  /*錢包紀錄_異動類型: 解凍*/
	TRANSACTION_TYPE_REVERSE  = "3"  /*錢包紀錄_異動類型: 沖正*/
	TRANSACTION_TYPE_REFUND   = "4"  /*錢包紀錄_異動類型: 出款退回*/
	TRANSACTION_TYPE_MAKE_UP  = "5"  /*錢包紀錄_異動類型: 補單*/
	TRANSACTION_TYPE_PAYMENT  = "11" /*錢包紀錄_異動類型: 出款*/
	TRANSACTION_TYPE_FREEZE   = "12" /*錢包紀錄_異動類型: 凍結*/
	TRANSACTION_TYPE_RECOVER  = "13" /*錢包紀錄_異動類型: 追回*/
	TRANSACTION_TYPE_ADJUST   = "20" /*錢包紀錄_異動類型: 調整*/
)