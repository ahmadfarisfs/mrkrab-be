package model

import "time"

type Transaction struct {
	BaseModel
	TransactionCode string
	Remarks         string
	Notes           string
	Meta            string
	FromAccountID   int64      `json:"from_account_id"`
	ToAccountID     int64      `json:"to_account_id"`
	TransactionType string     //enum: BANK or PROJECT or BOTH
	Mutation        []Mutation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TransactionDate time.Time
}
