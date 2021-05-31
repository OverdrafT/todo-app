package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Item struct {
	ID        string         `json:"-" gorm:"primaryKey"`
	Title     string         `json:"title,omitempty"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" sql:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (i *Item) BeforeCreate(tx *gorm.DB) error {
	i.ID = uuid.New().String()
	return nil
}
