package models

// По хорошему добавить что-то вроде expires_at, но в случае redis'а это было бы проще
type RefreshTokenModel struct {
	Model
	TokenHash string `gorm:"not null"`
	IP        string `gorm:"not null"`
}
