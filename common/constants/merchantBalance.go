package constants

const (
	TRANSACTION_TYPE_INTERNAL_CHARGE = "1"  /*錢包紀錄_異動類型: 內充*/
	TRANSACTION_TYPE_UNFREEZE        = "2"  /*錢包紀錄_異動類型: 解凍*/
	TRANSACTION_TYPE_REVERSE         = "3"  /*錢包紀錄_異動類型: 沖正*/
	TRANSACTION_TYPE_REFUND          = "4"  /*錢包紀錄_異動類型: 還款*/
	TRANSACTION_TYPE_MAKE_UP         = "5"  /*錢包紀錄_異動類型: 補單*/
	TRANSACTION_TYPE_PAY             = "6"  /*錢包紀錄_異動類型: 支付*/
	TRANSACTION_TYPE_PROXY_PAY       = "11" /*錢包紀錄_異動類型: 代付*/
	TRANSACTION_TYPE_FREEZE          = "12" /*錢包紀錄_異動類型: 凍結*/
	TRANSACTION_TYPE_RECOVER         = "13" /*錢包紀錄_異動類型: 追回*/
	TRANSACTION_TYPE_DEDUCT          = "14" /*錢包紀錄_異動類型: 扣回*/
	TRANSACTION_TYPE_ISSUED          = "15" /*錢包紀錄_異動類型: 下發*/
	TRANSACTION_TYPE_ADJUST          = "20" /*錢包紀錄_異動類型: 調整*/

	COMMISSION_TRANSACTION_TYPE_MONTHLY    = "1"   /*佣金紀錄_異動類型: 月結佣金*/
	COMMISSION_TRANSACTION_TYPE_WITHDRAWAL = "11"  /*佣金紀錄_異動類型: 佣金提領*/
)
