package model

type Project struct {
	BaseModel
	Name    string `gorm:"unique"`
	Amount  int
	Balance int
	// IsOpen      bool
	Description string
	Status      string
}

// type Budget struct {
// 	BaseModel
// 	Name      string
// 	ProjectID uint
// 	Project   Project `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
// 	AccountID uint
// 	Account   FinancialAccount `json:"Account" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
// 	Limit     *uint
// }
