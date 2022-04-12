package constants

var RegexpFeAccountName = "^[a-zA-Z0-9]{6,12}$"

var RegexpFePassword = "^[a-zA-Z0-9]{8,16}$"

var RegexpEmail = "^[_a-zA-Z0-9-]+([.][_a-zA-Z0-9-]+)*@[a-zA-Z0-9-]+([.][a-zA-Z0-9-]+)*$"

var RegexpVerificationCode = "^[A-Za-z0-9]{6}"

var RegexpMobilePhone = "1\\d{10}"

var RegexpChinese = "[\\u4E00-\\u9FA5]+"

var RegexpDate = "/^\\d{4}-\\d{2}-\\d{2}$/"

var RegexpTime = "\\d{14}"

var RegexpIpaddressPattern = "^([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\." +
	"([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\." +
	"([01]?\\d\\d?|2[0-4]\\d|25[0-5])\\." +
	"([01]?\\d\\d?|2[0-4]\\d|25[0-5])$"

var RegexpWebsiteUrl = "^https?://(.*):?[0-9]{0,4}/?\\w+"

/** ====================================渠道相關==============================================**/

var RegexpPayType = "^[A-Z0-9]{2}"

var RegPayTypeNum = "^[A-Z0-9]{2}[0-9]{1}"

var RegexpChnMerchantCoding = "[^\\u4E00-\\u9FA5]+"

var RegexpUrl = "^https?://(.*):[0-9]{1,4}/?\\w+"

var RegexpRateNumberic = "^[0-9]+(.[0-9]{1,})?$"

var RegexpMerchantId = "^[a-zA-Z0-9]*$"

var RegexpChannelUrl = "(^https?://(.*)[.:/]*|^$)"

var RegexpDateTime = "^((([0-9]{3}[1-9]|[0-9]{2}[1-9][0-9]{1}|[0-9]{1}[1-9][0-9]{2}|[1-9][0-9]{3})(((0[13578]|1[02])(0[1-9]|[12][0-9]|3[01]))|((0[469]|11)(0[1-9]|[12][0-9]|30))|(02(0[1-9]|[1][0-9]|2[0-8]))))|((([0-9]{2})(0[48]|[2468][048]|[13579][26])|((0[48]|[2468][048]|[3579][26])00))0229))([0-1]?[0-9]|2[0-3])([0-5][0-9])([0-5][0-9])$"

var RegexpMoneyFormat = "^(?:0\\.\\d{0,1}[1-9]|(?!0)\\d{1,6}(?:\\.\\d{0,1}[0-9])?)$"

var RegexpTaiwanPhone = "(\\+886)09\\d{8}"

var RegexpThailandPhone = "(\\+66)0\\d{9}"

var RegexpChinaPhone = "(\\+86)1\\d{10}"

var RegexpVietnamPhone = "(\\+84)0\\d{9}"

var RegexpCambodiaPhone = "(\\+855)0\\d{8}"

/** ====================================API提單==============================================**/
var REGEXP_WALLET_ERC = "[0][x](?=.*[a-zA-Z])(?=.*\\\\d)[a-zA-Z0-9]{40}$"

var REGEXP_WALLET_TRC = "^[T](?=.*[a-zA-Z])(?=.*\\\\d)[a-zA-Z0-9]{33}$"

var REGEXP_BANK_ID = "^[0-9]*$"

var REGEXP_URL = "(https?|http)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]"
