// Package dto holds the user module's request/response data-transfer objects.
// Validation tags mirror the contract in
// docs/02-system/api-specs/user-service/01. Onboarding User.md and the DB CHECK
// constraints in erd.md §4.3. Keep them in sync with the spec (contract-first).
package dto

// RegisterRequest — POST /auth/register (USR-1).
type RegisterRequest struct {
	Email           string `json:"email" binding:"required,email,max=255"`
	Password        string `json:"password" binding:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}

// LoginRequest — POST /auth/login (USR-2).
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// CreateProfileRequest — POST /me/profile (USR-3/4). Enum values mirror DB CHECKs.
type CreateProfileRequest struct {
	Goal          string  `json:"goal" binding:"required,oneof=lose maintain gain"`
	WeightKg      float64 `json:"weightKg" binding:"required,gt=0,lte=500"`
	HeightCm      float64 `json:"heightCm" binding:"required,gt=0,lte=300"`
	Age           int     `json:"age" binding:"required,gt=0,lt=150"`
	Sex           string  `json:"sex" binding:"required,oneof=male female"`
	ActivityLevel string  `json:"activityLevel" binding:"required,oneof=sedentary light moderate active very_active"`
}
