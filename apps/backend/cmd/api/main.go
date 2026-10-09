// Command api is the single entry point for the cal-lowrisk modular monolith.
// It wires modules in dependency order and registers their routes on one Gin engine.
// See docs/02-system/architecture.md §4.
package main

import (
	"log"
	"net/http"

	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/config"
	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/database"
	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("startup: %v", err)
	}

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("startup: database: %v", err)
	}
	_ = db // modules will receive this once they are wired below.

	if cfg.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Liveness/readiness. Also useful to warm Neon after auto-suspend (deployment.md §8.3).
	r.GET("/healthz", func(c *gin.Context) {
		response.OK(c, http.StatusOK, gin.H{"status": "ok", "env": cfg.AppEnv})
	})

	api := r.Group("/api/v1")

	// --- Module registration (dependency order: user -> food -> exercise ->
	//     nutrition -> workout -> summary). Wired as modules land. ---
	// user.New(db, cfg).RegisterRoutes(api)
	_ = api

	log.Printf("cal-lowrisk backend listening on :%s (env=%s)", cfg.Port, cfg.AppEnv)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server: %v", err)
	}
}
