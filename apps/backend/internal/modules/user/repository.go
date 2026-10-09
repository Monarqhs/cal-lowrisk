package user

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ErrNotFound is returned by the repository when a row does not exist.
var ErrNotFound = errors.New("not found")

// Repository is the data-access boundary for the user module (GORM only, no business
// logic). Service depends on this interface so it can be mocked in unit tests.
type Repository interface {
	// RoleIDByName resolves a seeded role's id by its name (e.g. "user").
	RoleIDByName(name string) (uuid.UUID, error)
	// EmailExists reports whether an email is already taken (active OR soft-deleted).
	EmailExists(email string) (bool, error)
	// CreateUser inserts a new user row.
	CreateUser(u *User) error
	// FindUserByEmail returns an active (non-soft-deleted) user with its role, or ErrNotFound.
	FindUserByEmail(email string) (*User, error)
	// HasProfile reports whether a profile row exists for the user.
	HasProfile(userID uuid.UUID) (bool, error)
	// CreateProfile inserts a new user_profile row.
	CreateProfile(p *UserProfile) error
}

type repository struct{ db *gorm.DB }

// NewRepository builds the GORM-backed Repository.
func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) RoleIDByName(name string) (uuid.UUID, error) {
	var role Role
	if err := r.db.Where("name = ?", name).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return uuid.Nil, ErrNotFound
		}
		return uuid.Nil, err
	}
	return role.ID, nil
}

func (r *repository) EmailExists(email string) (bool, error) {
	var count int64
	// Unscoped so a soft-deleted account still reserves its email (spec Endpoint 1 note).
	if err := r.db.Unscoped().Model(&User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) CreateUser(u *User) error {
	return r.db.Create(u).Error
}

func (r *repository) FindUserByEmail(email string) (*User, error) {
	var u User
	// Default scope excludes soft-deleted rows → deactivated accounts can't log in.
	if err := r.db.Preload("Role").Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *repository) HasProfile(userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&UserProfile{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) CreateProfile(p *UserProfile) error {
	return r.db.Create(p).Error
}
