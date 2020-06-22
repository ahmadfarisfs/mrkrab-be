package domain

import (
	"time"

	"gorm.io/gorm"
)

//BaseModel is basis for all MRKRAB models
type BaseModel struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
