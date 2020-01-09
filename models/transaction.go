package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type TransactionCategory struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string `gorm:"type:varchar(255)"`
	CreatedBy   uint   `gorm:"size:10;not null" migrator:"foreignkey:created_by;foreignkeyof:users(id);ondelete:RESTRICT;onupdate:CASCADE"`
	//IsDeleted   bool
}

type Transaction struct {
	gorm.Model
	TransactionDate time.Time
	Title           string              `gorm:"not null"`
	Description     string              `gorm:"type:varchar(max)"`
	Amount          int                 `gorm:"not null"`
	CreatedBy       uint                `gorm:"size:10;not null" migrator:"foreignkey:created_by;foreignkeyof:users(id);ondelete:RESTRICT;onupdate:CASCADE"`
	Project         uint                `gorm:"size:10;not null" migrator:"foreignkey:project;foreignkeyof:projects(id);ondelete:RESTRICT;onupdate:CASCADE"`
	ProjectName     string              `gorm:"type:varchar(255)"`
	Category        TransactionCategory `gorm:"type:integer"`
	CategoryName    string              `gorm:"type:varchar(255)"`
	Type            TransactionType     `gorm:"type:integer"`
	ImageURL        string              `gorm:"type:varchar(max)"`
	//	IsDeleted       bool                `gorm:"default:false;not null" json:"-"`
}

type Project struct {
	gorm.Model
	Title         string `gorm:"not null;unique"`
	Location      string `gorm:"not null"`
	Description   string `gorm:"not null"`
	StartDate     time.Time
	EndDate       time.Time
	Budget        int    `gorm:"not null"`
	PIC           uint   `gorm:"size:10;not null" migrator:"foreignkey:pic;foreignkeyof:users(id);ondelete:RESTRICT;onupdate:CASCADE"`
	CreatedBy     uint   `gorm:"size:10;not null" migrator:"foreignkey:created_by;foreignkeyof:users(id);ondelete:RESTRICT;onupdate:CASCADE"`
	CreatedByName string `gorm:"type:varchar(255);not null;unique"`
	PICName       string `gorm:"type:varchar(255);not null;unique"`
	//	IsDeleted     bool   `gorm:"default:false;not null"`
	IsClosed   bool   `gorm:"default:false;not null"`
	TeamMember []User `gorm:"many2many:project_member;"`
}
