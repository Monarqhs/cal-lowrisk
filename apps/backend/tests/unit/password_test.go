package unit

import (
	"testing"

	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TC-ON-U2 — password hashing round-trips; a wrong password fails verification,
// and the hash is never equal to the plaintext.
func TestHashAndVerifyPassword(t *testing.T) {
	const plaintext = "s3curePass"

	hash, err := auth.HashPassword(plaintext)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, plaintext, hash, "hash must not equal the plaintext")

	assert.True(t, auth.VerifyPassword(hash, plaintext), "correct password must verify")
	assert.False(t, auth.VerifyPassword(hash, "wrongPass9"), "wrong password must fail")
}

// Two hashes of the same password differ (bcrypt salts), but both verify.
func TestHashPassword_SaltedAndStillVerifies(t *testing.T) {
	const plaintext = "anotherPass1"

	h1, err := auth.HashPassword(plaintext)
	require.NoError(t, err)
	h2, err := auth.HashPassword(plaintext)
	require.NoError(t, err)

	assert.NotEqual(t, h1, h2, "bcrypt salts each hash, so two hashes differ")
	assert.True(t, auth.VerifyPassword(h1, plaintext))
	assert.True(t, auth.VerifyPassword(h2, plaintext))
}
