package model

import "time"

type Transaction struct {
	ID              uint `gorm:"primarykey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	TransactionCode string
	Remarks         string
	Mutation        []Mutation
}
