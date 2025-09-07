package models

import "opensesame/internal/service"

type Services struct {
	Config    *service.ConfigService
	Auth      *service.AuthService
	Entry     *service.EntryService
	Discovery *service.DiscoveryService
	Status    *service.StatusService
}
