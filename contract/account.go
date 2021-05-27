package contract

import (
	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
)

type FinancialAccountStore interface {
	ListAccount(req utils.CommonRequest) ([]model.FinancialAccount, int, error)
	CreateAccount(name string, description string, accountType string, parentID *uint) (model.FinancialAccount, error)
	GetAccountDetails(id int) (model.FinancialAccount, error)
	UpdateAccount(id int, name string, description string, accountType string) (model.FinancialAccount, error)
	DeleteAccount(id int) error
}

type BankAccountStore interface {
	ListAccount(req utils.CommonRequest) ([]model.BankAccount, int, error)
	CreateAccount(bankName string, holderName string, accountNumber string, description string, accountType string) (model.BankAccount, error)
	GetAccountDetails(id int) (model.BankAccount, error)
	UpdateAccount(id int, bankName string, holderName string, accountNumber string, description string, accountType string) (model.BankAccount, error)
	DeleteAccount(id int) error
}
