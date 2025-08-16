package service

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"opensesame/internal/constants"
	"opensesame/internal/models/db"
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
	cfg, err := a.configSvc.GetSystemConfigEntity(ctx)
	if err != nil || cfg == nil {
		return dto.SessionResponse{
			Message:       util.StrPtr(ErrNotConfigured.Error()),
			Authenticated: false,
			Configured:    false,
		}, nil, ErrNotConfigured
	}

	if err := a.validatePassword(cfg.AdminPasswordHash, req.Password); err != nil {
		return dto.SessionResponse{
			Message:       util.StrPtr("invalid credentials"),
			Authenticated: false,
			Configured:    true,
		}, nil, errors.New("invalid credentials")
	}

	cookie, err := a.createCookie(cfg)
	if err != nil {
		return dto.SessionResponse{
			Message:       util.StrPtr("failed to create session"),
			Authenticated: false,
			Configured:    true,
		}, nil, err
	}

	log.Printf("Setting cookie name: '%s'", constants.SessionCookieName)

	return dto.SessionResponse{
		Authenticated: true,
		Configured:    true,
	}, cookie, nil
}

func (a *AuthService) ValidateSession(ctx context.Context, tokenString string) (bool, error) {
	cfg, err := a.configSvc.GetSystemConfigEntity(ctx)
	if err != nil || cfg == nil {
		return false, ErrNotConfigured
	}

	token, err := a.parseToken(cfg, tokenString)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return false, nil
		}
		return false, err
	}

	return token.Valid, nil
}

// refreshes the session by creating a new cookie with a new token based on the system config
func (a *AuthService) RefreshSession(ctx context.Context, cookie *http.Cookie) (*http.Cookie, error) {
	cfg, err := a.configSvc.GetSystemConfigEntity(ctx)
	if err != nil || cfg == nil {
		return nil, ErrNotConfigured
	}

	// Validate the existing token
	token, err := a.parseToken(cfg, cookie.Value)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid session token")
	}

	newCookie, err := a.createCookie(cfg)
	if err != nil {
		return nil, err
	}

	return newCookie, nil
}

// generates a signed JWT and wraps it in an HTTP cookie
func (a *AuthService) createCookie(cfg *db.SystemConfig) (*http.Cookie, error) {
	expiry := time.Now().Add(
		time.Duration(cfg.SessionTimeoutSec) * time.Second,
	).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": expiry,
		"sub": "admin", // TODO: Use client ID
	})

	signed, err := token.SignedString([]byte(cfg.SystemSecret))
	if err != nil {
		return nil, errors.New("token signing failed")
	}

	cookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    signed,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // TODO: set true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(expiry, 0),
	}

	return cookie, nil
}

func (a *AuthService) parseToken(cfg *db.SystemConfig, tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(cfg.SystemSecret), nil
		},
	)
}

func (a *AuthService) validatePassword(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	)
}
