package dto

import "time"

// RegisterResponse — the created account (never password/hash). Matches spec Endpoint 1.
type RegisterResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

// LoginResponse — JWT + minimal user info. Matches spec Endpoint 2.
type LoginResponse struct {
	Token     string    `json:"token"`
	TokenType string    `json:"tokenType"` // always "Bearer"
	ExpiresIn int       `json:"expiresIn"` // seconds
	User      LoginUser `json:"user"`
}

// LoginUser is the compact user block inside LoginResponse.
type LoginUser struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	HasProfile bool   `json:"hasProfile"`
}

// CreateProfileResponse — saved profile + computed (not stored) target. Matches Endpoint 3.
type CreateProfileResponse struct {
	Profile ProfileView `json:"profile"`
	Target  TargetView  `json:"target"`
}

// ProfileView mirrors the user_profile row for clients.
type ProfileView struct {
	ID            string    `json:"id"`
	UserID        string    `json:"userId"`
	Goal          string    `json:"goal"`
	WeightKg      float64   `json:"weightKg"`
	HeightCm      float64   `json:"heightCm"`
	Age           int       `json:"age"`
	Sex           string    `json:"sex"`
	ActivityLevel string    `json:"activityLevel"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// TargetView is the on-the-fly target breakdown (erd.md §5). Computed server-side.
type TargetView struct {
	BMR                float64 `json:"bmr"`
	TDEE               float64 `json:"tdee"`
	DailyCalorieTarget float64 `json:"dailyCalorieTarget"`
	Formula            string  `json:"formula"` // "mifflin_st_jeor"
}
