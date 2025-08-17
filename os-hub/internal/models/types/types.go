package types

import (
	"opensesame/internal/repository"
	"opensesame/internal/service"
)

type Repositories struct {
	Config repository.ConfigRepository
	Entry  repository.EntryRepository
}

type Services struct {
	Config *service.ConfigService
	Auth   *service.AuthService
	Entry  *service.EntryService
}
