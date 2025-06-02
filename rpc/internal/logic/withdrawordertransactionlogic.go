package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/copo888/transaction_service/common/constants"
	"github.com/copo888/transaction_service/common/errorz"
	"github.com/copo888/transaction_service/common/response"
	"github.com/copo888/transaction_service/common/utils"
	"github.com/copo888/transaction_service/rpc/internal/model"
	"github.com/copo888/transaction_service/rpc/internal/service/merchantbalanceservice"
	"github.com/copo888/transaction_service/rpc/internal/types"
	"github.com/copo888/transaction_service/rpc/transactionclient"
	"github.com/gioco-play/easy-i18n/i18n"
	"github.com/neccoys/go-zero-extension/redislock"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/copo888/transaction_service/rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawOrderTransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWithdrawOrderTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawOrderTransactionLogic {
	return &WithdrawOrderTransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WithdrawOrderTransactionLogic) WithdrawOrderTransaction(in *transactionclient.WithdrawOrderRequest) (*transactionclient.WithdrawOrderResponse, error) {

	myDB := l.svcCtx.MyDB

	transferAmount := utils.FloatAddCWithTrancated(in.OrderAmount, in.HandlingFee, in.CurrencyCode)
	merchantOrderNo := "COPO_" + in.OrderNo
	if len(in.MerchantOrderNo) > 0 {
		merchantOrderNo = in.MerchantOrderNo
	}
	logx.WithContext(l.ctx).Infof("下发单交易初始化： %v, %v, %v", in.String())
	// 依商户是否给回调网址，决定是否回调商户flag
	var isMerchantCallback string //0：否、1:是、2:不需回调
	if len(in.NotifyUrl) > 0 {
		isMerchantCallback = constants.MERCHANT_CALL_BACK_NO
	} else {
		isMerchantCallback = constants.MERCHANT_CALL_BACK_DONT_USE
	}

	// 初始化订单
	txOrder := &types.Order{
		MerchantCode:         in.MerchantCode,
		MerchantOrderNo:      merchantOrderNo,
		OrderNo:              in.OrderNo,
		Type:                 constants.ORDER_TYPE_XF,
		Status:               constants.WAIT_PROCESS,
		IsMerchantCallback:   isMerchantCallback,
		IsLock:               constants.IS_LOCK_NO,
		IsCalculateProfit:    constants.IS_CALCULATE_PROFIT_NO,
		IsTest:               constants.IS_TEST_NO, //是否測試單
		PersonProcessStatus:  constants.PERSON_PROCESS_STATUS_NO_ROCESSING,
		CreatedBy:            in.UserAccount,
		UpdatedBy:            in.UserAccount,
		BalanceType:          "XFB",
		TransferAmount:       transferAmount,
		OrderAmount:          in.OrderAmount,
		TransferHandlingFee:  in.HandlingFee,
		HandlingFee:          in.HandlingFee,
		MerchantBankAccount:  in.MerchantBankeAccount,
		MerchantBankNo:       in.MerchantBankNo,
		MerchantBankProvince: in.MerchantBankProvince,
		MerchantBankCity:     in.MerchantBankCity,
		MerchantBankName:     in.MerchantBankName,
		MerchantAccountName:  in.MerchantAccountName,
		CurrencyCode:         in.CurrencyCode,
		Source:               in.Source,
		PageUrl:              in.PageUrl,
		NotifyUrl:            in.NotifyUrl,
		Memo:                 in.Memo,
		ChangeType:           in.ChangeType,
		IsUsdt:               in.IsUsdt,
		DefaultHandlingFee:   in.HandlingFee,
	}
	if len(in.MerchantOrderNo) > 0 {
		txOrder.MerchantOrderNo = in.MerchantOrderNo
	}
	redisKey := fmt.Sprintf("%s-%s", txOrder.MerchantCode, txOrder.CurrencyCode)
	redisLock := redislock.New(l.svcCtx.RedisClient, redisKey, "merchant-balance:")
	redisLock.SetExpire(8)
	if isOK, redisErr := redisLock.TryLockTimeout(8); isOK {
		defer redisLock.Release()

		txDB := l.svcCtx.MyDB.Begin()
		//新增收支记录，更新商户余额
		updateBalance := types.UpdateBalance{
			MerchantCode:    txOrder.MerchantCode,
			CurrencyCode:    txOrder.CurrencyCode,
			OrderNo:         txOrder.OrderNo,
			MerchantOrderNo: txOrder.MerchantOrderNo,
			OrderType:       txOrder.Type,
			TransactionType: "15",
			BalanceType:     txOrder.BalanceType,
			TransferAmount:  txOrder.TransferAmount,
			CreatedBy:       in.UserAccount,
			Comment:         txOrder.Memo,
			MerPtBalanceId:  in.PtBalanceId,
		}
		if in.Source == constants.API {
			isBlock, _ := model.NewBankBlockAccount(txDB).CheckIsBlockAccount(txOrder.MerchantBankAccount)
			if isBlock { //银行账号为黑名单
				logx.WithContext(l.ctx).Infof("交易账户%s-%s在黑名单内", txOrder.MerchantAccountName, txOrder.MerchantBankNo)
				updateBalance.TransferAmount = 0                           // 使用0元前往钱包扣款
				txOrder.ErrorType = constants.ERROR6_BANK_ACCOUNT_IS_BLACK //交易账户为黑名单
				txOrder.ErrorNote = constants.BANK_ACCOUNT_IS_BLACK        //失败原因：黑名单交易失败
				txOrder.Status = constants.FAIL                            //状态:失败
				txOrder.Fee = 0                                            //写入本次手续费(未发送到渠道的交易，都设为0元)
				txOrder.HandlingFee = 0
				//transAt = types.JsonTime{}.New()
				logx.WithContext(l.ctx).Infof("商户 %s，代付订单 %#v ，交易账户为黑名单", txOrder.MerchantCode, txOrder)
			}
		}

		//更新子钱包且新增商户子钱包异动记录
		if in.PtBalanceId > 0 {
			merchantPtBalanceRecord, errS := merchantbalanceservice.UpdateXF_Pt_Balance_Debit(l.ctx, txDB, &updateBalance)
			if errS != nil {
				logx.WithContext(l.ctx).Errorf("商户:%s，更新子錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, errS.Error(), updateBalance)
				txDB.Rollback()
				return &transactionclient.WithdrawOrderResponse{
					Code:    response.SYSTEM_ERROR,
					Message: "钱包异动失败",
					OrderNo: txOrder.OrderNo,
				}, nil
			} else {
				logx.WithContext(l.ctx).Infof("下发提单 %s，子錢包扣款成功", merchantPtBalanceRecord.OrderNo)
			}
		}

		//更新钱包且新增商户钱包异动记录
		merchantBalanceRecord, err1 := merchantbalanceservice.UpdateXFBalance_Debit(l.ctx, txDB, &updateBalance)
		if err1 != nil {
			logx.WithContext(l.ctx).Errorf("商户:%s，更新錢包紀錄錯誤:%s, updateBalance:%#v", updateBalance.MerchantCode, err1.Error(), updateBalance)
			//TODO  IF 更新钱包错误是response.DATABASE_FAILURE THEN return SYSTEM_ERROR
			txDB.Rollback()
			return &transactionclient.WithdrawOrderResponse{
				Code:    response.SYSTEM_ERROR,
				Message: "钱包异动失败",
				OrderNo: txOrder.OrderNo,
			}, nil
		} else {
			logx.WithContext(l.ctx).Infof("下发提单 %s，錢包扣款成功", merchantBalanceRecord.OrderNo)
			txOrder.BeforeBalance = merchantBalanceRecord.BeforeBalance // 商戶錢包異動紀錄
			txOrder.Balance = merchantBalanceRecord.AfterBalance
		}

		// 创建订单
		if err3 := txDB.Table("tx_orders").Create(&types.OrderX{
			Order: *txOrder,
		}).Error; err3 != nil {
			logx.WithContext(l.ctx).Errorf("新增下发提单失败，商户号: %s, 订单号: %s, err : %s", txOrder.MerchantCode, txOrder.OrderNo, err3.Error())
			txDB.Rollback()
			return &transactionclient.WithdrawOrderResponse{
				Code:    response.DATABASE_FAILURE,
				Message: "数据库错误 tx_orders Create",
				OrderNo: txOrder.OrderNo,
			}, nil
		}

		//計算商戶利潤（不报错）
		calculateProfit := types.CalculateProfit{
			MerchantCode: txOrder.MerchantCode,
			OrderNo:      txOrder.OrderNo,
			Type:         txOrder.Type,
			CurrencyCode: txOrder.CurrencyCode,
			BalanceType:  txOrder.BalanceType,
			OrderAmount:  txOrder.ActualAmount,
		}
		// Xuan 這裡使用 tx or ？
		if err4 := l.calculateOrderProfit(txDB, calculateProfit, in.HandlingFee); err4 != nil {
			logx.WithContext(l.ctx).Errorf("计算下发利润失败，商户号: %s, 订单号: %s, err : %s", txOrder.MerchantCode, txOrder.OrderNo, err4.Error())
		}

		if err4 := txDB.Commit().Error; err4 != nil {
			logx.WithContext(l.ctx).Errorf("最终新增下发提单失败，商户号: %s, 订单号: %s, err : %s", txOrder.MerchantCode, txOrder.OrderNo, err4.Error())
			txDB.Rollback()
		}

	} else {
		logx.WithContext(l.ctx).Errorf("商户钱包处理中，Err:%s。 %s", redisErr.Error(), redisKey)
		return &transactionclient.WithdrawOrderResponse{
			Code:    response.BALANCE_PROCESSING,
			Message: i18n.Sprintf(response.BALANCE_PROCESSING),
		}, nil
	}

	// 新單新增訂單歷程 (不抱錯) TODO: 異步??
	if err5 := myDB.Table("tx_order_actions").Create(&types.OrderActionX{
		OrderAction: types.OrderAction{
			OrderNo:     txOrder.OrderNo,
			Action:      "PLACE_ORDER",
			UserAccount: in.MerchantCode,
			Comment:     "",
		},
	}).Error; err5 != nil {
		logx.WithContext(l.ctx).Error("紀錄訂單歷程出錯:%s", err5.Error())
	}

	return &transactionclient.WithdrawOrderResponse{
		Code:    response.API_SUCCESS,
		Message: "操作成功",
		OrderNo: txOrder.OrderNo,
	}, nil
}

func (l *WithdrawOrderTransactionLogic) calculateOrderProfit(db *gorm.DB, calculateProfit types.CalculateProfit, handlingFee float64) (err error) {
	var merchant *types.Merchant
	var agentLayerCode string
	var agentParentCode string

	// 1. 不是算系統利潤時 要取當前計算商戶(或代理商戶)
	if calculateProfit.MerchantCode != "00000000" {
		if err = db.Table("mc_merchants").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("code = ?", calculateProfit.MerchantCode).Take(&merchant).Error; err != nil {
			return errorz.New(response.DATABASE_FAILURE, err.Error())
		}
		agentLayerCode = merchant.AgentLayerCode
		agentParentCode = merchant.AgentParentCode
	}

	// 2. 設定初始資料
	orderFeeProfit := types.OrderFeeProfit{
		OrderNo:             calculateProfit.OrderNo,
		MerchantCode:        calculateProfit.MerchantCode,
		AgentLayerNo:        agentLayerCode,
		AgentParentCode:     agentParentCode,
		BalanceType:         calculateProfit.BalanceType,
		Fee:                 0,
		HandlingFee:         handlingFee,
		TransferHandlingFee: handlingFee,
		ProfitAmount:        0,
	}

	// 3.设定商户费率
	var merchantCurrency *types.MerchantCurrency
	if err = db.Table("mc_merchant_currencies").
		Where("merchant_code = ? AND currency_code = ?", calculateProfit.MerchantCode, calculateProfit.CurrencyCode).
		Take(&merchantCurrency).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errorz.New(response.RATE_NOT_CONFIGURED, err.Error())
	} else if err != nil {
		return errorz.New(response.DATABASE_FAILURE, err.Error())
	}
	//orderFeeProfit.HandlingFee = merchantCurrency.WithdrawHandlingFee
	//  交易手續費總額 = 訂單金額 + 手續費
	//orderFeeProfit.TransferHandlingFee =
	//	utils.FloatAddC(calculateProfit.OrderAmount, orderFeeProfit.HandlingFee)

	// 4. 除存费率
	if err = db.Table("tx_orders_fee_profit").Create(&types.OrderFeeProfitX{
		OrderFeeProfit: orderFeeProfit,
	}).Error; err != nil {
		return errorz.New(response.DATABASE_FAILURE, err.Error())
	}
	return nil
}
