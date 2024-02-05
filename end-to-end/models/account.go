package models

import (
	"time"
)

type Account struct {
	ID        string `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string
	CreatedAt *time.Time
	UpdatedAt time.Time
}
