package contract

import (
	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
)

type MutationStore interface {
	ListBankAccountMutation(req utils.CommonRequest) ([]model.BankAccountMutation, int, error)
	ListFinancialAccountMutation(req utils.CommonRequest) ([]model.FinancialAccountMutation, int, error)
}
