package model

import (
	"time"
)

type FinancialAccountMutation struct {
	BaseModel
	AccountID             int
	Account               FinancialAccount `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	TransactionID         int
	Transaction           Transaction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	TransactionCode       string
	Amount                int
	ProjectID             int
	Project               Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	IsPaid                bool
	BankAccountMutationID *uint
	BankAccountMutation   *BankAccountMutation
}

type BankAccountMutation struct {
	BaseModel
	BankAccountID   int
	Bank            BankAccount `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	TransactionID   int
	Transaction     Transaction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
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
