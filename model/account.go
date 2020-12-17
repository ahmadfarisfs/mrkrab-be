package model

type Account struct {
	BaseModel
	AccountName  string `gorm:"unique"`
	ParentID     *uint
	Parent       *Account
	Balance      int
	TotalIncome  int
	TotalExpense int
}
