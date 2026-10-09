// Package e2e holds end-to-end API tests that drive the real HTTP router against a
// real PostgreSQL database (a disposable Neon branch). Each test maps 1:1 to a
// "Test cases" entry in docs/02-system/api-specs/user-service/01. Onboarding User.md
// (TC-ON-*). E2E tests skip automatically unless TEST_DATABASE_URL is set.
package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/Monarqhs/cal-lowrisk/apps/backend/tests/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// uniqueEmail returns a per-call unique email so parallel/sequential runs don't collide.
func uniqueEmail() string {
	return fmt.Sprintf("e2e.%d@example.com", time.Now().UnixNano())
}

func registerBody(email string) map[string]any {
	return map[string]any{"email": email, "password": "s3curePass", "confirmPassword": "s3curePass"}
}

func validProfileBody() map[string]any {
	return map[string]any{
		"goal": "lose", "weightKg": 72.5, "heightCm": 175, "age": 29,
		"sex": "male", "activityLevel": "moderate",
	}
}

// --- Register (TC-ON-01..07) ---

func TestRegister(t *testing.T) {
	h := helpers.New(t)
	defer h.Close()

	t.Run("TC-ON-01 valid register returns 201 with role user and no password", func(t *testing.T) {
		h.TruncateAll(t)
		email := uniqueEmail()
		res := h.POST(t, "/auth/register", registerBody(email), "")

		assert.Equal(t, 201, res.Status)
		assert.True(t, res.Body.Success)
		var data struct {
			ID    string `json:"id"`
			Email string `json:"email"`
			Role  string `json:"role"`
		}
		res.DecodeData(t, &data)
		assert.NotEmpty(t, data.ID)
		assert.Equal(t, "user", data.Role)
		assert.NotContains(t, string(res.Body.Data), "password")
	})

	t.Run("TC-ON-02 duplicate email returns 409 CONFLICT", func(t *testing.T) {
		h.TruncateAll(t)
		email := uniqueEmail()
		require.Equal(t, 201, h.POST(t, "/auth/register", registerBody(email), "").Status)

		res := h.POST(t, "/auth/register", registerBody(email), "")
		assert.Equal(t, 409, res.Status)
		require.NotNil(t, res.Body.Error)
		assert.Equal(t, "CONFLICT", res.Body.Error.Code)
	})

	t.Run("TC-ON-03 invalid email returns 400", func(t *testing.T) {
		h.TruncateAll(t)
		body := registerBody("not-an-email")
		res := h.POST(t, "/auth/register", body, "")
		assert.Equal(t, 400, res.Status)
		assert.Equal(t, "VALIDATION_ERROR", res.Body.Error.Code)
	})

	t.Run("TC-ON-04 short password returns 400", func(t *testing.T) {
		h.TruncateAll(t)
		body := map[string]any{"email": uniqueEmail(), "password": "short", "confirmPassword": "short"}
		res := h.POST(t, "/auth/register", body, "")
		assert.Equal(t, 400, res.Status)
		assert.Equal(t, "VALIDATION_ERROR", res.Body.Error.Code)
	})

	t.Run("TC-ON-05 confirmPassword mismatch returns 400", func(t *testing.T) {
		h.TruncateAll(t)
		body := map[string]any{"email": uniqueEmail(), "password": "s3curePass", "confirmPassword": "different1"}
		res := h.POST(t, "/auth/register", body, "")
		assert.Equal(t, 400, res.Status)
		assert.Equal(t, "VALIDATION_ERROR", res.Body.Error.Code)
	})

	t.Run("TC-ON-06 email stored and returned lowercased", func(t *testing.T) {
		h.TruncateAll(t)
		mixed := fmt.Sprintf("E2E.Mixed.%d@Example.COM", time.Now().UnixNano())
		res := h.POST(t, "/auth/register", registerBody(mixed), "")
		require.Equal(t, 201, res.Status)
		var data struct {
			Email string `json:"email"`
		}
		res.DecodeData(t, &data)
		assert.Equal(t, lower(mixed), data.Email)
	})

	t.Run("TC-ON-07 password persisted only as hash; login works with original", func(t *testing.T) {
		h.TruncateAll(t)
		email := uniqueEmail()
		require.Equal(t, 201, h.POST(t, "/auth/register", registerBody(email), "").Status)

		// hash is not plaintext in the DB
		var hash string
		require.NoError(t, h.DB.Raw(`SELECT password_hash FROM "user".users WHERE email = ?`, email).Scan(&hash).Error)
		assert.NotEqual(t, "s3curePass", hash)
		assert.NotEmpty(t, hash)

		// and the original plaintext still logs in
		login := h.POST(t, "/auth/login", map[string]any{"email": email, "password": "s3curePass"}, "")
		assert.Equal(t, 200, login.Status)
	})
}

// --- Login (TC-ON-08..12, TC-ON-23) ---

func TestLogin(t *testing.T) {
	h := helpers.New(t)
	defer h.Close()

	t.Run("TC-ON-08 correct creds return 200 with token and hasProfile=false", func(t *testing.T) {
		h.TruncateAll(t)
		email := uniqueEmail()
		require.Equal(t, 201, h.POST(t, "/auth/register", registerBody(email), "").Status)

		res := h.POST(t, "/auth/login", map[string]any{"email": email, "password": "s3curePass"}, "")
		assert.Equal(t, 200, res.Status)
		var data struct {
			Token string `json:"token"`
			User  struct {
				Role       string `json:"role"`
				HasProfile bool   `json:"hasProfile"`
			} `json:"user"`
		}
		res.DecodeData(t, &data)
		assert.NotEmpty(t, data.Token)
		assert.Equal(t, "user", data.User.Role)
		assert.False(t, data.User.HasProfile)
	})

	t.Run("TC-ON-09/10 wrong password and unknown email return identical 401", func(t *testing.T) {
		h.TruncateAll(t)
		email := uniqueEmail()
		require.Equal(t, 201, h.POST(t, "/auth/register", registerBody(email), "").Status)

		wrongPass := h.POST(t, "/auth/login", map[string]any{"email": email, "password": "wrongPass9"}, "")
		unknown := h.POST(t, "/auth/login", map[string]any{"email": uniqueEmail(), "password": "s3curePass"}, "")

		assert.Equal(t, 401, wrongPass.Status)
		assert.Equal(t, 401, unknown.Status)
		require.NotNil(t, wrongPass.Body.Error)
		require.NotNil(t, unknown.Body.Error)
		assert.Equal(t, "UNAUTHENTICATED", wrongPass.Body.Error.Code)
		// identical message → no account enumeration
		assert.Equal(t, wrongPass.Body.Error.Message, unknown.Body.Error.Message)
	})

	t.Run("TC-ON-12 deactivated (soft-deleted) account returns 401", func(t *testing.T) {
		h.TruncateAll(t)
		email := uniqueEmail()
		require.Equal(t, 201, h.POST(t, "/auth/register", registerBody(email), "").Status)
		// simulate admin deactivation via soft delete
		require.NoError(t, h.DB.Exec(`UPDATE "user".users SET deleted_at = now() WHERE email = ?`, email).Error)

		res := h.POST(t, "/auth/login", map[string]any{"email": email, "password": "s3curePass"}, "")
		assert.Equal(t, 401, res.Status)
		assert.Equal(t, "UNAUTHENTICATED", res.Body.Error.Code)
	})

	t.Run("TC-ON-23 after profile exists login returns hasProfile=true", func(t *testing.T) {
		h.TruncateAll(t)
		email := uniqueEmail()
		token := registerAndLogin(t, h, email)
		require.Equal(t, 201, h.POST(t, "/me/profile", validProfileBody(), token).Status)

		res := h.POST(t, "/auth/login", map[string]any{"email": email, "password": "s3curePass"}, "")
		require.Equal(t, 200, res.Status)
		var data struct {
			User struct {
				HasProfile bool `json:"hasProfile"`
			} `json:"user"`
		}
		res.DecodeData(t, &data)
		assert.True(t, data.User.HasProfile)
	})
}

// --- Create profile (TC-ON-11, TC-ON-13..22) ---

func TestCreateProfile(t *testing.T) {
	h := helpers.New(t)
	defer h.Close()

	t.Run("TC-ON-11/13 valid profile returns 201 with worked-example target", func(t *testing.T) {
		h.TruncateAll(t)
		token := registerAndLogin(t, h, uniqueEmail())

		res := h.POST(t, "/me/profile", validProfileBody(), token)
		assert.Equal(t, 201, res.Status)
		var data struct {
			Target struct {
				BMR                float64 `json:"bmr"`
				TDEE               float64 `json:"tdee"`
				DailyCalorieTarget float64 `json:"dailyCalorieTarget"`
			} `json:"target"`
		}
		res.DecodeData(t, &data)
		assert.Equal(t, 1678.75, data.Target.BMR)
		assert.Equal(t, 2602.06, data.Target.TDEE)
		assert.Equal(t, 2102.06, data.Target.DailyCalorieTarget)
	})

	t.Run("TC-ON-14 female formula branch", func(t *testing.T) {
		h.TruncateAll(t)
		token := registerAndLogin(t, h, uniqueEmail())
		body := validProfileBody()
		body["sex"] = "female"
		body["goal"] = "maintain"
		res := h.POST(t, "/me/profile", body, token)
		require.Equal(t, 201, res.Status)
		var data struct {
			Target struct {
				BMR  float64 `json:"bmr"`
				TDEE float64 `json:"tdee"`
			} `json:"target"`
		}
		res.DecodeData(t, &data)
		// female 72.5/175/29: 10*72.5 + 6.25*175 - 5*29 - 161 = 725 + 1093.75 - 145 - 161 = 1512.75
		assert.Equal(t, 1512.75, data.Target.BMR)
		assert.InDelta(t, 1512.75*1.55, data.Target.TDEE, 0.01)
	})

	t.Run("TC-ON-15 activity multipliers", func(t *testing.T) {
		mults := map[string]float64{"sedentary": 1.2, "light": 1.375, "moderate": 1.55, "active": 1.725, "very_active": 1.9}
		for level, mult := range mults {
			level, mult := level, mult
			t.Run(level, func(t *testing.T) {
				h.TruncateAll(t)
				token := registerAndLogin(t, h, uniqueEmail())
				body := validProfileBody()
				body["activityLevel"] = level
				body["goal"] = "maintain"
				res := h.POST(t, "/me/profile", body, token)
				require.Equal(t, 201, res.Status)
				var data struct {
					Target struct {
						BMR  float64 `json:"bmr"`
						TDEE float64 `json:"tdee"`
					} `json:"target"`
				}
				res.DecodeData(t, &data)
				assert.InDelta(t, data.Target.BMR*mult, data.Target.TDEE, 0.01)
			})
		}
	})

	t.Run("TC-ON-16 goal adjustments", func(t *testing.T) {
		deltas := map[string]float64{"lose": -500, "maintain": 0, "gain": 300}
		for goal, delta := range deltas {
			goal, delta := goal, delta
			t.Run(goal, func(t *testing.T) {
				h.TruncateAll(t)
				token := registerAndLogin(t, h, uniqueEmail())
				body := validProfileBody()
				body["goal"] = goal
				res := h.POST(t, "/me/profile", body, token)
				require.Equal(t, 201, res.Status)
				var data struct {
					Target struct {
						TDEE               float64 `json:"tdee"`
						DailyCalorieTarget float64 `json:"dailyCalorieTarget"`
					} `json:"target"`
				}
				res.DecodeData(t, &data)
				assert.InDelta(t, data.Target.TDEE+delta, data.Target.DailyCalorieTarget, 0.01)
			})
		}
	})

	t.Run("TC-ON-17 no token returns 401", func(t *testing.T) {
		h.TruncateAll(t)
		res := h.POST(t, "/me/profile", validProfileBody(), "")
		assert.Equal(t, 401, res.Status)
		assert.Equal(t, "UNAUTHENTICATED", res.Body.Error.Code)
	})

	t.Run("TC-ON-18 invalid/expired token returns 401", func(t *testing.T) {
		h.TruncateAll(t)
		// malformed
		assert.Equal(t, 401, h.POST(t, "/me/profile", validProfileBody(), "garbage.token.here").Status)
		// expired (issued with negative TTL)
		expired := h.IssueToken(t, uuid.New(), "user", -time.Hour)
		res := h.POST(t, "/me/profile", validProfileBody(), expired)
		assert.Equal(t, 401, res.Status)
		assert.Equal(t, "UNAUTHENTICATED", res.Body.Error.Code)
	})

	t.Run("TC-ON-19 admin-role token returns 403", func(t *testing.T) {
		h.TruncateAll(t)
		adminToken := h.IssueToken(t, uuid.New(), "admin", time.Hour)
		res := h.POST(t, "/me/profile", validProfileBody(), adminToken)
		assert.Equal(t, 403, res.Status)
		assert.Equal(t, "FORBIDDEN", res.Body.Error.Code)
	})

	t.Run("TC-ON-20 second profile returns 409", func(t *testing.T) {
		h.TruncateAll(t)
		token := registerAndLogin(t, h, uniqueEmail())
		require.Equal(t, 201, h.POST(t, "/me/profile", validProfileBody(), token).Status)
		res := h.POST(t, "/me/profile", validProfileBody(), token)
		assert.Equal(t, 409, res.Status)
		assert.Equal(t, "CONFLICT", res.Body.Error.Code)
	})

	t.Run("TC-ON-21 out-of-range age returns 400", func(t *testing.T) {
		h.TruncateAll(t)
		token := registerAndLogin(t, h, uniqueEmail())
		for _, age := range []int{0, 200} {
			body := validProfileBody()
			body["age"] = age
			res := h.POST(t, "/me/profile", body, token)
			assert.Equal(t, 400, res.Status, "age=%d", age)
			assert.Equal(t, "VALIDATION_ERROR", res.Body.Error.Code)
		}
	})

	t.Run("TC-ON-22 invalid enum values return 400", func(t *testing.T) {
		h.TruncateAll(t)
		token := registerAndLogin(t, h, uniqueEmail())
		for field, bad := range map[string]string{"sex": "other", "goal": "bulk", "activityLevel": "extreme"} {
			body := validProfileBody()
			body[field] = bad
			res := h.POST(t, "/me/profile", body, token)
			assert.Equal(t, 400, res.Status, "field=%s", field)
			assert.Equal(t, "VALIDATION_ERROR", res.Body.Error.Code)
		}
	})
}

// registerAndLogin registers a fresh user and returns a valid user-role JWT.
func registerAndLogin(t *testing.T, h *helpers.Harness, email string) string {
	t.Helper()
	require.Equal(t, 201, h.POST(t, "/auth/register", registerBody(email), "").Status)
	res := h.POST(t, "/auth/login", map[string]any{"email": email, "password": "s3curePass"}, "")
	require.Equal(t, 200, res.Status)
	var data struct {
		Token string `json:"token"`
	}
	res.DecodeData(t, &data)
	require.NotEmpty(t, data.Token)
	return data.Token
}

func lower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}
