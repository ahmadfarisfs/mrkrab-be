package contract

import (
	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
)

type ProjectStore interface {
	CreateProject(name string, amount int, description string, status string) (model.Project, error)
	GetProjectDetails(projectID int) (model.Project, error)
	ListProject(req utils.CommonRequest) ([]model.Project, int, error)
	UpdateProject(id int, name string, amount int, description string, status string) error
	DeleteProject(id int) error

	GetProjectAnalysis(id int) (map[string]interface{}, error)
}

// p ProjectStore
