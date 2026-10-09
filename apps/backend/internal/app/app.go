// Package app builds the HTTP router for the modular monolith, wiring modules in
// dependency order onto a single Gin engine. Extracted from cmd/api so both the
// server entry point and the E2E tests construct the exact same routes.
// See docs/02-system/architecture.md §4.
package app

import (
	"net/http"

	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/modules/user"
	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/config"
	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NewRouter constructs the Gin engine with all modules registered. The caller owns
// the DB connection and config. Gin mode is left to the caller (main sets release
// in prod; tests use test mode).
func NewRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Liveness/readiness. Also useful to warm Neon after auto-suspend (deployment.md §8.3).
	r.GET("/healthz", func(c *gin.Context) {
		response.OK(c, http.StatusOK, gin.H{"status": "ok", "env": cfg.AppEnv})
	})

	api := r.Group("/api/v1")

	// --- Module registration (dependency order: user -> food -> exercise ->
	//     nutrition -> workout -> summary). Wired as modules land. ---
	userModule := user.New(db, cfg)
	userModule.RegisterRoutes(api)

	return r
}
