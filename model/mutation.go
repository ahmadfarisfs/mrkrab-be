package model

import "time"

type Mutation struct {
	ID            uint      `gorm:"primarykey"` // json:"id"`
	CreatedAt     time.Time //`json:"created_at"`
	UpdatedAt     time.Time
	AccountID     int
	Account       Account
	TransactionID int
	Transaction   Transaction
	Amount        int //deltas
}
