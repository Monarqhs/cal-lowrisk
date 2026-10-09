// Command api is the single entry point for the cal-lowrisk modular monolith.
// Router construction lives in internal/app so the server and the E2E tests share it.
// See docs/02-system/architecture.md §4.
package main

import (
	"log"

	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/app"
	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/config"
	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/database"
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

	if cfg.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := app.NewRouter(db, cfg)

	log.Printf("cal-lowrisk backend listening on :%s (env=%s)", cfg.Port, cfg.AppEnv)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server: %v", err)
	}
}
