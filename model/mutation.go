package model

import (
	"time"
)

type Mutation struct {
	BaseModel
	AccountID     int
	Account       Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TransactionID int
	Transaction   Transaction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Amount        int         //deltas
	SoD           string
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
	// Mutation
	// Project Project `gorm:"-"`
	// Budget  *Budget `gorm:"-"`
}
