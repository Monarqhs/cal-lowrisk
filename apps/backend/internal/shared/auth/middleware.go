package auth

import (
	"net/http"
	"strings"

	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Context keys for values the middleware stores on the Gin context.
const (
	ctxUserID = "auth.userID"
	ctxRole   = "auth.role"
)

// RequireAuth validates the Bearer token and stores the user id + role on the context.
// On failure it writes 401 UNAUTHENTICATED and aborts.
func RequireAuth(issuer *Issuer) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(raw, prefix) {
			response.Fail(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Missing or malformed Authorization header.")
			c.Abort()
			return
		}
		claims, err := issuer.Parse(strings.TrimPrefix(raw, prefix))
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Invalid or expired token.")
			c.Abort()
			return
		}
		uid, err := uuid.Parse(claims.Subject)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Invalid token subject.")
			c.Abort()
			return
		}
		c.Set(ctxUserID, uid)
		c.Set(ctxRole, claims.Role)
		c.Next()
	}
}

// RequireRole enforces that the authenticated caller has one of the allowed roles.
// Must run after RequireAuth. On mismatch it writes 403 FORBIDDEN and aborts.
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role, _ := c.Get(ctxRole)
		rs, _ := role.(string)
		if _, ok := allowed[rs]; !ok {
			response.Fail(c, http.StatusForbidden, "FORBIDDEN", "You do not have access to this resource.")
			c.Abort()
			return
		}
		c.Next()
	}
}

// UserID returns the authenticated caller's id from the context (set by RequireAuth).
func UserID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get(ctxUserID)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}
