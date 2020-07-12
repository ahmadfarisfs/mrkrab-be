package domain

import (
	"context"
	//"github.com/jinzhu/gorm"
)

// Transaction ...
type Transaction struct {
	//	gorm.Model
	BaseModel
	Project   Project `gorm:"foreignkey:ProjectID" json:"-"`
	ProjectID int     `gorm:"not null" json:"project_id" validate:"required"`

	Creator   User `gorm:"foreignkey:CreatorID" json:"-"`
	CreatorID int  `gorm:"not null" json:"creator_id" validate:"required"`

	Category   Category `gorm:"foreignkey:CategoryID" json:"-"`
	CategoryID int      `gorm:"not null"`

	Description string `gorm:"not null" json:"description" validate:"required"`
	Amount      int    `gorm:"not null" json:"amount" validate:"required"`
	Type        string `gorm:"not null;type:enum('credit','debit')" json:"type" validate:"required"`
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
