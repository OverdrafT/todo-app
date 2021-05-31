package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string     `json:"-" gorm:"primaryKey"`
	Name      string     `json:"name,omitempty"`
	Email     string     `json:"email,omitempty" gorm:"type:varchar(100);unique_index"`
	Gender    string     `json:"gender"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

func (i *User) BeforeCreate(tx *gorm.DB) error {
	i.ID = uuid.New().String()
	return nil
}
