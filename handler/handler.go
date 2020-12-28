package handler

import (
	"github.com/ahmadfarisfs/krab-core/contract"
)

type Handler struct {
	accountStore     contract.AccountStore
	transactionStore contract.TransactionStore
	projectStore     contract.ProjectStore
	userStore        contract.UserStore
	mutationStore    contract.MutationStore
	payRecStore      contract.PayRecStore
}

func NewHandler(as contract.AccountStore, ts contract.TransactionStore,
	ps contract.ProjectStore, us contract.UserStore,
	ms contract.MutationStore, prs contract.PayRecStore) *Handler {

	return &Handler{
		accountStore:     as,
		transactionStore: ts,
		projectStore:     ps,
		userStore:        us,
		mutationStore:    ms,
		payRecStore:      prs,
	}
}
