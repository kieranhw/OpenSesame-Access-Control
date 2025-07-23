package service

import (
	"context"
	"errors"
	"time"

	"opensesame/internal/models/db"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (a *AuthService) CreateSession(ctx context.Context) (*db.Session, error) {
	token := uuid.NewString()
	sess := &db.Session{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := a.db.WithContext(ctx).Create(sess).Error; err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *AuthService) ValidateSession(ctx context.Context, tokenString string, systemSecret string) (bool, error) {
	return s.validateJWT(tokenString, systemSecret)
}

func (s *AuthService) validateJWT(tokenString string, systemSecret string) (bool, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) {
			// ensure the signing method is HMAC
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(systemSecret), nil
		},
	)

	if err != nil {
		// return false instead of error if token is expired
		if errors.Is(err, jwt.ErrTokenExpired) {
			return false, nil
		}
		return false, err
	}

	if !token.Valid {
		return false, nil
	}

	return true, nil
}

func (a *AuthService) DeleteSession(ctx context.Context, token string) error {
	return a.db.WithContext(ctx).Delete(&db.Session{}, "token = ?", token).Error
}
