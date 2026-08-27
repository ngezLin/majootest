package services

import (
	"context"
	"errors"
	"fmt"
	"majootest/case2-go/internal/config"
	"majootest/case2-go/internal/models"
	"majootest/case2-go/internal/repositories"
	"majootest/case2-go/internal/utils"
)

// AuthService defines the business logic for user authentication.
type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
	GetProfile(ctx context.Context, userID int64) (*models.User, error)
}

type authService struct {
	userRepo repositories.UserRepository
	cfg      *config.Config
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(userRepo repositories.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, models.ErrUserAlreadyExists
	}
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT Token
	token, err := utils.GenerateJWT(user.ID, user.Email, s.cfg.JWTSecret, s.cfg.JWTExpiryHours)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, models.ErrInvalidCredentials
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, s.cfg.JWTSecret, s.cfg.JWTExpiryHours)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) GetProfile(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
