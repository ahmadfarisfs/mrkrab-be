package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/ahmadfarisfs/mrkrab-be/domain"
	//	"gorm.io/gorm"
)

type transactionUseCase struct {
	userRepo        domain.UserRepository
	timeout         time.Duration
	transactionRepo domain.TransactionRepository
}

func NewTransactionUseCase(p domain.TransactionRepository, u domain.UserRepository, timeout time.Duration) domain.TransactionUsecase {
	return &transactionUseCase{
		transactionRepo: p,
		userRepo:        u,
		timeout:         timeout,
	}
}

func (p *transactionUseCase) Fetch(ctx context.Context, limitPerPage int64, page int64, filter *domain.Transaction) (res []domain.Transaction, totalRecord int, totalPage int, err error) {
	if limitPerPage == 0 {
		limitPerPage = 10
	}
	if page == 0 {
		page = 1
	}

	return nil, 0, 0, nil
}

func (p *transactionUseCase) GetByID(ctx context.Context, id int64) (domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	return p.transactionRepo.GetByID(ctx, id)
}

func (p *transactionUseCase) Update(ctx context.Context, ar *domain.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	return p.transactionRepo.Update(ctx, ar)
}

func (p *transactionUseCase) Add(ctx context.Context, a *domain.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if a.SoFType == "user" {
		if a.SoFUserID == nil {
			return errors.New("Source of Fund userID must be set if SoF type is user")
		}
	} else if a.SoFType == "other" {
		if a.SoFAccount == "" {
			return errors.New("Source of Fund account must be set if SoF type is other")
		}
	} else {
		return errors.New("Invalid SoF type")
	}

	if a.Paid {
		timeNow := time.Now()
		a.PaidOn = &timeNow
	}

	/*
		if a.Type == "credit" {

		} else if a.Type == "debit" {
			//
		} else {
			return errors.New("Invalid Transaction Type")
		}*/

	return p.transactionRepo.Store(ctx, a)
}

func (p *transactionUseCase) Delete(ctx context.Context, id int64) error {
	return nil
}
