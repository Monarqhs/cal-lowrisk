package user

import (
	"errors"
	"net/http"

	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/modules/user/dto"
	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/auth"
	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/response"
	"github.com/gin-gonic/gin"
)

// Controller holds the user module's HTTP handlers. No business logic lives here —
// it parses/validates input, calls the service, and formats the response envelope.
type Controller struct{ svc Service }

// NewController builds the user Controller.
func NewController(svc Service) *Controller { return &Controller{svc: svc} }

// Register handles POST /auth/register.
func (ctl *Controller) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	res, err := ctl.svc.Register(req)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailTaken):
			response.Fail(c, http.StatusConflict, "CONFLICT", "Email is already registered.")
		default:
			response.Fail(c, http.StatusInternalServerError, "INTERNAL", "Could not create the account.")
		}
		return
	}
	response.OK(c, http.StatusCreated, res)
}

// Login handles POST /auth/login.
func (ctl *Controller) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	res, err := ctl.svc.Login(req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCreds):
			// Same response for unknown email and wrong password (no enumeration).
			response.Fail(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Invalid email or password.")
		default:
			response.Fail(c, http.StatusInternalServerError, "INTERNAL", "Could not log in.")
		}
		return
	}
	response.OK(c, http.StatusOK, res)
}

// CreateProfile handles POST /me/profile (requires auth + user role).
func (ctl *Controller) CreateProfile(c *gin.Context) {
	userID, ok := auth.UserID(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, "UNAUTHENTICATED", "Authentication required.")
		return
	}
	var req dto.CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	res, err := ctl.svc.CreateProfile(userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrProfileExists):
			response.Fail(c, http.StatusConflict, "CONFLICT", "A profile already exists for this user.")
		default:
			response.Fail(c, http.StatusInternalServerError, "INTERNAL", "Could not create the profile.")
		}
		return
	}
	response.OK(c, http.StatusCreated, res)
}
