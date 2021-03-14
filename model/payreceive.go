package model

import "time"

type PayRec struct {
	BaseModel
	Remarks         string
	Approved        bool
	TransactionCode *string //not null if approved
	ProjectID       uint
	Project         Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	// PocketID        *uint   //can be adjusted when approval
	// Pocket          *Budget `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Amount int

	SourceProjectAccountID *int
	TargetProjectAccountID *int

	SourceBankAccountID *int
	TargetBankAccountID *int

	// SoD             string
	Email           string
	Meta            string
	Notes           string
	TransactionDate time.Time
	TargetUserID    *int //only for reimbursement

}
