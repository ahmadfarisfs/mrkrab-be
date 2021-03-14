package model

type User struct {
	BaseModel
	Fullname  string
	Username  string `gorm:"unique"`
	Role      string
	Password  string `json:"-"`
	Email     string `gorm:"unique"`
	Account   *Account
	AccountID uint
}
