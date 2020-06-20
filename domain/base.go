package domain

import (
	"time"

	"gorm.io/gorm"
)

//BaseModel is basis for all MRKRAB models
type BaseModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
