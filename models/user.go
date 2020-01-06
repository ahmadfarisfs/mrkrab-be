package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name              string
	ProfilePictureURL string
	Role              UserRoleType
	Email             string
	Phone             int
	IsDeleted         bool
}
