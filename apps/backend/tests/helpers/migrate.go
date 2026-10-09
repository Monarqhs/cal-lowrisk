package helpers

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// applyUserSchema runs the user module's *.up.sql migrations against the test DB, in
// filename order, so a freshly-created (empty) Neon branch gets the full "user" schema
// + seeded roles. Idempotent: migrations use CREATE ... IF NOT EXISTS / ON CONFLICT.
//
// This mirrors what golang-migrate / the Neon Power does in real environments, but
// keeps E2E tests self-contained (no external migrate CLI needed).
func applyUserSchema(t *testing.T, db *gorm.DB) {
	t.Helper()
	dir := userMigrationsDir(t)

	entries, err := os.ReadDir(dir)
	require.NoError(t, err, "read migrations dir")

	var ups []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".up.sql") {
			ups = append(ups, e.Name())
		}
	}
	sort.Strings(ups) // 000001.. before 000002.. etc.
	require.NotEmpty(t, ups, "expected user up-migrations")

	for _, name := range ups {
		sqlBytes, err := os.ReadFile(filepath.Join(dir, name))
		require.NoError(t, err, "read %s", name)
		require.NoError(t, db.Exec(string(sqlBytes)).Error, "apply %s", name)
	}
}

// userMigrationsDir resolves apps/backend/migrations/user relative to this source file,
// so the helper works regardless of the test's working directory.
func userMigrationsDir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtime.Caller")
	// this file: apps/backend/tests/helpers/migrate.go
	backendRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")
	return filepath.Join(backendRoot, "migrations", "user")
}
