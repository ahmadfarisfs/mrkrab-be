package store

import (
	"errors"
	"time"

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
	//transaksi tiap proyek
	//info mutasi setiap akun under proyek
	//1. find all pockets under projects
	//2 find all ackun id of pockets
	//3 . find akun id of projects
	//4. find in mutation table yang punya id tersebut semua, join dengan tabel transaksu

	ret := []model.Transaction{}
	var count int64
	//query builder

	initQuery := ac.db.Preload("Mutation").Preload("Mutation.Account")
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

//CreateTransaction should not be used when using accrual basis !
func (ac *TransactionStore) CreateTransaction(accountID int, amount int, remarks string, SoD string, trxTime time.Time) (model.Transaction, error) {
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

			// return model.Transaction{}, errors.New("Account does not have enough balance")
		}
	}
	err = ac.db.Transaction(func(tx *gorm.DB) error {
		//create entry in transaction db
		trxCode := uuid.New().String()
		trxEntry := model.Transaction{
			Remarks:         remarks,
			TransactionCode: trxCode,
			TransactionDate: trxTime,
		}
		if err := tx.Create(&trxEntry).Error; err != nil {
			return err
		}
		//create enty account cred
		if err := tx.Create(&model.Mutation{
			AccountID:     accountID,
			Amount:        amount,
			TransactionID: int(trxEntry.ID),
			// SoD:           SoD,
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

//CreateBankTransaction only transfer bank, subset of create project
func (ac *TransactionStore) CreateBankTransaction(sourceBankAccountID uint, destBankAccountID uint, amount uint, remarks string, meta string, notes string, trxDate time.Time) (model.Transaction, error) {
	accountIDs := []uint{sourceBankAccountID, destBankAccountID}
	resAccounts := []model.Account{}
	if err := ac.db.Model(&model.Account{}).Where("id IN (?)", accountIDs).Find(&resAccounts).Error; err != nil {
		return model.Transaction{}, err
	}
	if len(resAccounts) != 2 {
		return model.Transaction{}, errors.New("Invalid accountID")
	}

	//another sanity check
	for _, v := range resAccounts {
		if (v.ID == sourceBankAccountID || v.ID == destBankAccountID) && v.AccountType != "BANK" {
			return model.Transaction{}, errors.New("Source account bank ID type must be BANK")
		}
	}
	//create entry in ledger
	trxCode := uuid.New().String()
	err := ac.db.Transaction(func(tx *gorm.DB) error {
		//create entry in transaction db. one transaction for both bank and project

		trxEntry := model.Transaction{
			Remarks:         remarks,
			TransactionCode: trxCode,
			Notes:           notes,
			Meta:            meta,
			TransactionType: "BANK",
			TransactionDate: trxDate,
		}
		if err := tx.Create(&trxEntry).Error; err != nil {
			return err
		}

		//create mutation for bank
		mutationEntryBankSource := model.Mutation{
			AccountID:       int(sourceBankAccountID),
			TransactionID:   int(trxEntry.ID),
			TransactionCode: trxCode,
			Amount:          -int(amount),
			Meta:            meta,
			TransactionType: "BANK",
		}
		if err := tx.Create(&mutationEntryBankSource).Error; err != nil {
			return err
		}
		mutationEntryBankTarget := model.Mutation{
			AccountID:       int(destBankAccountID),
			TransactionID:   int(trxEntry.ID),
			TransactionCode: trxCode,
			Amount:          int(amount),
			Meta:            meta,
			TransactionType: "BANK",
		}
		if err := tx.Create(&mutationEntryBankTarget).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.Transaction{}, err
	}
	ret, err := ac.GetTransactionDetailsbyCode(trxCode)
	if err != nil {
		return model.Transaction{}, err
	}
	return ret, nil
}

//CreateProjectTransaction use when income or expense happen
func (ac *TransactionStore) CreateProjectTransaction(sourceIncomeAccountID uint, targetAccountID uint, sourceBankAccountID uint, destBankAccountID uint, amount uint, remarks string, meta string, notes string, trxDate time.Time) (model.Transaction, error) {
	accountIDs := []uint{sourceIncomeAccountID, targetAccountID, sourceBankAccountID, destBankAccountID}
	resAccounts := []model.Account{}
	if err := ac.db.Model(&model.Account{}).Where("id IN (?)", accountIDs).Find(&resAccounts).Error; err != nil {
		return model.Transaction{}, err
	}
	if len(resAccounts) != 4 {
		return model.Transaction{}, errors.New("Invalid accountID")
	}
	//check source and dest type
	fromAccountType := ""
	toAccountType := ""
	for _, v := range resAccounts {
		if v.ID == sourceIncomeAccountID {
			fromAccountType = v.AccountType
		}
		if v.ID == targetAccountID {
			toAccountType = v.AccountType
		}
	}
	//another sanity check
	for _, v := range resAccounts {
		if (v.ID == sourceBankAccountID || v.ID == destBankAccountID) && v.AccountType != "BANK" {
			return model.Transaction{}, errors.New("Source account bank ID type must be BANK")
		}
		if v.ID == sourceIncomeAccountID && v.AccountType == "PROJECT" {
			//expense from project
			//target account must be expense account
			if toAccountType != "EXPENSE" {
				return model.Transaction{}, errors.New("Target account must be expense account")
			}
		}
		if v.ID == targetAccountID && v.AccountType == "PROJECT" {
			//income to project
			//source account must be INCOME account
			if fromAccountType != "INCOME" {
				return model.Transaction{}, errors.New("From account must be income account")
			}

		}

	}
	//create entry in ledger
	trxCode := uuid.New().String()

	err := ac.db.Transaction(func(tx *gorm.DB) error {
		//create entry in transaction db. one transaction for both bank and project

		trxEntry := model.Transaction{
			Remarks:         remarks,
			TransactionCode: trxCode,
			Notes:           notes,
			Meta:            meta,
			TransactionType: "BOTH",
			TransactionDate: trxDate,
		}
		if err := tx.Create(&trxEntry).Error; err != nil {
			return err
		}

		//create entry in mutation for project
		mutationEntrySource := model.Mutation{
			AccountID:       int(sourceIncomeAccountID),
			TransactionID:   int(trxEntry.ID),
			TransactionCode: trxCode,
			Amount:          -int(amount),
			Meta:            meta,
			TransactionType: "PROJECT",
		}
		if err := tx.Create(&mutationEntrySource).Error; err != nil {
			return err
		}
		mutationEntryTarget := model.Mutation{
			AccountID:       int(targetAccountID),
			TransactionID:   int(trxEntry.ID),
			TransactionCode: trxCode,
			Amount:          int(amount),
			Meta:            meta,
			TransactionType: "PROJECT",
		}
		if err := tx.Create(&mutationEntryTarget).Error; err != nil {
			return err
		}

		//create mutation for bank
		mutationEntryBankSource := model.Mutation{
			AccountID:       int(sourceBankAccountID),
			TransactionID:   int(trxEntry.ID),
			TransactionCode: trxCode,
			Amount:          -int(amount),
			Meta:            meta,
			TransactionType: "BANK",
		}
		if err := tx.Create(&mutationEntryBankSource).Error; err != nil {
			return err
		}
		mutationEntryBankTarget := model.Mutation{
			AccountID:       int(destBankAccountID),
			TransactionID:   int(trxEntry.ID),
			TransactionCode: trxCode,
			Amount:          int(amount),
			Meta:            meta,
			TransactionType: "BANK",
		}
		if err := tx.Create(&mutationEntryBankTarget).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.Transaction{}, err
	}
	ret, err := ac.GetTransactionDetailsbyCode(trxCode)
	if err != nil {
		return model.Transaction{}, err
	}
	return ret, nil
}
