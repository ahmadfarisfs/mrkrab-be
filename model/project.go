package model

type Project struct {
	BaseModel
	Name        string
	AccountID   int
	Account     Account `json:"-"`
	Amount      *uint
	IsOpen      bool
	Description *string
	Budgets     []Budget //`gorm:"many2many:project_budgets;"`
}

type Budget struct {
	BaseModel
	Name      string
	ProjectID uint
	Project   Project `json:"-"`
	AccountID uint
	Account   Account `json:"-"`
	Limit     *uint
}
