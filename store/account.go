package store

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
	"gorm.io/gorm"
)

type AccountStore struct {
	db *gorm.DB
}

func NewAccountStore(db *gorm.DB) *AccountStore {
	return &AccountStore{
		db: db,
	}
}

func (ac *AccountStore) ListAccount(req utils.CommonRequest) ([]model.Account, int, error) {
	ret := []model.Account{}
	var count int64
	query := ac.db

	err := query.Model(&ret).Count(&count).Error
	if err != nil {
		return ret, int(count), err
	}
	err = utils.AppendCommonRequest(query, req).Find(&ret).Error
	return ret, int(count), err
}

func (ac *AccountStore) CreateBankAccount(bankName string,
	bankNumber string,
	bankHoldername string,
	meta string, isInternal bool) (model.BankAccount, error) {
	accountName := "BANK-" + strings.ToUpper(bankHoldername) + "-" + strings.ToUpper(bankName) + "-" + strconv.Itoa(int(time.Now().Unix()))

	retc := model.BankAccount{
		BankName:       bankName,
		BankNumber:     bankNumber,
		BankHolderName: bankHoldername,
		Internal:       isInternal,
		// AccountID:      int(ret.ID),
	}
	if errx := ac.db.Transaction(func(tx *gorm.DB) error {
		ret := model.Account{
			AccountType: "BANK",
			AccountName: accountName,
		}

		err := ac.db.Model(&model.Account{}).Create(&ret).Error
		if err != nil {
			return err
		}
		retc.AccountID = int(ret.ID)
		if err := ac.db.Model(&model.BankAccount{}).Create(&retc).Error; err != nil {
			return err
		}
		return nil
	}); errx != nil {
		return model.BankAccount{}, errx
	}
	return retc, nil
}
func (ac *AccountStore) CreateFinancialAccount(name string, accountType string, meta string, parentID *uint) error {
	ret := model.Account{
		AccountType: accountType,
		AccountName: name,
		ParentID:    parentID,
	}

	if errx := ac.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.Account{}).Create(&ret).Error
		if err != nil {
			return err
		}
		accountName := accountType + "-" + strings.ToUpper(name) + "-" + strconv.Itoa(int(time.Now().Unix()))
		if accountType == "EXPENSE" {
			if err := tx.Model(&model.ExpenseAccount{}).Create(&model.ExpenseAccount{
				Name:      accountName,
				AccountID: int(ret.ID),
			}).Error; err != nil {
				return err
			}
		} else if accountType == "INCOME" {
			if err := tx.Model(&model.IncomeAccount{}).Create(&model.IncomeAccount{
				Name:      accountName,
				AccountID: int(ret.ID),
			}).Error; err != nil {
				return err
			}
		}
		return nil
	}); errx != nil {
		return errx
	}

	return nil

}

func (ac *AccountStore) CreateAccount(name string, accountType string, meta string, parentID *uint) (model.Account, error) {
	ret := model.Account{
		AccountType: accountType,
		AccountName: name,
		ParentID:    parentID,
	}
	if parentID != nil {
		_, err := ac.GetAccountDetails(int(*parentID))
		if err != nil {
			return model.Account{}, errors.New("ParentID does not exist")
		}
	}

	err := ac.db.Model(&model.Account{}).Create(&ret).Error
	if err != nil {
		return model.Account{}, err
	}

	if accountType == "BANK" {

	} else {

	}

	return ac.GetAccountDetails(int(ret.ID))
}

func (ac *AccountStore) GetAccountDetails(id int) (model.Account, error) {
	ret := model.Account{}
	err := ac.db.Preload("Parent").First(&ret, "id = ?", id).Error
	return ret, err
}
