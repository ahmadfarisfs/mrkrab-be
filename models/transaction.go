package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type TransactionCategory struct {
	gorm.Model
	Name        string
	Description string
	CreatedBy   User
	IsDeleted   bool
}

type Transaction struct {
	gorm.Model
	TransactionDate time.Time
	Title           string
	Description     string
	Amount          int
	CreatedBy       User
	Project         Project
	Category        TransactionCategory
	Type            TransactionType
	ImageURL        string
	IsDeleted       bool
}

type Project struct {
	gorm.Model
	Title       string
	Location    string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	Budget      int
	PIC         User
	CreatedBy   User
	IsDeleted   bool
	IsClosed    bool
	TeamMember  []User `gorm:"many2many:project_member;"`
}
