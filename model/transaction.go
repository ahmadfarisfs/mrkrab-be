package model

import "time"

type Transaction struct {
	BaseModel
	TransactionCode string
	Remarks         string
	IsTransfer      bool
	Mutation        []Mutation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TransactionDate time.Time
	// SoD             string
}
