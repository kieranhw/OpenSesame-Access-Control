package types

import (
	"opensesame/internal/repository"
	"opensesame/internal/service"
)

type Repositories struct {
	Config repository.ConfigRepository
}

type Services struct {
	Config *service.ConfigService
	Auth   *service.AuthService
}
