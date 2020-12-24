package model

type Transaction struct {
	BaseModel
	TransactionCode string
	Remarks         string
	Mutation        []Mutation `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
