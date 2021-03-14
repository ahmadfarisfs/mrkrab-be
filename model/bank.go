package model

type BankAccount struct {
	BaseModel
	Internal       bool
	BankName       string
	BankNumber     string
	BankHolderName string
	AccountID      int
	Account        Account `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ExpenseAccount struct {
	BaseModel
	Name      string
	AccountID int
	Account   Account `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type IncomeAccount struct {
	BaseModel
	Name      string
	AccountID int
	Account   Account `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
