package store

import (
	"time"

	"github.com/ahmadfarisfs/krab-core/contract"
	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionStore struct {
	db *gorm.DB
}

func NewTransactionStore(db *gorm.DB) *TransactionStore {
	return &TransactionStore{
		db: db,
	}
}

func (ts *TransactionStore) CreateIncomeTransaction(amount int, remarks string, projectID int, incomeFinancialAccountID int, sourceBankAccountID int, destinationBankAccountID int, isPaid bool, transactionTime time.Time) (model.Transaction, error) {

	newTransaction := model.Transaction{
		Remarks:         remarks,
		Amount:          amount,
		TransactionType: "income",
		TransactionTime: transactionTime,
		IsPaid:          isPaid,
		TransactionCode: utils.GenerateTrxCode("income"),
	}

	err := ts.db.Transaction(func(tx *gorm.DB) error {
		//create new transaction
		err := tx.Model(&model.Transaction{}).Create(&newTransaction).Error
		if err != nil {
			return err
		}

		//create project - account mutation
		err = tx.Model(&model.FinancialAccountMutation{}).Create(model.FinancialAccountMutation{
			TransactionID:   int(newTransaction.ID),
			ProjectID:       projectID,
			AccountID:       incomeFinancialAccountID,
			Amount:          amount,
			IsPaid:          isPaid,
			TransactionCode: newTransaction.TransactionCode,
		}).Error
		if err != nil {
			return err
		}

		//create bank mutation (-)
		err = tx.Model(&model.BankAccountMutation{}).Create(model.BankAccountMutation{
			TransactionID:   int(newTransaction.ID),
			BankAccountID:   sourceBankAccountID,
			Amount:          -amount,
			IsPaid:          isPaid,
			TransactionCode: newTransaction.TransactionCode,
		}).Error
		if err != nil {
			return err
		}

		//create bank mutation (+)
		err = tx.Model(&model.BankAccountMutation{}).Create(model.BankAccountMutation{
			TransactionID:   int(newTransaction.ID),
			BankAccountID:   destinationBankAccountID,
			Amount:          amount,
			IsPaid:          isPaid,
			TransactionCode: newTransaction.TransactionCode,
		}).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return model.Transaction{}, err
	}

	return newTransaction, nil
}
func (ts *TransactionStore) CreateExpenseTransaction(amount int, remarks string, projectID int, sourceBankAccountID int, expenses []contract.ExpenseDestination, isPaid bool, transactionTime time.Time) (model.Transaction, error) {
	newTransaction := model.Transaction{
		Remarks:         remarks,
		Amount:          amount,
		TransactionType: "expense",
		TransactionTime: transactionTime,
		IsPaid:          isPaid,
		TransactionCode: utils.GenerateTrxCode("expense"),
	}
	err := ts.db.Transaction(func(tx *gorm.DB) error {
		//create new transaction
		err := tx.Create(&newTransaction).Error
		if err != nil {
			return err
		}

		//create bank mutation (-)
		err = tx.Model(&model.BankAccountMutation{}).Create(model.BankAccountMutation{
			TransactionID: int(newTransaction.ID),
			BankAccountID: sourceBankAccountID,
			Amount:        -amount,
			IsPaid:        isPaid,
		}).Error
		if err != nil {
			return err
		}

		for _, expense := range expenses {
			//create bank mutation (+)
			err = tx.Model(&model.BankAccountMutation{}).Create(model.BankAccountMutation{
				TransactionID: int(newTransaction.ID),
				BankAccountID: expense.DestinationBankAccountID,
				Amount:        expense.Amount,
				IsPaid:        isPaid,
			}).Error
			if err != nil {
				return err
			}
			//create project - account mutation
			err = tx.Model(&model.FinancialAccountMutation{}).Create(model.FinancialAccountMutation{
				TransactionID: int(newTransaction.ID),
				ProjectID:     projectID,
				AccountID:     expense.ExpenseFinancialAccountID,
				Amount:        expense.Amount,
				IsPaid:        isPaid,
			}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return newTransaction, err
	}
	return newTransaction, nil
}
func (ts *TransactionStore) CreateBankTransferTransaction(amount int, remarks string, projectID int, sourceBankAccountID int, destination []contract.TransferDestination, transferFee contract.ExpenseDestination, isPaid bool, transactionTime time.Time) (model.Transaction, error) {
	newTransaction := model.Transaction{
		Remarks:         remarks,
		Amount:          amount,
		TransactionType: "btransfer",
		TransactionTime: transactionTime,
		IsPaid:          isPaid,
		TransactionCode: utils.GenerateTrxCode("btransfer"),
	}
	err := ts.db.Transaction(func(tx *gorm.DB) error {
		//create new transaction
		err := tx.Create(&newTransaction).Error
		if err != nil {
			return err
		}

		//create bank mutation (-)
		err = tx.Model(&model.BankAccountMutation{}).Create(model.BankAccountMutation{
			TransactionID: int(newTransaction.ID),
			BankAccountID: sourceBankAccountID,
			Amount:        -amount,
			IsPaid:        isPaid,
		}).Error
		if err != nil {
			return err
		}

		for _, expense := range destination {
			//create bank mutation (+)
			err = tx.Model(&model.BankAccountMutation{}).Create(model.BankAccountMutation{
				TransactionID: int(newTransaction.ID),
				BankAccountID: expense.DestinationBankAccountID,
				Amount:        expense.Amount,
				IsPaid:        isPaid,
			}).Error
			if err != nil {
				return err
			}
		}

		//create transfer fee expense
		err = tx.Model(&model.FinancialAccountMutation{}).Create(model.FinancialAccountMutation{
			TransactionID: int(newTransaction.ID),
			ProjectID:     projectID,
			AccountID:     transferFee.ExpenseFinancialAccountID,
			Amount:        transferFee.Amount,
			IsPaid:        isPaid,
		}).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return newTransaction, err
	}
	return newTransaction, nil
}
func (ts *TransactionStore) CreateProjectTransferTransaction(amount int, remarks string, sourceProjectID int, sourceFinancialAccountID int, destinationProjectID int, destinationFinancialAccountID int, isPaid bool, transactionTime time.Time) (model.Transaction, error) {
	newTransaction := model.Transaction{
		Remarks:         remarks,
		Amount:          amount,
		TransactionType: "ptransfer",
		TransactionTime: transactionTime,
		IsPaid:          isPaid,
		TransactionCode: utils.GenerateTrxCode("ptransfer"),
	}
	err := ts.db.Transaction(func(tx *gorm.DB) error {
		//create new transaction
		err := tx.Create(&newTransaction).Error
		if err != nil {
			return err
		}

		//create project - account mutation (-)
		err = tx.Model(&model.FinancialAccountMutation{}).Create(model.FinancialAccountMutation{
			TransactionID: int(newTransaction.ID),
			ProjectID:     sourceProjectID,
			AccountID:     sourceFinancialAccountID,
			Amount:        -amount,
			IsPaid:        isPaid,
		}).Error
		if err != nil {
			return err
		}

		//create project - account mutation (+)
		err = tx.Model(&model.FinancialAccountMutation{}).Create(model.FinancialAccountMutation{
			TransactionID: int(newTransaction.ID),
			ProjectID:     destinationProjectID,
			AccountID:     destinationFinancialAccountID,
			Amount:        amount,
			IsPaid:        isPaid,
		}).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return newTransaction, err
	}
	return newTransaction, nil
}

func (ts *TransactionStore) UpdateIncomeTransaction(id int, amount int, remarks string, projectID int, incomeFinancialAccountID int, sourceBankAccountID int, destinationBankAccountID int, isPaid bool, transactionTime time.Time) (model.Transaction, error) {
	trx, err := ts.GetTransactionDetailsbyID(id)
	if err != nil {
		return trx, err
	}
	trx.Amount = amount
	trx.Remarks = remarks
	trx.IsPaid = isPaid
	trx.TransactionTime = transactionTime

	err = ts.db.Transaction(func(tx *gorm.DB) error {
		err = ts.db.Model(&model.Transaction{}).Updates(&trx).Error
		if err != nil {
			return err
		}

		// bank source
		{
			//get mutation: bank source mutation
			bankMut := model.BankAccountMutation{}
			err = ts.db.Model(&model.BankAccountMutation{}).Where("transaction_id = ? and amount < 0", trx.ID).First(&bankMut).Error
			if err != nil {
				return err
			}

			//fill new data
			bankMut.Amount = -amount
			bankMut.BankAccountID = sourceBankAccountID
			bankMut.IsPaid = isPaid

			//update the mutation
			err = ts.db.Model(&model.BankAccountMutation{}).Updates(&bankMut).Error
			if err != nil {
				return err
			}
		}

		//destination bank
		{
			//get mutation: bank destination mutation
			bankMutDest := model.BankAccountMutation{}
			err = ts.db.Model(&model.BankAccountMutation{}).Where("transaction_id = ? and amount > 0", trx.ID).First(&bankMutDest).Error
			if err != nil {
				return err
			}

			//fill new data
			bankMutDest.Amount = amount
			bankMutDest.BankAccountID = destinationBankAccountID
			bankMutDest.IsPaid = isPaid

			//update the mutation
			err = ts.db.Model(&model.BankAccountMutation{}).Updates(&bankMutDest).Error
			if err != nil {
				return err
			}
		}

		//destination project - account
		{
			//get mutation
			finMut := model.FinancialAccountMutation{}
			err = ts.db.Model(&model.FinancialAccountMutation{}).Where("transaction_id = ?", trx.ID).First(&finMut).Error
			if err != nil {
				return err
			}

			//fill new data
			finMut.Amount = amount
			finMut.AccountID = incomeFinancialAccountID
			finMut.IsPaid = isPaid

			//update the mutation
			err = ts.db.Model(&model.FinancialAccountMutation{}).Updates(&finMut).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	return model.Transaction{}, nil
}
func (ts *TransactionStore) UpdateExpenseTransaction(id int, amount int, remarks string, projectID int, sourceBankAccountID int, expenses []contract.ExpenseDestination, isPaid bool, transactionTime time.Time) (model.Transaction, error) {
	// trx, err := ts.GetTransactionDetailsbyID(id)
	// if err != nil {
	// 	return trx, err
	// }
	// trx.Amount = amount
	// trx.Remarks = remarks
	// trx.IsPaid = isPaid
	// trx.TransactionTime = transactionTime
	// trx.AccountMutation =

	// err = ts.db.Transaction(func(tx *gorm.DB) error {
	// 	tx
	// })
	return model.Transaction{}, nil
}

func (ts *TransactionStore) UpdateBankTransferTransaction(id int, amount int, remarks string, sourceBankAccountID int, destination []contract.TransferDestination, transferFee contract.ExpenseDestination, isPaid bool, transactionTime time.Time) (trx model.Transaction, err error) {
	return
}
func (ts *TransactionStore) UpdateFinancialAccountTransferTransaction(id int, amount int, remarks string, sourceFinancialAccountID int, destinationFinancialAccountID int, isPaid bool, transactionTime time.Time) (trx model.Transaction, err error) {
	return
}

func (ts *TransactionStore) UpdateTransactionRaw(id int, trx model.Transaction) (model.Transaction, error) {
	trx, err := ts.GetTransactionDetailsbyID(id)
	if err != nil {
		return trx, err
	}
	err = ts.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&trx).Error
	if err != nil {
		return trx, err
	}
	return trx, nil
}

func (ts *TransactionStore) DeleteTransaction(id int) error {
	return ts.db.Model(&model.Transaction{}).Delete(&model.Transaction{}, id).Error
}

func (ts *TransactionStore) GetTransactionDetailsbyID(transactionID int) (model.Transaction, error) {
	trx := model.Transaction{}
	err := ts.db.Model(&model.Transaction{}).Preload(clause.Associations).Where("id= ?", transactionID).First(&trx).Error
	return trx, err
}
func (ts *TransactionStore) GetTransactionDetailsbyCode(transactionCode string) (model.Transaction, error) {
	trx := model.Transaction{}
	err := ts.db.Model(&model.Transaction{}).Where("transaction_code= ?", transactionCode).First(&trx).Error
	return trx, err
}

func (ts *TransactionStore) ListTransaction(req utils.CommonRequest) ([]model.Transaction, int, error) {
	ret := []model.Transaction{}
	var count int64
	query := ts.db

	err := query.Model(&ret).Count(&count).Error
	if err != nil {
		return ret, int(count), err
	}
	err = utils.AppendCommonRequest(query, req).Find(&ret).Error
	return ret, int(count), err
}

// func (ac *TransactionStore) ListTransaction(req utils.CommonRequest) ([]model.Transaction, int, error) {
// 	//transaksi tiap proyek
// 	//info mutasi setiap akun under proyek
// 	//1. find all pockets under projects
// 	//2 find all ackun id of pockets
// 	//3 . find akun id of projects
// 	//4. find in mutation table yang punya id tersebut semua, join dengan tabel transaksu

// 	ret := []model.Transaction{}
// 	var count int64
// 	//query builder

// 	initQuery := ac.db.Preload("Mutation").Preload("Mutation.Account")
// 	//	Where("created_date BETWEEN ? and ?", startTime, endTime).
// 	//Where("account_id in (?)", accountID)
// 	//count total data
// 	err := initQuery.Model(&ret).Count(&count).Error
// 	if err != nil {
// 		return ret, int(count), err
// 	}
// 	//actually fetch data with limit and offset
// 	//err = initQuery.
// 	err = utils.AppendCommonRequest(initQuery, req).
// 		//	Offset(offset).Limit(limit).
// 		Find(&ret).Error
// 	return ret, int(count), err
// }

// func (ac *TransactionStore) GetTransactionDetailsbyID(transactionID int) (model.Transaction, error) {
// 	ret := model.Transaction{}
// 	err := ac.db.Preload("Mutation").Preload("Mutation.Transaction").First(&ret, "id = ?", transactionID).Error
// 	return ret, err
// }

// func (ac *TransactionStore) GetTransactionDetailsbyCode(transactionCode string) (model.Transaction, error) {
// 	ret := model.Transaction{}
// 	err := ac.db.Preload("Mutation").First(&ret, "transaction_code = ?", transactionCode).Error
// 	return ret, err
// }

// //CreateTransaction should not be used when using accrual basis !
// func (ac *TransactionStore) CreateTransaction(accountID int, amount int, remarks string, SoD string, trxTime time.Time) (model.Transaction, error) {
// 	var transactionID int
// 	//check account must be valid
// 	ret := model.FinancialAccount{}
// 	err := ac.db.Model(&model.FinancialAccount{}).First(&ret, "id = ?", accountID).Error
// 	if err != nil {
// 		return model.Transaction{}, err
// 	}
// 	if amount < 0 {
// 		if ret.Balance+amount < 0 {
// 			//cannot do

// 			// return model.Transaction{}, errors.New("Account does not have enough balance")
// 		}
// 	}
// 	err = ac.db.Transaction(func(tx *gorm.DB) error {
// 		//create entry in transaction db
// 		trxCode := uuid.New().String()
// 		trxEntry := model.Transaction{
// 			Remarks:         remarks,
// 			TransactionCode: trxCode,
// 			TransactionDate: trxTime,
// 		}
// 		if err := tx.Create(&trxEntry).Error; err != nil {
// 			return err
// 		}
// 		//create enty account cred
// 		if err := tx.Create(&model.Mutation{
// 			AccountID:     accountID,
// 			Amount:        amount,
// 			TransactionID: int(trxEntry.ID),
// 			SoD:           SoD,
// 		}).Error; err != nil {
// 			return err
// 		}

// 		transactionID = int(trxEntry.ID)
// 		// return nil will commit the whole transaction
// 		return nil
// 	})
// 	if err != nil {
// 		return model.Transaction{}, err
// 	}
// 	return ac.GetTransactionDetailsbyID(transactionID)
// }

// func (ac *TransactionStore) CreateTransfer(accountFrom int, accountTo int, amount uint, remarks string, trxDate time.Time, isTransfer bool) (model.Transaction, error) {
// 	var transactionID int
// 	if accountTo == accountFrom {
// 		return model.Transaction{}, errors.New("Cannot Transfer to the same account")
// 	}
// 	//check both account must be valid
// 	ret := []model.FinancialAccount{}
// 	err := ac.db.Model(&model.FinancialAccount{}).Find(&ret, "id IN (?)", []int{accountFrom, accountTo}).Error
// 	if err != nil {
// 		return model.Transaction{}, err
// 	}
// 	if len(ret) != 2 {
// 		return model.Transaction{}, errors.New("Invalid Account ID")
// 	}
// 	//check amount
// 	accountFromDetails := model.FinancialAccount{}
// 	accountToDetails := model.FinancialAccount{}
// 	for _, v := range ret {
// 		if v.ID == uint(accountFrom) {
// 			accountFromDetails = v
// 			if v.Balance < int(amount) {
// 				// return model.Transaction{}, errors.New("Source Account does not have enough balance")
// 			}
// 		} else {
// 			//accounto
// 			accountToDetails = v
// 		}
// 	}
// 	//	ac.db.Model((&model.Account))
// 	err = ac.db.Transaction(func(tx *gorm.DB) error {
// 		//create entry in transaction db
// 		trxCode := uuid.New().String()
// 		trxEntry := model.Transaction{IsTransfer: isTransfer, Remarks: remarks, TransactionCode: trxCode, TransactionDate: trxDate}
// 		if err := tx.Create(&trxEntry).Error; err != nil {
// 			return err
// 		}
// 		//create enty account cred
// 		if err := tx.Create(&model.Mutation{
// 			AccountID:     accountFrom,
// 			Amount:        -int(amount),
// 			TransactionID: int(trxEntry.ID),
// 			SoD:           accountToDetails.AccountName, //trf to account name
// 			// Remarks:       "TRF OUT: " + remarks,
// 		}).Error; err != nil {
// 			return err
// 		}
// 		//create enty account deb
// 		if err := tx.Create(&model.Mutation{
// 			AccountID:     accountTo,
// 			Amount:        int(amount),
// 			TransactionID: int(trxEntry.ID),
// 			SoD:           accountFromDetails.AccountName,
// 			// Remarks:       "TRF IN: " + remarks,
// 		}).Error; err != nil {
// 			return err
// 		}
// 		transactionID = int(trxEntry.ID)
// 		// return nil will commit the whole transaction
// 		return nil
// 	})
// 	if err != nil {
// 		return model.Transaction{}, err
// 	}

// 	return ac.GetTransactionDetailsbyID(transactionID)
// }
