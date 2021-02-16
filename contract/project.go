package contract

import (
	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
)

type ProjectStore interface {
	CreateProject(name string, accountID int, budget *uint, description *string) (model.Project, error)
	CreatePocket(projectID int, name string, accountID uint, limit *uint) (model.Project, error)
	//	GetProjectAccounts(projectID int) ( error)
	GetProjectDetails(projectID int) (model.Project, uint, []uint, error)
	CheckBudgetIDValidity(budgetID int, projectID int) (model.Budget, error)
	ListProject(req utils.CommonRequest) ([]model.Project, int, error)
	UpdateProject(prj model.Project) error
	DeleteProject(id int) error

	GetProjectAnalysis(id int) (map[string]interface{}, error)
}
