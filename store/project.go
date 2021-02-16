package store

import (
	"log"
	"math"

	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
	"gorm.io/gorm"
)

type ProjectStore struct {
	db *gorm.DB
}

func NewProjectStore(db *gorm.DB) *ProjectStore {
	return &ProjectStore{
		db: db,
	}
}

func (ps *ProjectStore) CreateProject(name string, accountID int, budget *uint, description *string) (model.Project, error) {

	ret := model.Project{Name: name, AccountID: accountID, Amount: budget, Description: description}
	err := ps.db.Model(&model.Project{}).Create(&ret).Error
	if err != nil {
		return model.Project{}, err
	}

	return ret, err
}

func (ps *ProjectStore) CreatePocket(projectID int, name string, accountID uint, limit *uint) (model.Project, error) {
	mdl := model.Budget{
		Name:      name,
		ProjectID: uint(projectID),
		Limit:     limit,
		AccountID: accountID,
	}

	//err := ps.db.Model(&model.Project{}).Association("Budgets").Append(&mdl)//idk why this wont work
	err := ps.db.Model(&model.Budget{}).Create(&mdl).Error
	if err != nil {
		return model.Project{}, err
	}
	prj, _, _, err := ps.GetProjectDetails(projectID)
	if err != nil {
		return model.Project{}, err
	}
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

func (ps *ProjectStore) CheckBudgetIDValidity(budgetID int, projectID int) (model.Budget, error) {
	ret := model.Budget{}
	err := ps.db.Model(&model.Budget{}).First(&ret, "id = ? and project_id = ?", budgetID, projectID).Error
	if err != nil {
		return model.Budget{}, err
	}
	return ret, err
}

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

func (ps *ProjectStore) DeleteProject(id int) error {
	return ps.db.Where("id = ?", id).Delete(&model.Project{}).Error
}

func (ps *ProjectStore) UpdateProject(prj model.Project) error {
	log.Println(prj)
	editPayload := map[string]interface{}{"is_open": prj.IsOpen}
	if prj.Description != nil {
		editPayload["description"] = prj.Description
	}
	return ps.db.Model(&model.Project{}).Where("id = ?", prj.ID).Updates(editPayload).Error
}
func (ps *ProjectStore) GetProjectAnalysis(id int) (map[string]interface{}, error) {
	prj, accountID, _, err := ps.GetProjectDetails(id)
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
	}

	return retDat, nil
}
