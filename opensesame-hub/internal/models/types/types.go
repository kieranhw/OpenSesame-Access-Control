package types

import (
	"opensesame/internal/repository"
	"opensesame/internal/service"
)

type Repositories struct {
	Config repository.ConfigRepository
	// Add other repositories as needed:
	// Session repository.SessionRepository
	// Entry   repository.EntryRepository
}

type Services struct {
	Config *service.ConfigService
	Auth   *service.AuthService
	// Add other services as needed
}
