// service/auth_service.go
package service

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"opensesame/internal/constants"
	"opensesame/internal/models/dto"
	"opensesame/internal/util"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	configSvc *ConfigService
}

func NewAuthService(configSvc *ConfigService) *AuthService {
	return &AuthService{
		configSvc: configSvc,
	}
}

func (a *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.SessionResponse, *http.Cookie, error) {
	// Get the full config entity (with sensitive fields) through ConfigService
	cfg, err := a.configSvc.GetSystemConfigEntity(ctx)
	if err != nil {
		return dto.SessionResponse{
			Message:       util.StrPtr(ErrNotConfigured.Error()),
			Authenticated: false,
			Configured:    false,
		}, nil, ErrNotConfigured
	}
	if cfg == nil {
		return dto.SessionResponse{
			Message:       util.StrPtr(ErrNotConfigured.Error()),
			Authenticated: false,
			Configured:    false,
		}, nil, ErrNotConfigured
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(cfg.AdminPasswordHash),
		[]byte(req.Password),
	); err != nil {
		return dto.SessionResponse{
			Message:       util.StrPtr("invalid credentials"),
			Authenticated: false,
			Configured:    true,
		}, nil, errors.New("invalid credentials")
	}

	expiry := time.Now().Add(time.Duration(cfg.SessionTimeoutSec) * time.Second).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": expiry,
		"sub": "admin",
	})

	signed, err := token.SignedString([]byte(cfg.SystemSecret))
	if err != nil {
		msg := "could not sign token"
		return dto.SessionResponse{
			Message:       &msg,
			Authenticated: false,
			Configured:    true,
		}, nil, errors.New("token signing failed")
	}

	cookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    signed,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(expiry, 0),
	}

	log.Printf("Setting cookie name: '%s'", constants.SessionCookieName)

	return dto.SessionResponse{
		Authenticated: true,
		Configured:    true,
	}, cookie, nil
}

func (a *AuthService) ValidateSession(ctx context.Context, tokenString string) (bool, error) {
	cfg, err := a.configSvc.GetSystemConfigEntity(ctx)
	if err != nil {
		return false, ErrNotConfigured
	}
	if cfg == nil {
		return false, ErrNotConfigured
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(cfg.SystemSecret), nil
		},
	)

	if err != nil {
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
