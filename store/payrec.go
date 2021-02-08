package store

import (
	"errors"
	"log"
	"time"

	"github.com/ahmadfarisfs/krab-core/contract"
	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
	"gorm.io/gorm"
)

type PayRecStore struct {
	db *gorm.DB
	ts contract.TransactionStore
	ps contract.ProjectStore
}

func NewPayRecStore(ts contract.TransactionStore, ps contract.ProjectStore, db *gorm.DB) *PayRecStore {
	return &PayRecStore{
		db: db,
		ts: ts,
		ps: ps,
	}
}

func (ps *PayRecStore) CreatePayRec(remarks string, amount int, projectID uint, pocketID *uint) (model.PayRec, error) {
	data := model.PayRec{
		Remarks:   remarks,
		Amount:    amount,
		ProjectID: projectID,
		PocketID:  pocketID,
	}
	err := ps.db.Model(&model.PayRec{}).Create(&data).Error
	if err != nil {
		return model.PayRec{}, err
	}
	return data, nil
}
func (ps *PayRecStore) Approve(id uint) (model.PayRec, error) {

	payRecDetails := model.PayRec{}
	err := ps.db.Model(&model.PayRec{}).Where("id = ? and transaction_code is null", id).First(&payRecDetails).Error
	if err != nil {
		return payRecDetails, err
	}
	//cannot add receivables to pocket account
	if payRecDetails.Amount > 0 && payRecDetails.PocketID != nil {
		return payRecDetails, errors.New("Cannot approve receivables to pocket account")
	}

	prjDet, prjAccountID, _, err := ps.ps.GetProjectDetails(int(payRecDetails.ProjectID))
	var accountID uint
	if payRecDetails.PocketID != nil {
		isValid := false
		for _, budget := range prjDet.Budgets {
			if *payRecDetails.PocketID == budget.ID {
				//valid
				isValid = true
				accountID = budget.AccountID
			}
		}
		if !isValid {
			return payRecDetails, errors.New("invalid pocket ID")
		}
	} else {
		accountID = prjAccountID
	}

	trx, err := ps.ts.CreateTransaction(int(accountID), payRecDetails.Amount, payRecDetails.Remarks, time.Now())
	if err != nil {
		return model.PayRec{}, err
	}

	err = ps.db.Model(&model.PayRec{}).Where("id = ?", id).Update("transaction_code", trx.TransactionCode).Error
	if err != nil {
		return model.PayRec{}, err
	}

	return payRecDetails, nil
}
func (ps *PayRecStore) Reject(id uint) (model.PayRec, error) {
	err := ps.db.Model(&model.PayRec{}).Where("id = ? and transaction_code is null", id).Delete(&model.PayRec{}).Error //.Update("transaction_code", trx.TransactionCode).Error
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
