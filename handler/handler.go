package handler

import (
	"github.com/ahmadfarisfs/krab-core/store"
)

type Handler struct {
	accountStore     *store.AccountStore
	transactionStore *store.TransactionStore
	projectStore     *store.ProjectStore
	userStore        *store.UserStore
	mutationStore    *store.MutationStore
	payRecStore      *store.PayRecStore
}

func NewHandler(as *store.AccountStore, ts *store.TransactionStore,
	ps *store.ProjectStore, us *store.UserStore,
	ms *store.MutationStore, prs *store.PayRecStore) *Handler {

	return &Handler{
		accountStore:     as,
		transactionStore: ts,
		projectStore:     ps,
		userStore:        us,
		mutationStore:    ms,
		payRecStore:      prs,
	}
}
