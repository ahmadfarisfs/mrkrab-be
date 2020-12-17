package handler

import (
	"github.com/ahmadfarisfs/krab-core/contract"
)

type Handler struct {
	accountStore     contract.AccountStore
	transactionStore contract.TransactionStore
	projectStore     contract.ProjectStore
}

func NewHandler(as contract.AccountStore, ts contract.TransactionStore, ps contract.ProjectStore) *Handler {
	return &Handler{
		accountStore:     as,
		transactionStore: ts,
		projectStore:     ps,
	}
}
