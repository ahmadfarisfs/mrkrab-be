package store

import (
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
	"gorm.io/gorm"
)

type ProjectStore struct {
	db *gorm.DB
	// as *AccountStore
}

func NewProjectStore(db *gorm.DB) *ProjectStore {
	return &ProjectStore{
		db: db,
		// as: account,
	}
}

func (ps *ProjectStore) CreateProject(name string, budget *uint, description string, budgetIDs []uint, limits []uint) (model.Project, error) {
	var prjID uint
	err := ps.db.Transaction(func(tx *gorm.DB) error {
		prjAccountName := "PROJECT-" + strings.ToUpper(name) + "-" + strconv.Itoa(int(time.Now().Unix()))
		ret := model.Account{
			AccountName: prjAccountName,
			AccountType: "PROJECT",
		}
		if err := tx.Model(&model.Account{}).Create(&ret).Error; err != nil {
			return err
		}

		rets := model.Project{Name: name, AccountID: int(ret.ID), Amount: budget, Description: description}
		if err := tx.Model(&model.Project{}).Create(&rets).Error; err != nil {
			return err
		}

		for idx, v := range budgetIDs {
			mdl := model.Budget{
				ProjectID:        ret.ID,
				ExpenseAccountID: v,
				Limit:            limits[idx],
			}
			if err := tx.Model(&model.Budget{}).Create(&mdl).Error; err != nil {
				return err
			}
		}
		prjID = ret.ID
		return nil
	})
	if err != nil {
		return model.Project{}, err
	}

	prj, _, _, err := ps.GetProjectDetails(int(prjID))
	return prj, err
}

//  func(ps *ProjectStore) DeletePocket
func (ps *ProjectStore) CreatePocket(projectID int, accountID uint, limit uint) (model.Project, error) {
	mdl := model.Budget{
		ProjectID:        uint(projectID),
		Limit:            limit,
		ExpenseAccountID: accountID,
	}

	err := ps.db.Model(&model.Budget{}).Create(&mdl).Error
	if err != nil {
		return model.Project{}, err
	}
	prj, _, _, err := ps.GetProjectDetails(projectID)
	return prj, err
}

func (ps *ProjectStore) GetProjectDetails(projectID int) (model.Project, uint, []uint, error) {
	ret := model.Project{}
	err := ps.db.Preload("Budgets").Preload("Account").First(&ret, "id = ?", projectID).Error
	budgetIDs := []uint{}
	for _, v := range ret.Budgets {
		budgetIDs = append(budgetIDs, v.ID)
	}
	return ret, uint(ret.AccountID), budgetIDs, err
}

// func (ps *ProjectStore) CheckBudgetIDValidity(budgetID int, projectID int) (model.Budget, error) {
// 	ret := model.Budget{}
// 	err := ps.db.Model(&model.Budget{}).First(&ret, "id = ? and project_id = ?", budgetID, projectID).Error
// 	if err != nil {
// 		return model.Budget{}, err
// 	}
// 	return ret, err
// }
// func (ps *ProjectStore) CreateProjectIncome()

func (ps *ProjectStore) ListProject(req utils.CommonRequest) ([]model.Project, int, error) {
	ret := []model.Project{}
	var count int64
	//query builder
	initQuery := ps.db

	err := initQuery.Model(&model.Project{}).Count(&count).Error
	if err != nil {
		return ret, int(count), err
	}
	log.Println(req)
	//actually fetch data with limit and offset
	quer := utils.AppendCommonRequest(initQuery, req)
	err = quer.Preload("Budgets").Find(&ret).Error
	return ret, int(count), err
}
func (ps *ProjectStore) UpdatePocket(id int, limit int) error {
	return ps.db.Model(&model.Budget{}).Where("id = ?", id).Update("limit", limit).Error
}
func (ps *ProjectStore) DeletePocket(id int) error {
	return ps.db.Where("id = ?", id).Delete(&model.Budget{}).Error
}
func (ps *ProjectStore) DeleteProject(id int) error {
	return ps.db.Where("id = ?", id).Delete(&model.Project{}).Error
}

func (ps *ProjectStore) UpdateProject(prj model.Project) error {
	log.Println(prj)
	editPayload := map[string]interface{}{"is_open": prj.IsOpen}
	// if prj.Description != nil {
	editPayload["description"] = prj.Description
	// }
	return ps.db.Model(&model.Project{}).Where("id = ?", prj.ID).Updates(editPayload).Error
}
func (ps *ProjectStore) GetProjectAnalysis(id int) (map[string]interface{}, error) {
	prj, accountID, _, err := ps.GetProjectDetails(id)
	if err != nil {
		return nil, err
	}
	budgets := []model.Budget{}
	err = ps.db.Model(&model.Budget{}).Preload("Account").Where("project_id = ?", id).Find(&budgets).Error
	if err != nil {
		return nil, err
	}
	accountDet := model.Account{}
	err = ps.db.Model(&model.Account{}).Where("id = ?", accountID).First(&accountDet).Error
	if err != nil {
		return nil, err
	}
	subAccountsDet := []model.Account{}
	err = ps.db.Model(&model.Account{}).Where("parent_id is not null and parent_id = ?", accountDet.ID).
		Find(&subAccountsDet).Error
	if err != nil {
		return nil, err
	}
	pocketTotalExpense := 0
	pocketTotalIncome := 0
	for _, v := range subAccountsDet {
		pocketTotalExpense += v.TotalExpense
		pocketTotalIncome += v.TotalIncome
	}
	payRec := []model.PayRec{}
	err = ps.db.Model(&model.PayRec{}).Where("project_id = ? ", id).Find(&payRec).Error
	if err != nil {
		return nil, err
	}
	totalPayables := 0
	totalReceivables := 0
	for _, v := range payRec {
		if v.Amount < 0 {
			totalPayables += int(math.Abs(float64(v.Amount)))
		} else {
			totalReceivables += int(math.Abs(float64(v.Amount)))
		}
	}

	//seharusnya tidak ada transfer dari project ke pocket
	retDat := map[string]interface{}{
		// "ProjectTotalExpense": accountDet.TotalExpense,
		// "ProjectTotalIncome":  accountDet.TotalIncome,
		// "General":             accountDet,
		"Project":          prj,
		"ProjectAccount":   accountDet,
		"Pockets":          subAccountsDet,
		"TotalPayables":    totalPayables,
		"TotalReceivables": totalReceivables,
		"Budgets":          budgets,
	}

	return retDat, nil
}
