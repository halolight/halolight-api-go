package services

import (
	"errors"
	"strings"

	"github.com/halolight/halolight-api-go/internal/models"
	"github.com/halolight/halolight-api-go/internal/repository"
	"github.com/halolight/halolight-api-go/pkg/config"
	"github.com/halolight/halolight-api-go/pkg/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already exists")
	ErrUsernameExists     = errors.New("username already exists")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrWeakPassword       = errors.New("password is too weak")
)

type AuthService interface {
	Register(email, username, password string) (*models.User, string, error)
	Login(email, password string) (*models.User, string, error)
}

type authService struct {
	cfg  config.Config
	repo repository.UserRepository
}

func NewAuthService(cfg config.Config, repo repository.UserRepository) AuthService {
	return &authService{cfg: cfg, repo: repo}
}

func (s *authService) Register(email, username, password string) (*models.User, string, error) {
	// Validate email
	email = strings.TrimSpace(strings.ToLower(email))
	if !strings.Contains(email, "@") {
		return nil, "", ErrInvalidEmail
	}

	// Validate username
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 64 {
		return nil, "", errors.New("username must be between 3 and 64 characters")
	}

	// Validate password
	if len(password) < 6 {
		return nil, "", ErrWeakPassword
	}

	// Check if email already exists
	if _, err := s.repo.GetByEmail(email); err == nil {
		return nil, "", ErrEmailExists
	}

	// Check if username already exists
	if _, err := s.repo.GetByUsername(username); err == nil {
		return nil, "", ErrUsernameExists
	}

	// Hash password
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user := &models.User{
		Email:    email,
		Username: username,
		Password: hash,
	}

	if err := s.repo.Create(user); err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return nil, "", ErrEmailExists
		}
		return nil, "", err
	}

	// Generate JWT token
	token, err := utils.GenerateAccessToken(user.ID, s.cfg.JWTSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *authService) Login(email, password string) (*models.User, string, error) {
	// Normalize email
	email = strings.TrimSpace(strings.ToLower(email))

	// Get user by email
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	// Check password
	if err := utils.CheckPassword(password, user.Password); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateAccessToken(user.ID, s.cfg.JWTSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
