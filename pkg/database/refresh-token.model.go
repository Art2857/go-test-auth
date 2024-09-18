package database

import (
	"time"

	"gorm.io/gorm"
)

// По хорошему добавить что-то вроде expires_at, но в случае redis'а это было бы проще
type RefreshTokenModel struct {
	gorm.Model
	TokenHash string    `gorm:"not null"`
	IP        string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
