package utils

import (
	"gorm.io/gorm"
	"time"
)

// GormModel is used instead of default gorm.Model
// because default does not have json fields for
// CreatedAt and UpdatedAt
type GormModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
