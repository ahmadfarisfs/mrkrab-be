package model

type Account struct {
	BaseModel
	AccountName  string `gorm:"unique"`
	ParentID     *uint
	Parent       *Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Balance      int
	TotalIncome  int
	TotalExpense int
}
