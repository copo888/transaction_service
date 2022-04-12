package merchantbalanceservice

import (
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"gorm.io/gorm"
)

// GetBalanceType
// orderType (代付 DF 支付 ZF 下發 XF 內充 NC)
func GetBalanceType(db *gorm.DB, channelCode, orderType string) (balanceType string, err error) {
	balanceType = "XFB"

	// 支付&下發 一定是異動下發餘額
	if orderType == "ZF" || orderType == "XF" {
		return
	}

	// 取得渠道資訊
	var channel types.ChannelData
	if err = db.Table("ch_channels").Where("code = ?", channelCode).Take(&channel).Error; err != nil {
		return "", errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	// 純代付渠道 代付 內充 異動代付餘額
	if channel.IsProxy == "0" && (orderType == "DF" || orderType == "NC") {
		balanceType = "DFB"
	}

	return
}
