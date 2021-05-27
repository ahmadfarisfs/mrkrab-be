package store

import (
	"errors"

	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
	"gorm.io/gorm"
)

type BankAccountStore struct {
	db *gorm.DB
}

func NewBankAccountStore(db *gorm.DB) *BankAccountStore {
	return &BankAccountStore{
		db: db,
	}
}

func (ac *BankAccountStore) ListAccount(req utils.CommonRequest) ([]model.BankAccount, int, error) {
	ret := []model.BankAccount{}
	var count int64
	query := ac.db

	err := query.Model(&ret).Count(&count).Error
	if err != nil {
		return ret, int(count), err
	}
	err = utils.AppendCommonRequest(query, req).Find(&ret).Error
	return ret, int(count), err
}

func (ac *BankAccountStore) CreateAccount(bankName string, holderName string, accountNumber string, description string, accountType string) (model.BankAccount, error) {
	if accountType != "internal" && accountType != "external" {
		return model.BankAccount{}, errors.New("invalid account type")
	}
	ret := model.BankAccount{
		BankName:        bankName,
		HolderName:      holderName,
		Description:     description,
		BankNumber:      accountNumber,
		BankAccountType: accountType,
	}

	err := ac.db.Model(&model.BankAccount{}).Create(&ret).Error
	if err != nil {
		return model.BankAccount{}, err
	}
	return ac.GetAccountDetails(int(ret.ID))
}

func (ac *BankAccountStore) GetAccountDetails(id int) (model.BankAccount, error) {
	ret := model.BankAccount{}
	err := ac.db.Preload("Parent").First(&ret, "id = ?", id).Error
	return ret, err
}

func (ac *BankAccountStore) UpdateAccount(id int, bankName string, holderName string, accountNumber string, description string, accountType string) (model.BankAccount, error) {
	account, err := ac.GetAccountDetails(id)
	if err != nil {
		return account, err
	}
	account.BankName = bankName
	account.HolderName = holderName
	account.BankNumber = accountNumber
	account.Description = description
	account.BankAccountType = accountType
	err = ac.db.Model(&model.BankAccount{}).Updates(&account).Error
	if err != nil {
		return model.BankAccount{}, err
	}

	return account, nil
}

func (ac BankAccountStore) DeleteAccount(id int) error {
	return ac.db.Delete(&model.BankAccount{}, id).Error
}
