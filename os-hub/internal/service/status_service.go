package service

import (
	"context"
	"opensesame/internal/models/dto"
	"opensesame/internal/repository"
)

type StatusService struct {
	configRepo           repository.ConfigRepository
	entryRepo            repository.EntryRepository
	discoveredDeviceRepo repository.DiscoveredDeviceRepository
	// accessRepo repository.AccessRepository // if/when you add it
}

func NewStatusService(
	configRepo repository.ConfigRepository,
	entryRepo repository.EntryRepository,
	discoveredRepo repository.DiscoveredDeviceRepository,
) *StatusService {
	return &StatusService{
		configRepo:           configRepo,
		entryRepo:            entryRepo,
		discoveredDeviceRepo: discoveredRepo,
	}
}

func (s *StatusService) GetStatus(ctx context.Context) (*dto.StatusResponse, error) {
	// System config
	cfg, err := s.configRepo.GetSystemConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Entry devices
	entryDevices, err := s.entryRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	entrySummaries := make([]dto.EntryStatus, 0, len(entryDevices))
	for _, e := range entryDevices {
		entrySummaries = append(entrySummaries, dto.EntryStatus{
			ID:        e.EntryID,
			IPAddress: e.IP,
		})
	}

	// Discovered devices
	discovered, err := s.discoveredDeviceRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	discoveredSummaries := make([]dto.DiscoveryStatus, 0, len(discovered))
	for _, d := range discovered {
		discoveredSummaries = append(discoveredSummaries, dto.DiscoveryStatus{
			ID:         d.ID,
			MacAddress: d.MacAddress,
			Instance:   d.Instance,
			IPAddress:  d.IPv4,
			DeviceType: d.DeviceType,
		})
	}

	resp := &dto.StatusResponse{
		Configured:        cfg != nil,
		EntryDevices:      entrySummaries,
		DiscoveredDevices: discoveredSummaries,
		// AccessDevices:     accessSummaries,
	}
	if cfg != nil {
		resp.SystemName = &cfg.SystemName
	}

	return resp, nil
}
