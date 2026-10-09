package user

import (
	"errors"
	"strings"

	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/modules/user/dto"
	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/auth"
	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/model"
	"github.com/google/uuid"
)

// Service-level sentinel errors. The controller maps these to the spec's HTTP
// status + error.code (see docs/02-system/api-specs/.../01. Onboarding User.md).
var (
	ErrEmailTaken        = errors.New("email already registered")
	ErrInvalidCreds      = errors.New("invalid credentials")
	ErrProfileExists     = errors.New("profile already exists")
	ErrRoleMisconfigured = errors.New("default role not found (seed missing)")
)

// Service holds the user module's business logic. It depends on the Repository
// interface (not GORM) so it is unit-testable with a mock (spec TC-ON-U*).
type Service interface {
	Register(req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	CreateProfile(userID uuid.UUID, req dto.CreateProfileRequest) (*dto.CreateProfileResponse, error)
}

type service struct {
	repo   Repository
	issuer *auth.Issuer
}

// NewService wires the user service with its repository and the JWT issuer.
func NewService(repo Repository, issuer *auth.Issuer) Service {
	return &service{repo: repo, issuer: issuer}
}

// Register creates a new account with the default "user" role (USR-1/6). Email is
// normalised to lowercase; the password is stored only as a bcrypt hash.
func (s *service) Register(req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	email := normalizeEmail(req.Email)

	exists, err := s.repo.EmailExists(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailTaken
	}

	roleID, err := s.repo.RoleIDByName(RoleUserName)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrRoleMisconfigured
		}
		return nil, err
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &User{Email: email, PasswordHash: hash, RoleID: roleID}
	if err := s.repo.CreateUser(u); err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		ID:        u.ID.String(),
		Email:     u.Email,
		Role:      RoleUserName,
		CreatedAt: u.CreatedAt,
	}, nil
}

// Login verifies credentials and issues a JWT (USR-2). Unknown email and wrong
// password both return ErrInvalidCreds (no account enumeration). Soft-deleted
// accounts are excluded by the repository, so they also fail with ErrInvalidCreds.
func (s *service) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	email := normalizeEmail(req.Email)

	u, err := s.repo.FindUserByEmail(email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrInvalidCreds
		}
		return nil, err
	}
	if !auth.VerifyPassword(u.PasswordHash, req.Password) {
		return nil, ErrInvalidCreds
	}

	token, err := s.issuer.Issue(u.ID, u.Role.Name)
	if err != nil {
		return nil, err
	}

	hasProfile, err := s.repo.HasProfile(u.ID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: int(s.issuer.TTL().Seconds()),
		User: dto.LoginUser{
			ID:         u.ID.String(),
			Email:      u.Email,
			Role:       u.Role.Name,
			HasProfile: hasProfile,
		},
	}, nil
}

// CreateProfile stores the caller's 1:1 profile and returns it with the computed
// (not stored) daily target (USR-3/4, erd.md §5). Enforces the 1:1 rule.
func (s *service) CreateProfile(userID uuid.UUID, req dto.CreateProfileRequest) (*dto.CreateProfileResponse, error) {
	has, err := s.repo.HasProfile(userID)
	if err != nil {
		return nil, err
	}
	if has {
		return nil, ErrProfileExists
	}

	p := &UserProfile{
		ID:            model.NewID(),
		UserID:        userID,
		WeightKg:      req.WeightKg,
		HeightCm:      req.HeightCm,
		Age:           req.Age,
		Sex:           req.Sex,
		ActivityLevel: req.ActivityLevel,
		Goal:          req.Goal,
	}
	if err := s.repo.CreateProfile(p); err != nil {
		return nil, err
	}

	t := ComputeTarget(p.Sex, p.WeightKg, p.HeightCm, p.Age, p.ActivityLevel, p.Goal)

	return &dto.CreateProfileResponse{
		Profile: dto.ProfileView{
			ID:            p.ID.String(),
			UserID:        p.UserID.String(),
			Goal:          p.Goal,
			WeightKg:      p.WeightKg,
			HeightCm:      p.HeightCm,
			Age:           p.Age,
			Sex:           p.Sex,
			ActivityLevel: p.ActivityLevel,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
		},
		Target: dto.TargetView{
			BMR:                t.BMR,
			TDEE:               t.TDEE,
			DailyCalorieTarget: t.DailyCalorieTarget,
			Formula:            "mifflin_st_jeor",
		},
	}, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
