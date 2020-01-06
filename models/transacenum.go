package model

import "database/sql/driver"

type TransactionType int64

const (
	TransactionDebit TransactionType = iota
	TransactionCredit
)

var transactionTypes = [...]string{
	"DEBIT",
	"CREDIT",
}

func (t TransactionType) String() string {
	return transactionTypes[t]
}

func (u *TransactionType) Scan(value interface{}) error {
	*u = TransactionType(value.(int64))
	return nil
}

func (u TransactionType) Value() (driver.Value, error) { return int64(u), nil }
