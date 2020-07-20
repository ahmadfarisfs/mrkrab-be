package domain

import (
	"context"
	"time"
	//"github.com/jinzhu/gorm"
)

// Transaction ...
type Transaction struct {
	//	gorm.Model
	BaseModel
	Project   Project `gorm:"foreignkey:ProjectID" json:"-" validate:"-"`
	ProjectID int     `gorm:"not null" json:"project_id" validate:"required"`

	Creator   User `gorm:"foreignkey:CreatorID" json:"-" validate:"-" `
	CreatorID int  `gorm:"not null" json:"creator_id" validate:"required"`

	Category   Category `gorm:"foreignkey:CategoryID" json:"-" validate:"-"`
	CategoryID int      `gorm:"not null" validate:"required" json:"category_id"`

	Paid    bool   `gorm:"not null" json:"paid"`
	SoFType string `gorm:"null;type:enum('project','user','other')" json:"sof_type" validator:"oneof=project user other"` //this field must be set if paid is false

	SoFUser      *User    `gorm:"null;foreignkey:SoFUserID" json:"-"  validate:"-"`
	SoFUserID    *int     `gorm:"null" json:"sof_user_id"` //this field must be set if sof_type = user
	SoFProject   *Project `gorm:"null;foreignkey:SoFProjectID" json:"-"  validate:"-"`
	SoFProjectID *int     `gorm:"null" json:"sof_project_id"` //this field must be set if sof_type = user

	SoFAccount *string    `gorm:"null" json:"sof_account"` //this field must be set if sof_type = other
	PaidOn     *time.Time `gorm:"null" json:"paid_on"`     //this field should be set when paid become true

	Description string `gorm:"not null" json:"description" validate:"required"`
	Amount      uint   `gorm:"not null" json:"amount" validate:"required,min=0"`
	Type        string `gorm:"not null;type:enum('credit','debit')" json:"type" validate:"required,oneof=credit debit"`

	Approved bool `gorm:"not null;default:false" json:"approved"`
}

// TransactionUsecase represent the Transaction's usecases (business process)
type TransactionUsecase interface {
	Fetch(ctx context.Context, limitPerPage int64, page int64, filter *Transaction) (res []Transaction, totalRecord int, totalPage int, err error)
	GetByID(ctx context.Context, id int64) (Transaction, error)
	Update(ctx context.Context, Transaction *Transaction) error
	Delete(ctx context.Context, id int64) error
	Add(context.Context, *Transaction) error
}

// TransactionRepository represent the Transactions's repository contract -> implemented in db conn
type TransactionRepository interface {
	Fetch(ctx context.Context, limitPerPage int64, page int64, filter *Transaction) (res []Transaction, totalRecord int, totalPage int, err error)
	GetByID(ctx context.Context, id int64) (Transaction, error)
	Update(ctx context.Context, ar *Transaction) error
	Store(ctx context.Context, a *Transaction) error
	Delete(ctx context.Context, id int64) error
}
