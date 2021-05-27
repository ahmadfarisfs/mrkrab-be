package model

import (
	"time"
)

type FinancialAccountMutation struct {
	BaseModel
	AccountID       int
	Account         FinancialAccount `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TransactionID   int
	Transaction     Transaction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TransactionCode string
	Amount          int
	ProjectID       int
	Project         Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	IsPaid          bool
}

type BankAccountMutation struct {
	BaseModel
	BankAccountID   int
	Bank            BankAccount `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TransactionID   int
	Transaction     Transaction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TransactionCode string
	Amount          int
	IsPaid          bool
}

type MutationExtended struct {
	ID                 uint
	PocketID           *uint
	ProjectID          uint
	IsOpen             bool
	CreatedAt          time.Time
	Amount             int
	Remarks            string
	TransactionCode    string
	ProjectDescription string
	ProjectName        string
	PocketName         *string
	PocketLimit        *int
	TransactionDate    time.Time
	SoD                string
}
