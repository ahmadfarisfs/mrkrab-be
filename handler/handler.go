package handler

import (
	"github.com/ahmadfarisfs/krab-core/contract"
)

type Handler struct {
	accountStore     contract.FinancialAccountStore
	bankStore        contract.BankAccountStore
	transactionStore contract.TransactionStore
	projectStore     contract.ProjectStore
	userStore        contract.UserStore
	mutationStore    contract.MutationStore
	authStore        contract.AuthStore
}

func NewHandler(as contract.FinancialAccountStore, bs contract.BankAccountStore,
	ts contract.TransactionStore,
	ps contract.ProjectStore,
	us contract.UserStore,
	ms contract.MutationStore,
	// hs contract.AuthStore,
) *Handler {

	return &Handler{
		accountStore:     as,
		transactionStore: ts,
		projectStore:     ps,
		userStore:        us,
		mutationStore:    ms,
		// authStore:        hs,
		bankStore: bs,
	}
}
