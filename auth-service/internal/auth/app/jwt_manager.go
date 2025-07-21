package app

import (
	"errors"
	"os"
	"time"

	"auth-service/internal/auth/domain"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// Buat JWTManager baru dengan secret dari env dan durasi 1 jam
func NewJWTManager() *JWTManager {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "mysecretkey123" // fallback default
	}
	return &JWTManager{
		secretKey:     secret,
		tokenDuration: time.Hour,
	}
}

// Generate token JWT dari user
func (j *JWTManager) Generate(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(j.tokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// Validasi dan ambil data user dari token
func (j *JWTManager) Verify(tokenString string) (domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi method yang digunakan
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil || !token.Valid {
		return domain.User{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return domain.User{}, errors.New("invalid claims")
	}

	// Ambil informasi user dari claim
	user := domain.User{
		ID:    claims["id"].(string),
		Name:  claims["name"].(string),
		Email: claims["email"].(string),
	}

	return user, nil
}
