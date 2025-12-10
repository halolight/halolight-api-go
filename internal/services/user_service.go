package services

import (
	"errors"
	"strings"

	"github.com/halolight/halolight-api-go/internal/models"
	"github.com/halolight/halolight-api-go/internal/repository"
	"github.com/halolight/halolight-api-go/pkg/utils"
)

var ErrUserNotFound = errors.New("user not found")

type UserService interface {
	List(page, pageSize int) ([]models.User, int64, error)
	Get(id uint) (*models.User, error)
	GetByID(id string) (*models.User, error)
	Create(email, username, password string) (*models.User, error)
	Update(id uint, email, username, password string) (*models.User, error)
	UpdateStatus(id string, status string) (*models.User, error)
	Delete(id uint) error
	BatchDelete(ids []string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) List(page, pageSize int) ([]models.User, int64, error) {
	// Set default values
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize)
}

func (s *userService) Get(id uint) (*models.User, error) {
	u, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (s *userService) Create(email, username, password string) (*models.User, error) {
	// Validate email
	email = strings.TrimSpace(strings.ToLower(email))
	if !strings.Contains(email, "@") {
		return nil, errors.New("invalid email format")
	}

	// Validate username
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 64 {
		return nil, errors.New("username must be between 3 and 64 characters")
	}

	// Validate password
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	// Hash password
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	u := &models.User{
		Email:    email,
		Username: username,
		Password: hash,
	}

	if err := s.repo.Create(u); err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return nil, errors.New("email or username already exists")
		}
		return nil, err
	}

	return u, nil
}

func (s *userService) Update(id uint, email, username, password string) (*models.User, error) {
	// Get existing user
	u, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	// Update email
	if email != "" {
		email = strings.TrimSpace(strings.ToLower(email))
		if !strings.Contains(email, "@") {
			return nil, errors.New("invalid email format")
		}
		u.Email = email
	}

	// Update username
	if username != "" {
		username = strings.TrimSpace(username)
		if len(username) < 3 || len(username) > 64 {
			return nil, errors.New("username must be between 3 and 64 characters")
		}
		u.Username = username
	}

	// Update password if provided
	if password != "" {
		if len(password) < 6 {
			return nil, errors.New("password must be at least 6 characters")
		}
		hash, err := utils.HashPassword(password)
		if err != nil {
			return nil, err
		}
		u.Password = hash
	}

	// Save changes
	if err := s.repo.Update(u); err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return nil, errors.New("email or username already exists")
		}
		return nil, err
	}

	return u, nil
}

func (s *userService) Delete(id uint) error {
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func (s *userService) GetByID(id string) (*models.User, error) {
	return s.repo.GetByStringID(id)
}

func (s *userService) UpdateStatus(id string, status string) (*models.User, error) {
	user, err := s.repo.GetByStringID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	user.Status = models.UserStatus(status)
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) BatchDelete(ids []string) error {
	return s.repo.BatchDelete(ids)
}
