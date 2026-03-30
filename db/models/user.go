package models

import (
	"time"

	"github.com/google/uuid"
)

type GormModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time  `gorm:"type:timestamptz;default:now()"`
	UpdatedAt time.Time  `gorm:"type:timestamptz;default:now()"`
	DeletedAt *time.Time `gorm:"type:timestamptz;index"`
}

type RefreshToken struct {
	GormModel
	Token     string    `gorm:"uniqueIndex;not null"`
	UserID    string    `gorm:"index;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Revoked   bool      `gorm:"default:false"`
	RevokedAt *time.Time
	Email     string `gorm:"not null"`
}

type User struct {
	GormModel
	ProfileID        uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	Profile          Profile   `gorm:"foreignKey:ProfileID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Provider         string    `gorm:"not null"`
	ProviderClientID string    `gorm:"not null"`
	LastLogin        time.Time `gorm:"type:timestamp;default:now()"`
}

type Profile struct {
	GormModel
	Email      string `gorm:"unique;uniqueIndex;not null"`
	Username   string `gorm:"unique;uniqueIndex;not null"`
	FirstName  string `gorm:"not null"`
	MiddleName string
	LastName   string
	Avatar     string `gorm:"not null"`
}
