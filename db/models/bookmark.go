package models

import (
	"github.com/google/uuid"
)

type EnrichmentStatus string

const (
	EnrichmentPending EnrichmentStatus = "pending"
	EnrichmentDone    EnrichmentStatus = "done"
	EnrichmentFailed  EnrichmentStatus = "failed"
)

type Bookmark struct {
	GormModel
	ProfileID uuid.UUID `gorm:"type:uuid;not null;index"`
	Profile   Profile   `gorm:"foreignKey:ProfileID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	URL       string  `gorm:"type:text;not null"`
	Title     *string `gorm:"type:text"`
	Notes     *string `gorm:"type:text"`
	CursorKey string  `gorm:"type:text;not null"`

	Status EnrichmentStatus `gorm:"type:text;not null;default:'pending'"`
	Tags   []Tag            `gorm:"many2many:bookmark_tags;"`
}

type Tag struct {
	GormModel
	ProfileID uuid.UUID `gorm:"type:uuid;not null;index"`
	Profile   Profile   `gorm:"foreignKey:ProfileID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Name      string     `gorm:"type:text;not null"`
	Bookmarks []Bookmark `gorm:"many2many:bookmark_tags;"`
}

type BookmarkTag struct {
	BookmarkID uuid.UUID `gorm:"type:uuid;not null;primaryKey"`
	TagID      uuid.UUID `gorm:"type:uuid;not null;primaryKey"`

	Bookmark Bookmark `gorm:"foreignKey:BookmarkID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tag      Tag      `gorm:"foreignKey:TagID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
