package model

import (
	"github.com/copo888/transaction_service/rpc/internal/types"
	"gorm.io/gorm"
)

type bankBlockAccount struct {
	MyDB  *gorm.DB
	Table string
}

func NewBankBlockAccount(mydb *gorm.DB) *bankBlockAccount {
	return &bankBlockAccount{
		MyDB:  mydb,
		Table: "bk_block_account",
	}
}

func (b *bankBlockAccount) GetAll() (resp []types.BankBlockAccount, err error) {
	err = b.MyDB.Table(b.Table).Find(&resp).Error
	return
}

func (b *bankBlockAccount) CheckIsBlockAccount(account string) (isExist bool, err error) {
	err = b.MyDB.Table(b.Table).
		Select("count(*) > 0").
		Where("bank_account = ?", account).
		Find(&isExist).Error
	return
}
