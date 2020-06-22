package domain

import (
	"context"
	//"github.com/jinzhu/gorm"
)

// Project ...
type Project struct {
	BaseModel
	Name        string          `gorm:"not null" json:"name" validate:"required"`
	ProjectType string          `gorm:"not null;type:enum('onetime','business_unit')" json:"project_type"`
	Status      string          `gorm:"not null;type:enum('offering','ongoing','close')" json:"status"`
	PICID       int             `json:"pic_id"`
	PIC         User            `json:"pic_details" gorm:"foreignkey:PICID"`
	Budgets     []ProjectBudget `json:"-" gorm:"foreignkey:ProjectID"`
	Members     []User          `json:"-" gorm:"many2many:user_projects;"`
}

// ProjectBudget define budget foreach cateory on a project
type ProjectBudget struct {
	BaseModel
	ProjectID int     `gorm:"not null" json:"project_id"`
	Project   Project `gorm:"foreignkey:ProjectID" json:"-"`

	CategoryID int      `gorm:"not null" json:"category_id"`
	Category   Category `gorm:"foreignkey:CategoryID" json:"-"`
	Amount     int      `gorm:"not null" json:"amount"`
}

//ProjectStatus ...
type ProjectStatus int

const (
	Offering ProjectStatus = iota
	Ongoing
	Closed
)

func (d ProjectStatus) String() string {
	return [...]string{"offering", "ongoing", "closed"}[d]
}

//ProjectMemberRole ...
type ProjectMemberRole int

const (
	PIC ProjectMemberRole = iota
	Member
)

func (d ProjectMemberRole) String() string {
	return [...]string{"pic", "member"}[d]
}

// ProjectUsecase represent the Project's usecases (business process)
type ProjectUsecase interface {
	Fetch(ctx context.Context, limitPerPage int64, page int64, filter map[string]string) ([]Project, error)
	GetByID(ctx context.Context, id int64) (Project, error)
	Update(ctx context.Context, Project *Project) error
	Delete(ctx context.Context, id int64) error
	Add(context.Context, *Project) error

	GetProjectMember(ctx context.Context, projectID int64) ([]User, error)
	AssignPIC(ctx context.Context, projectID int64, userID int64) error
	AssignMember(ctx context.Context, projectID int64, userID int64) error
	RemoveMember(ctx context.Context, projectID int64, userID int64) error

	SetStatus(ctx context.Context, projectID int64, status ProjectStatus) error
}

// ProjectRepository represent the Projects's repository contract -> implemented in db conn
type ProjectRepository interface {
	Fetch(ctx context.Context, limitPerPage int64, page int64) (res []Project, totalRecord int, totalPage int, err error)
	GetByID(ctx context.Context, id int64) (Project, error)
	Update(ctx context.Context, ar *Project) error
	Store(ctx context.Context, a *Project) error
	Delete(ctx context.Context, id int64) error
	//IsPICAssigned(ctx context.Context, projectID int64) (bool, error)
	GetProjectsByUser(ctx context.Context, userID int64) (map[ProjectMemberRole][]Project, error)
	GetProjectMember(ctx context.Context, projectID int64) (map[ProjectMemberRole][]User, error)
	AddMember(ctx context.Context, projectID int64, userID int64, role ProjectMemberRole) error
	RemoveMember(ctx context.Context, projectID int64, userID int64) error
}
