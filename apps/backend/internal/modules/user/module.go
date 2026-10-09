package user

import (
	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/auth"
	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Module is the user bounded context: it owns construction of its own
// repo→service→controller graph and registers its routes. main.go constructs it
// first (other modules depend on it) and calls RegisterRoutes.
type Module struct {
	ctl    *Controller
	issuer *auth.Issuer
}

// New builds the user module from the shared DB and config.
func New(db *gorm.DB, cfg *config.Config) *Module {
	issuer := auth.NewIssuer(cfg.JWTSecret, cfg.JWTTTL)
	repo := NewRepository(db)
	svc := NewService(repo, issuer)
	return &Module{ctl: NewController(svc), issuer: issuer}
}

// Issuer exposes the JWT issuer so other modules can build auth middleware with the
// same signing key (cross-module reuse through this module, not shared global state).
func (m *Module) Issuer() *auth.Issuer { return m.issuer }

// RegisterRoutes mounts the user module's endpoints on the given /api/v1 group.
// Paths match docs/02-system/api-specs/user-service/01. Onboarding User.md.
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
	authGrp := api.Group("/auth")
	{
		authGrp.POST("/register", m.ctl.Register)
		authGrp.POST("/login", m.ctl.Login)
	}

	me := api.Group("/me")
	me.Use(auth.RequireAuth(m.issuer), auth.RequireRole(RoleUserName))
	{
		me.POST("/profile", m.ctl.CreateProfile)
	}
}
