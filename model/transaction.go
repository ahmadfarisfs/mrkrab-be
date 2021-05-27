package model

import "time"

type Transaction struct {
	BaseModel
	TransactionCode string `gorm:"unique"`
	Remarks         string
	TransactionType string //income, expense, bank transfer or account transfer
	TransactionTime time.Time
	BankMutation    []BankAccountMutation
	AccountMutation []FinancialAccountMutation
	IsPaid          bool
	Amount          int
	//redundant field for frontend clarity
	// SourceBankID *int
	// SourceBank   *BankAccount

	// SourceProjectID *int
	// SourceProject   *Project

	// DestinationBankID *int
	// DestinationBank   *BankAccount

	// DestinationProjectID *int
	// DestinationProject   *Project
}
