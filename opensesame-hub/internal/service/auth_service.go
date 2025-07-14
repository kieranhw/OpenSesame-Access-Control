package service

import (
	"context"
	"time"

	"opensesame/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (a *AuthService) CreateSession(ctx context.Context) (*model.Session, error) {
	token := uuid.NewString()
	sess := &model.Session{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := a.db.WithContext(ctx).Create(sess).Error; err != nil {
		return nil, err
	}
	return sess, nil
}

func (a *AuthService) ValidateSession(ctx context.Context, token string) (bool, error) {
	var sess model.Session
	err := a.db.WithContext(ctx).
		First(&sess, "token = ?", token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	if time.Now().After(sess.ExpiresAt) {
		return false, nil
	}
	return true, nil
}

func (a *AuthService) DeleteSession(ctx context.Context, token string) error {
	return a.db.WithContext(ctx).Delete(&model.Session{}, "token = ?", token).Error
}
