package model

type FinancialAccount struct {
	BaseModel
	AccountName  string `gorm:"unique"`
	ParentID     *uint
	Parent       *FinancialAccount `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Balance      int
	TotalIncome  int
	TotalExpense int
	AccountType  string //expense or income
	Description  string
}

type BankAccount struct {
	BaseModel
	BankName        string
	HolderName      string
	BankNumber      string
	BankAccountType string //external or internal
	Description     string
	Balance         int
}
