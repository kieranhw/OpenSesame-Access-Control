package types

import (
	"opensesame/internal/repository"
	"opensesame/internal/service"
)

type Repositories struct {
	Config           repository.ConfigRepository
	Entry            repository.EntryRepository
	DiscoveredDevice repository.DiscoveredDeviceRepository
}

type Services struct {
	Config    *service.ConfigService
	Auth      *service.AuthService
	Entry     *service.EntryService
	Discovery *service.DiscoveryService
	Status    *service.StatusService
}
