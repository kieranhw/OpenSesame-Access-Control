package models

import (
	"opensesame/internal/repository"
)

type Repositories struct {
	Config           repository.ConfigRepository
	Entry            repository.EntryRepository
	DiscoveredDevice repository.DiscoveredDeviceRepository
}
