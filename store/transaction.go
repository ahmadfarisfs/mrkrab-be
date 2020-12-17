package store

import (
	"errors"

	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionStore struct {
	db *gorm.DB
}

func NewTransactionStore(db *gorm.DB) *TransactionStore {
	return &TransactionStore{
		db: db,
	}
}

func (ac *TransactionStore) ListTransaction(req utils.CommonRequest) ([]model.Transaction, int, error) {
	ret := []model.Transaction{}
	var count int64
	//query builder

	initQuery := ac.db.Preload("Mutation")
	//	Where("created_date BETWEEN ? and ?", startTime, endTime).
	//Where("account_id in (?)", accountID)
	//count total data
	err := initQuery.Model(&ret).Count(&count).Error
	if err != nil {
		return ret, int(count), err
	}
	//actually fetch data with limit and offset
	//err = initQuery.
	err = utils.AppendCommonRequest(initQuery, req).
		//	Offset(offset).Limit(limit).
		Find(&ret).Error
	return ret, int(count), err
}

func (ac *TransactionStore) GetTransactionDetailsbyID(transactionID int) (model.Transaction, error) {
	ret := model.Transaction{}
	err := ac.db.Preload("Mutation").Preload("Mutation.Transaction").First(&ret, "id = ?", transactionID).Error
	return ret, err
}

func (ac *TransactionStore) GetTransactionDetailsbyCode(transactionCode string) (model.Transaction, error) {
	ret := model.Transaction{}
	err := ac.db.Preload("Mutation").First(&ret, "transaction_code = ?", transactionCode).Error
	return ret, err
}

func (ac *TransactionStore) CreateTransaction(accountID int, amount int, remarks string) (model.Transaction, error) {
	var transactionID int
	//check account must be valid
	ret := model.Account{}
	err := ac.db.Model(&model.Account{}).First(&ret, "id = ?", accountID).Error
	if err != nil {
		return model.Transaction{}, err
	}
	if amount < 0 {
		if ret.Balance+amount < 0 {
			//cannot do
			return model.Transaction{}, errors.New("Account does not have enough balance")
		}
	}
	err = ac.db.Transaction(func(tx *gorm.DB) error {
		//create entry in transaction db
		trxCode := uuid.New().String()
		trxEntry := model.Transaction{Remarks: remarks, TransactionCode: trxCode}
		if err := tx.Create(&trxEntry).Error; err != nil {
			return err
		}
		//create enty account cred
		if err := tx.Create(&model.Mutation{
			AccountID:     accountID,
			Amount:        amount,
			TransactionID: int(trxEntry.ID),
		}).Error; err != nil {
			return err
		}

		transactionID = int(trxEntry.ID)
		// return nil will commit the whole transaction
		return nil
	})
	if err != nil {
		return model.Transaction{}, err
	}
	return ac.GetTransactionDetailsbyID(transactionID)
}

func (ac *TransactionStore) CreateTransfer(accountFrom int, accountTo int, amount uint, remarks string) (model.Transaction, error) {
	var transactionID int

	//check both account must be valid
	ret := []model.Account{}
	err := ac.db.Model(&model.Account{}).Find(&ret, "id IN (?)", []int{accountFrom, accountTo}).Error
	if err != nil {
		return model.Transaction{}, err
	}
	if len(ret) != 2 {
		return model.Transaction{}, errors.New("Invalid Account ID")
	}
	//check amount
	for _, v := range ret {
		if v.ID == uint(accountFrom) {
			if v.Balance < int(amount) {
				return model.Transaction{}, errors.New("Source Account does not have enough balance")
			}
		}
	}
	err = ac.db.Transaction(func(tx *gorm.DB) error {
		//create entry in transaction db
		trxCode := uuid.New().String()
		trxEntry := model.Transaction{Remarks: remarks, TransactionCode: trxCode}
		if err := tx.Create(&trxEntry).Error; err != nil {
			return err
		}
		//create enty account cred
		if err := tx.Create(&model.Mutation{
			AccountID:     accountFrom,
			Amount:        -int(amount),
			TransactionID: int(trxEntry.ID),
		}).Error; err != nil {
			return err
		}
		//create enty account deb
		if err := tx.Create(&model.Mutation{
			AccountID:     accountTo,
			Amount:        int(amount),
			TransactionID: int(trxEntry.ID),
		}).Error; err != nil {
			return err
		}
		transactionID = int(trxEntry.ID)
		// return nil will commit the whole transaction
		return nil
	})
	if err != nil {
		return model.Transaction{}, err
	}

	return ac.GetTransactionDetailsbyID(transactionID)
}
