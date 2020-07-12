package mysql

import (
	"context"

	"github.com/ahmadfarisfs/mrkrab-be/domain"

	"gorm.io/gorm"
)

type mysqlTransactionRepo struct {
	DB *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) domain.TransactionRepository {
	return &mysqlTransactionRepo{
		DB: db,
	}
}

func (m *mysqlTransactionRepo) Fetch(ctx context.Context, limitPerPage int64, page int64, filter *domain.Transaction) (res []domain.Transaction, totalRecord int, totalPage int, err error) {
	return nil, 0, 0, nil
}

func (m *mysqlTransactionRepo) GetByID(ctx context.Context, id int64) (domain.Transaction, error) {
	return domain.Transaction{}, nil
}

func (m *mysqlTransactionRepo) Update(ctx context.Context, ar *domain.Transaction) error {
	return nil
}
func (m *mysqlTransactionRepo) Store(ctx context.Context, a *domain.Transaction) error {
	return m.DB.Create(a).Error
}
func (m *mysqlTransactionRepo) Delete(ctx context.Context, id int64) error {
	return nil
}
