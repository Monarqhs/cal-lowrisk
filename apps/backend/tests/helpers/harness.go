// Package helpers provides shared E2E test infrastructure: a test HTTP server backed
// by a real PostgreSQL database, plus request/JSON helpers. E2E tests connect to a
// disposable database (a Neon branch created from uat, dropped after). See
// docs/02-system/api-specs/README.md and the test suite README.
//
// E2E tests are SKIPPED unless TEST_DATABASE_URL is set, so `go test ./...` on a
// machine without a test DB still runs the unit tests cleanly.
package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/app"
	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// EnvTestDBURL is the env var holding the E2E test database connection string.
const EnvTestDBURL = "TEST_DATABASE_URL"

// Harness is a running test server plus its DB handle.
type Harness struct {
	Server *httptest.Server
	DB     *gorm.DB
	cfg    *config.Config
}

// SkipIfNoDB skips the calling test when TEST_DATABASE_URL is not configured.
func SkipIfNoDB(t *testing.T) string {
	t.Helper()
	url := os.Getenv(EnvTestDBURL)
	if url == "" {
		t.Skipf("%s not set — skipping E2E test (set it to a disposable Neon branch URL)", EnvTestDBURL)
	}
	return url
}

// New boots a Harness: connects to the test DB, applies the user-module schema, and
// starts an httptest server with the real router. Call Close() when done.
func New(t *testing.T) *Harness {
	t.Helper()
	url := SkipIfNoDB(t)
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger:                 gormlogger.Default.LogMode(gormlogger.Silent),
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err, "connect test DB")

	cfg := &config.Config{
		AppEnv:    "test",
		Port:      "0",
		JWTSecret: "test-secret",
		JWTTTL:    24 * time.Hour,
	}

	applyUserSchema(t, db)

	router := app.NewRouter(db, cfg)
	srv := httptest.NewServer(router)

	return &Harness{Server: srv, DB: db, cfg: cfg}
}

// Close shuts down the server and DB connection.
func (h *Harness) Close() {
	h.Server.Close()
	if sqlDB, err := h.DB.DB(); err == nil {
		_ = sqlDB.Close()
	}
}

// TruncateAll clears user data between tests, keeping seeded roles. Order respects FKs.
func (h *Harness) TruncateAll(t *testing.T) {
	t.Helper()
	require.NoError(t, h.DB.Exec(`TRUNCATE "user".user_profile, "user".users RESTART IDENTITY CASCADE`).Error)
}

// --- HTTP helpers ---

// Response is a decoded API envelope for assertions.
type Response struct {
	Status int
	Body   struct {
		Success bool            `json:"success"`
		Data    json.RawMessage `json:"data"`
		Error   *struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
}

// POST sends a JSON POST to a path under /api/v1 with an optional bearer token.
func (h *Harness) POST(t *testing.T, path string, body any, bearer string) Response {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, json.NewEncoder(&buf).Encode(body))

	req, err := http.NewRequest(http.MethodPost, h.Server.URL+"/api/v1"+path, &buf)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	var out Response
	out.Status = resp.StatusCode
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&out.Body))
	return out
}

// DecodeData unmarshals the envelope's data field into v.
func (r Response) DecodeData(t *testing.T, v any) {
	t.Helper()
	require.NoError(t, json.Unmarshal(r.Body.Data, v))
}

// IssueToken mints a JWT with the harness's secret for a given user id + role.
// Used to construct tokens the API won't otherwise hand out (e.g. an admin-role
// token for the FORBIDDEN test, or a deliberately-expired token).
func (h *Harness) IssueToken(t *testing.T, userID uuid.UUID, role string, ttl time.Duration) string {
	t.Helper()
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":  userID.String(),
		"role": role,
		"iat":  now.Unix(),
		"exp":  now.Add(ttl).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString([]byte(h.cfg.JWTSecret))
	require.NoError(t, err)
	return s
}
