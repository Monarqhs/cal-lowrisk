// Package model holds the shared base entity reused by every module's GORM models,
// keeping id/timestamps uniform across modules (see add-module skill + erd.md §2).
package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base is embedded by every table's model. Primary keys are native uuid, generated
// in the app as UUIDv7 (time-ordered) — see erd.md §2. Schema itself is defined by
// migrations, not GORM; these tags only describe existing columns for queries.
type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

// BeforeCreate assigns a UUIDv7 if the caller did not set one.
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		b.ID = id
	}
	return nil
}

// NewID returns a fresh UUIDv7, for app-side id generation outside GORM hooks
// (e.g. seeds, services). Falls back to a random UUID only if v7 generation fails.
func NewID() uuid.UUID {
	if id, err := uuid.NewV7(); err == nil {
		return id
	}
	return uuid.New()
}
