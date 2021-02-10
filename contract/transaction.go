package contract

import (
	"time"

	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
)

//TransactionStore create model for cash basis logging
type TransactionStore interface {
	CreateTransaction(accountID int, amount int, remarks string, SoD string, trxTime time.Time) (model.Transaction, error)
	CreateTransfer(accountFrom int, accountTo int, amount uint, remarks string) (model.Transaction, error)
	GetTransactionDetailsbyID(transactionID int) (model.Transaction, error)
	GetTransactionDetailsbyCode(transactionCode string) (model.Transaction, error)
	ListTransaction(req utils.CommonRequest) ([]model.Transaction, int, error)
}
