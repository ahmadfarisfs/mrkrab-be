package model

type Project struct {
	BaseModel
	Name        string `gorm:"unique"`
	AccountID   int
	Account     Account `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Amount      *uint
	IsOpen      bool
	Description string
	// IsPooling   bool
	Budgets []Budget `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` //`gorm:"many2many:project_budgets;"`
}

type Budget struct {
	BaseModel
	// Name             string
	ProjectID        uint
	Project          Project `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ExpenseAccountID uint
	ExpenseAccount   Account `json:"ExpenseAccount" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Limit            uint
}
