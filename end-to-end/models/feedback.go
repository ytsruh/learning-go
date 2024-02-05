package models

import (
	"time"

	"gorm.io/gorm"
)

type Feedback struct {
	ID        string `gorm:"type:uuid;default:uuid_generate_v4()"`
	Title     string
	Body      *string
	UserId    string `gorm:"type:uuid"`
	User      User   `gorm:"foreignKey:UserId;references:ID"`
	CreatedAt *time.Time
	UpdatedAt time.Time
}

func (f *Feedback) TableName() string {
	return "feedback"
}

func (f *Feedback) Create(db *gorm.DB) error {
	tx := db.Create(f)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
