package model

import (
	"time"
)

type BaseModel struct {
	ID        uint      `gorm:"primarykey"` // json:"id"`
	CreatedAt time.Time //`json:"created_at"`
	UpdatedAt time.Time
}
