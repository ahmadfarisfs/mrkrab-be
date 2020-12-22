package model

import "time"

type Project struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string `gorm:"unique"`
	AccountID   int
	Account     Account `json:"-"`
	Amount      *uint
	IsOpen      bool
	Description *string
	Budgets     []Budget //`gorm:"many2many:project_budgets;"`
}

type Budget struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	ProjectID uint
	Project   Project `json:"-"`
	AccountID uint
	Account   Account `json:"-"`
	Limit     *uint
}
