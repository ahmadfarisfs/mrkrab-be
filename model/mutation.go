package model

type Mutation struct {
	BaseModel
	AccountID     int
	Account       Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TransactionID int
	Transaction   Transaction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Amount        int         //deltas
}
