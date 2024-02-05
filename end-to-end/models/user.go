package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID            string     `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name          string     `gorm:"default:'User Name'" json:"name"`
	Email         string     `gorm:"unique" json:"email"`
	Password      string     `json:"-"`
	ProfileImage  *string    `json:"profileImage"`
	ShowBooks     bool       `gorm:"default:true" json:"showBooks"`
	ShowDocuments bool       `gorm:"default:true" json:"showDocuments"`
	AccountId     string     `json:"accountId"`
	Account       Account    `gorm:"foreignKey:AccountId;references:ID" json:"account"`
	CreatedAt     *time.Time `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

func Login(db *gorm.DB, email string, password string) (User, error) {
	var user User
	// Find User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return user, errors.New("user not found")
	}
	// Compare Passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, errors.New("invalid password")
	}
	return user, nil
}
