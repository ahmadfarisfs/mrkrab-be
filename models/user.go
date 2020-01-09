package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name              string       `gorm:"type:varchar(255);not null" json:"name"`
	Username          string       `gorm:"type:varchar(255);unique;not null" json:"username"`
	ProfilePictureURL string       `gorm:"type:varchar(255)" json:"image_url"`
	Role              UserRoleType `gorm:"type:integer" json:"role"`
	Email             string       `gorm:"type:varchar(255);not null" json:"email"`
	Phone             int          `gorm:"not null" json:"phone"`
	Password          string       `gorm:"type:varchar(255);not null" json:"password"`
	//	IsDeleted         bool         `gorm:"default:false;not null" json:"-"`
}
