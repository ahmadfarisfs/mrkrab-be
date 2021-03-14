package store

import (
	"errors"
	"log"
	"time"

	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
	"gorm.io/gorm"
)

type PayRecStore struct {
	db *gorm.DB
	ts *TransactionStore
	ps *ProjectStore
}

func NewPayRecStore(ts *TransactionStore, ps *ProjectStore, db *gorm.DB) *PayRecStore {
	return &PayRecStore{
		db: db,
		ts: ts,
		ps: ps,
	}
}

// func (ps *PayRecStore) CreatePayables()

func (ps *PayRecStore) CreatePayRec(remarks string, notes string, meta string, isReimburse bool, amount int, projectID int, sourceProjectAccountID *int, destProjectAccountID *int, sourceBankAccountID *int, destBankAccountID *int) (model.PayRec, error) {
	if amount < 0 {
		//payables: must be defined: destProjectAccount (expense accont)
		if destBankAccountID == nil || sourceProjectAccountID == nil {
			return model.PayRec{}, errors.New("Dest Project Account must be defined for payables")
		}
		if isReimburse && destBankAccountID == nil {
			//reimburse dest bank account must be user
			return model.PayRec{}, errors.New("Dest Bank account must be defined")
		}
	} else {
		//receivables
		if destProjectAccountID == nil {
			return model.PayRec{}, errors.New("Dest Project Account must be defined for receivables")
		}
	}
	data := model.PayRec{
		ProjectID:              uint(projectID),
		Remarks:                remarks,
		Notes:                  notes,
		Meta:                   meta,
		Amount:                 amount,
		SourceBankAccountID:    sourceBankAccountID,
		TargetBankAccountID:    destBankAccountID,
		SourceProjectAccountID: sourceProjectAccountID, //can be project account (payables) or income account (receivables)
		TargetProjectAccountID: destProjectAccountID,   //can be project account (receivables) or expense account (payables)
	}
	err := ps.db.Model(&model.PayRec{}).Create(&data).Error
	if err != nil {
		return model.PayRec{}, err
	}
	return data, nil
}
func (ps *PayRecStore) Approve(id uint, sourceProjectAccountID int, destProjectAccountID int, sourceBankAccountID int, destBankAccountID int, trxDate *time.Time) (model.PayRec, error) {
	payRecDetails := model.PayRec{}
	err := ps.db.Model(&model.PayRec{}).Where("id = ? and transaction_code is null", id).First(&payRecDetails).Error
	if err != nil {
		return payRecDetails, err
	}
	var overDate time.Time
	if trxDate == nil {
		overDate = time.Now()
	} else {
		overDate = *trxDate
	}
	// if payRecDetails.Amount < 0 {
	//payables
	ret, err := ps.ts.CreateProjectTransaction(uint(sourceProjectAccountID), uint(destProjectAccountID), uint(sourceBankAccountID), uint(destBankAccountID), uint(payRecDetails.Amount),
		payRecDetails.Remarks,
		payRecDetails.Meta,
		payRecDetails.Notes,
		overDate)
	if err != nil {
		return model.PayRec{}, err
	}
	// } else {
	// 	//receivables
	// }
	//cannot add receivables to pocket account
	// if payRecDetails.Amount > 0 && payRecDetails.PocketID != nil {
	// 	return payRecDetails, errors.New("Cannot approve receivables to pocket account")
	// }

	// prjDet, prjAccountID, budgetAccountIDs, err := ps.ps.GetProjectDetails(int(payRecDetails.ProjectID))
	// // var accountID uint
	// // if payRecDetails.PocketID != nil {
	// // 	isValid := false
	// // 	for _, budget := range prjDet.Budgets {
	// // 		if *payRecDetails.PocketID == budget.ID {
	// // 			//valid
	// // 			isValid = true
	// // 			accountID = budget.AccountID
	// // 		}
	// // 	}
	// // 	if !isValid {
	// // 		return payRecDetails, errors.New("invalid pocket ID")
	// // 	}
	// // } else {
	// // 	// accountID = prjAccountID
	// // }

	// trx := model.Transaction{}

	// if payRecDetails.TargetUserAccountID != nil {
	// 	//reimburse approved, jadi pengeluaran di pockets
	// 	// ps.ts.CreateProjectTransaction(int(prjAccountID))
	// 	// trx, err = ps.ts.CreateTransfer(int(prjAccountID), int(accountID), uint(math.Abs(float64(payRecDetails.Amount))), payRecDetails.Remarks, payRecDetails.TransactionDate, false)
	// 	if err != nil {
	// 		return model.PayRec{}, err
	// 	}
	// } else {
	// 	//normal
	// 	if payRecDetails.Amount > 0 {
	// 		//income can only come from REVENUE ACCOUNT with prefix name: ACCOUNT-REVENUE and id 0
	// 		// to projects
	// 		//find revenue account
	// 		trx, err = ps.ts.CreateTransfer(0, int(prjAccountID), uint(math.Abs(float64(payRecDetails.Amount))), payRecDetails.Remarks, payRecDetails.TransactionDate, false)
	// 		if err != nil {
	// 			return model.PayRec{}, err
	// 		}
	// 		// return trx,nil
	// 	} else {
	// 		//expense
	// 		var accountID int
	// 		if payRecDetails.PocketID != nil {
	// 			//check that account should be under projects
	// 			isValid := false
	// 			for _, v := range budgetAccountIDs {
	// 				if v == *payRecDetails.PocketID {
	// 					//good
	// 					isValid = true
	// 				}
	// 			}
	// 			if !isValid {
	// 				return model.PayRec{}, errors.New("invalid budget id")
	// 			}
	// 			//harus transalate dari budgetID ke accountID
	// 			isValid = false
	// 			for _, budget := range prjDet.Budgets {
	// 				if budget.ID == uint(*payRecDetails.PocketID) {
	// 					accountID = int(budget.AccountID)
	// 					isValid = true
	// 				}
	// 			}
	// 			if !isValid {
	// 				return model.PayRec{}, errors.New("Invalid budget account ID")
	// 			}
	// 		} else {
	// 			return model.PayRec{}, errors.New("Expense account must be defined")
	// 			// accountID = int(projAccountID)
	// 		}

	// 		trx, err = ps.ts.CreateTransfer(int(prjAccountID), accountID, uint(math.Abs(float64(payRecDetails.Amount))), payRecDetails.Remarks, payRecDetails.TransactionDate, false)
	// 		// trx, err := h.transactionStore.CreateTransaction(accountID, req.Amount, req.Remarks, req.SoD, req.TransactionDate)
	// 		if err != nil {
	// 			return model.PayRec{}, err
	// 		}

	// 		// return c.JSON(http.StatusOK, trx)
	// 	}
	// 	// _,err=ps.ts.CreateTransfer(int(prjAccountID),int(accountID),uint(math.Abs(float64(payRecDetails.Amount))),payRecDetails.Remarks,payRecDetails.TransactionDate,false)
	// 	// if err != nil {
	// 	// 	return model.PayRec{}, err
	// 	// }
	// }

	// trx, err := ps.ts.CreateTransaction(int(accountID),
	// 	payRecDetails.Amount,
	// 	payRecDetails.Remarks,
	// 	payRecDetails.SoD,
	// 	time.Now())
	// if err != nil {
	// 	return model.PayRec{}, err
	// }

	err = ps.db.Model(&model.PayRec{}).Where("id = ?", id).Update("transaction_code", ret.TransactionCode).Error
	if err != nil {
		return model.PayRec{}, err
	}

	return payRecDetails, nil
}
func (ps *PayRecStore) Reject(id uint) (model.PayRec, error) {
	// err := ps.db.Model(&model.PayRec{}).Where("id = ? and transaction_code is null", id).Delete(&model.PayRec{}).Error //.Update("transaction_code", trx.TransactionCode).Error
	err := ps.db.Model(&model.PayRec{}).Where("id = ? and transaction_code is null", id).Update("approved", true).Error

	if err != nil {
		return model.PayRec{}, err
	}

	return model.PayRec{}, nil
}
func (ps *PayRecStore) ListPayRec(req utils.CommonRequest) ([]model.PayRec, int, error) {
	ret := []model.PayRec{}
	//query builder
	var count int64
	initQuery := ps.db

	err := initQuery.Model(&model.PayRec{}).Count(&count).Error
	if err != nil {
		return ret, int(count), err
	}
	log.Println(req)
	//actually fetch data with limit and offset
	quer := utils.AppendCommonRequest(initQuery, req)
	err = quer.Preload("Project").Preload("Pocket").Find(&ret).Error
	return ret, int(count), err
}
