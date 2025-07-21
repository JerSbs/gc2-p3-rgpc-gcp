package app

import (
	"context"
	"errors"
	"time"

	"auth-service/internal/auth/domain"

	"github.com/google/uuid"
)

type AuthService struct {
	repo   domain.UserRepository
	jwt    *JWTManager
	hasher PasswordHasher
}

// NewAuthService buat service auth dengan dependency injection
func NewAuthService(repo domain.UserRepository, jwt *JWTManager, hasher PasswordHasher) *AuthService {
	return &AuthService{
		repo:   repo,
		jwt:    jwt,
		hasher: hasher,
	}
}

// RegisterUser untuk register user baru
func (s *AuthService) RegisterUser(ctx context.Context, user domain.User) (domain.User, error) {
	existing, _ := s.repo.GetByEmail(user.Email)
	if existing.Email != "" {
		return domain.User{}, errors.New("email already registered")
	}

	hashedPassword, err := s.hasher.Hash(user.Password)
	if err != nil {
		return domain.User{}, err
	}

	user.ID = uuid.NewString()
	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err = s.repo.Save(user)
	if err != nil {
		return domain.User{}, err
	}

	user.Password = ""
	return user, nil
}

// LoginUser untuk login user dan return token
func (s *AuthService) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}

	err = s.hasher.Compare(user.Password, password)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return s.jwt.Generate(user)
}

// ValidateToken untuk ambil data user dari token
func (s *AuthService) ValidateToken(ctx context.Context, token string) (domain.User, error) {
	return s.jwt.Verify(token)
}
