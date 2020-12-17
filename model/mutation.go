package model

type Mutation struct {
	BaseModel
	AccountID     int
	Account       Account
	TransactionID int
	Transaction   Transaction
	Amount        int //deltas
}
