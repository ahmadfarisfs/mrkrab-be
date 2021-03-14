package model

type BankAccount struct {
	BaseModel
	BankName   string
	BankNumber string
	AccountID  int
	Account    Account `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
