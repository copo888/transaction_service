package orderfeeprofitservice

import (
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"gorm.io/gorm"
)

// CalculateSubOrderProfit 從舊單產生新單時 計算利潤要沿用舊單費率 (EX: 補單,追回)
func CalculateSubOrderProfit(db *gorm.DB, calculateProfit types.CalculateSubOrderProfit, currency string) (err error) {
	var oldFeeProfits []types.OrderFeeProfit
	newFeeProfits := make([]types.OrderFeeProfit, 2)
	// 取得舊單費率列表 按 agent_layer_no 降序排序 (第一筆是自己 中間幾筆代理 最後一筆系統利潤)
	if err = db.Table("tx_orders_fee_profit").Where("order_no = ?", calculateProfit.OldOrderNo).Order("agent_layer_no desc").Find(&oldFeeProfits).Error; err != nil {
		return errorz.New(response.DATABASE_FAILURE, err.Error())
	}

	if len(oldFeeProfits) < 2 {
		return errorz.New(response.DATABASE_FAILURE, "原單利潤錯誤")
	}

	if calculateProfit.IsCalculateCommission {
		newFeeProfits = oldFeeProfits
	} else {
		//不計算傭金的話 只取第一筆(本身) 第二筆(利潤)
		newFeeProfits[0] = oldFeeProfits[0]
		newFeeProfits[1] = oldFeeProfits[len(oldFeeProfits)-1]
	}

	for i, profit := range newFeeProfits {

		newFeeProfits[i].ID = 0
		// 將舊單利潤替換成新單去做改變
		newFeeProfits[i].OrderNo = calculateProfit.NewOrderNo

		// 交易手續費總額 = 訂單金額 / 100 * 費率 + 手續費
		newFeeProfits[i].TransferHandlingFee = utils.FloatAddCWithTrancated(utils.FloatMulC(utils.FloatDivC(calculateProfit.OrderAmount, 100, currency), profit.Fee, currency), profit.HandlingFee, currency)

		if i == 0 {
			// 第一筆沒有傭金利潤
			newFeeProfits[i].ProfitAmount = 0
		} else {
			//  計算利潤 (上一筆交易手續費 - 當前交易手續費 = 利潤(或傭金))
			newFeeProfits[i].ProfitAmount = newFeeProfits[i-1].TransferHandlingFee - newFeeProfits[i].TransferHandlingFee
		}

		// 保存利潤
		if err = db.Table("tx_orders_fee_profit").Create(&types.OrderFeeProfitX{
			OrderFeeProfit: newFeeProfits[i],
		}).Error; err != nil {
			return errorz.New(response.DATABASE_FAILURE, err.Error())
		}
	}

	return
}
