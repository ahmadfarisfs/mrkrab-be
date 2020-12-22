package model

import "time"

type Account struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	AccountName  string `gorm:"unique"`
	ParentID     *uint
	Parent       *Account
	Balance      int
	TotalIncome  int
	TotalExpense int
}
