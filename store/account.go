package store

import (
	"errors"

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

func (ac *AccountStore) CreateAccount(name string, parentID *uint) (model.Account, error) {
	ret := model.Account{
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
	return ac.GetAccountDetails(int(ret.ID))
}

func (ac *AccountStore) GetAccountDetails(id int) (model.Account, error) {
	ret := model.Account{}
	err := ac.db.Preload("Parent").First(&ret, "id = ?", id).Error
	return ret, err
}
