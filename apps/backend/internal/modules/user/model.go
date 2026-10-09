// Package user is the user module (bounded context): auth, profile, role.
// It owns the "user" PostgreSQL schema. Schema (DDL) is defined by migrations
// (migrations/user/*), not GORM AutoMigrate; these structs describe existing
// columns for runtime queries only. See docs/02-system/erd.md §4.1–4.3.
package user

import (
	"time"

	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Seeded role ids (fixed UUIDs — see migrations/user/000004_seed_roles.up.sql + erd.md §2).
const (
	RoleUserName  = "user"
	RoleAdminName = "admin"
)

// Role is a first-class entity (not a boolean), per BRD §4 / USR-6.
type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	Description string    `gorm:"column:description" json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName pins Role to the "user" schema (reserved word → always quoted).
func (Role) TableName() string { return model.Qualify("user", "role") }

// User is the auth identity (USR-1/2) with a role FK (USR-6) and soft delete
// (admin activate/deactivate, ADM-5). PasswordHash is never serialised to clients.
type User struct {
	model.Base
	Email        string    `gorm:"column:email" json:"email"`
	PasswordHash string    `gorm:"column:password_hash" json:"-"`
	RoleID       uuid.UUID `gorm:"column:role_id;type:uuid" json:"roleId"`
	Role         Role      `gorm:"foreignKey:RoleID" json:"-"`
}

// TableName pins User to the "user" schema.
func (User) TableName() string { return model.Qualify("user", "users") }

// UserProfile is the 1:1 body-data/goal profile (USR-3/4). Enum-like fields mirror
// the DB CHECK constraints (erd.md §4.3). No soft delete (tied to the user lifecycle).
type UserProfile struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"column:user_id;type:uuid" json:"userId"`
	WeightKg      float64   `gorm:"column:weight_kg" json:"weightKg"`
	HeightCm      float64   `gorm:"column:height_cm" json:"heightCm"`
	Age           int       `gorm:"column:age" json:"age"`
	Sex           string    `gorm:"column:sex" json:"sex"`
	ActivityLevel string    `gorm:"column:activity_level" json:"activityLevel"`
	Goal          string    `gorm:"column:goal" json:"goal"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName pins UserProfile to the "user" schema.
func (UserProfile) TableName() string { return model.Qualify("user", "user_profile") }

// BeforeCreate assigns a UUIDv7 id if unset (profiles are created outside Base's hook).
func (p *UserProfile) BeforeCreate(_ *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = model.NewID()
	}
	return nil
}
