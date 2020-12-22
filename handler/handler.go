package handler

import (
	"github.com/ahmadfarisfs/krab-core/contract"
)

type Handler struct {
	accountStore     contract.AccountStore
	transactionStore contract.TransactionStore
	projectStore     contract.ProjectStore
	userStore        contract.UserStore
}

func NewHandler(as contract.AccountStore, ts contract.TransactionStore, ps contract.ProjectStore, us contract.UserStore) *Handler {
	return &Handler{
		accountStore:     as,
		transactionStore: ts,
		projectStore:     ps,
		userStore:        us,
	}
}
