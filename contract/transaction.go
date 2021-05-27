package contract

import (
	"time"

	"github.com/ahmadfarisfs/krab-core/model"
	"github.com/ahmadfarisfs/krab-core/utils"
)

type ExpenseDestination struct {
	DestinationBankAccountID  int
	Amount                    int
	ExpenseFinancialAccountID int
}

type TransferDestination struct {
	DestinationBankAccountID int
	Amount                   int
}

//TransactionStore create model for cash basis logging
type TransactionStore interface {
	CreateIncomeTransaction(amount int, remarks string, destinationProjectID int, incomeFinancialAccountID int, sourceBankAccountID int, destinationBankAccountID int, isPaid bool, transactionTime time.Time) (model.Transaction, error)
	CreateExpenseTransaction(amount int, remarks string, sourceProjectID int, sourceBankAccountID int, expenses []ExpenseDestination, isPaid bool, transactionTime time.Time) (model.Transaction, error)
	CreateBankTransferTransaction(amount int, remarks string, transferFeeProjectID int, sourceBankAccountID int, destination []TransferDestination, transferFee ExpenseDestination, isPaid bool, transactionTime time.Time) (model.Transaction, error)
	CreateProjectTransferTransaction(amount int, remarks string, sourceProjectID int, sourceFinancialAccountID int, destinationProjectID int, destinationFinancialAccountID int, isPaid bool, transactionTime time.Time) (model.Transaction, error)

	UpdateIncomeTransaction(id int, amount int, remarks string, destinationProjectID int, incomeFinancialAccountID int, sourceBankAccountID int, destinationBankAccountID int, isPaid bool, transactionTime time.Time) (model.Transaction, error)
	UpdateExpenseTransaction(id int, amount int, remarks string, sourceProjectID int, sourceBankAccountID int, expenses []ExpenseDestination, isPaid bool, transactionTime time.Time) (model.Transaction, error)
	UpdateBankTransferTransaction(id int, amount int, remarks string, sourceBankAccountID int, destination []TransferDestination, transferFee ExpenseDestination, isPaid bool, transactionTime time.Time) (model.Transaction, error)
	UpdateFinancialAccountTransferTransaction(id int, amount int, remarks string, sourceFinancialAccountID int, destinationFinancialAccountID int, isPaid bool, transactionTime time.Time) (model.Transaction, error)

	UpdateTransactionRaw(id int, trx model.Transaction) (model.Transaction, error)

	DeleteTransaction(id int) error

	GetTransactionDetailsbyID(transactionID int) (model.Transaction, error)
	GetTransactionDetailsbyCode(transactionCode string) (model.Transaction, error)

	ListTransaction(req utils.CommonRequest) ([]model.Transaction, int, error)
}
