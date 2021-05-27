package store

import (
	"errors"

	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
	"gorm.io/gorm"
)

type FinancialAccountStore struct {
	db *gorm.DB
}

func NewFinancialAccountStore(db *gorm.DB) *FinancialAccountStore {
	return &FinancialAccountStore{
		db: db,
	}
}

func (ac *FinancialAccountStore) ListAccount(req utils.CommonRequest) ([]model.FinancialAccount, int, error) {
	ret := []model.FinancialAccount{}
	var count int64
	query := ac.db

	err := query.Model(&ret).Count(&count).Error
	if err != nil {
		return ret, int(count), err
	}
	err = utils.AppendCommonRequest(query, req).Find(&ret).Error
	return ret, int(count), err
}

func (ac *FinancialAccountStore) CreateAccount(name string, description string, accountType string, parentID *uint) (model.FinancialAccount, error) {
	if accountType != "internal" && accountType != "external" {
		return model.FinancialAccount{}, errors.New("invalid account type")
	}
	ret := model.FinancialAccount{
		AccountName: name,
		ParentID:    parentID,
		Description: description,
		AccountType: accountType,
	}
	if parentID != nil {
		_, err := ac.GetAccountDetails(int(*parentID))
		if err != nil {
			return model.FinancialAccount{}, errors.New("ParentID does not exist")
		}
	}

	err := ac.db.Model(&model.FinancialAccount{}).Create(&ret).Error
	if err != nil {
		return model.FinancialAccount{}, err
	}
	return ac.GetAccountDetails(int(ret.ID))
}

func (ac *FinancialAccountStore) GetAccountDetails(id int) (model.FinancialAccount, error) {
	ret := model.FinancialAccount{}
	err := ac.db.Preload("Parent").First(&ret, "id = ?", id).Error
	return ret, err
}

func (ac *FinancialAccountStore) UpdateAccount(id int, name string, description string, accountType string) (model.FinancialAccount, error) {
	account, err := ac.GetAccountDetails(id)
	if err != nil {
		return account, err
	}
	account.AccountName = name
	account.Description = description
	account.AccountType = accountType
	err = ac.db.Model(&model.FinancialAccount{}).Updates(&account).Error
	if err != nil {
		return model.FinancialAccount{}, err
	}

	return account, nil
}

func (ac FinancialAccountStore) DeleteAccount(id int) error {
	return ac.db.Delete(&model.FinancialAccount{}, id).Error
}
