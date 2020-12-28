package model

type PayRec struct {
	BaseModel
	Remarks         string
	TransactionCode *string
	ProjectID       uint
	Project         Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PocketID        *uint
	Pocket          *Budget `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Amount          int
	Email           string
}
