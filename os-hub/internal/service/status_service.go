package service

import (
	"context"
	"opensesame/internal/models/db"
	"opensesame/internal/models/dto"
	"opensesame/internal/repository"

	"golang.org/x/sync/errgroup"
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
	var (
		cfg                 *db.SystemConfig
		entrySummaries      []dto.EntryDevice
		discoveredSummaries []dto.DiscoveryStatus
	)

	g, ctx := errgroup.WithContext(ctx)

	// System config
	g.Go(func() error {
		c, err := s.configRepo.GetSystemConfig(ctx)
		if err != nil {
			return err
		}
		cfg = c
		return nil
	})

	// Entry devices
	g.Go(func() error {
		devices, err := s.entryRepo.List(ctx)
		if err != nil {
			return err
		}
		summaries := make([]dto.EntryDevice, 0, len(devices))
		for _, d := range devices {
			summaries = append(summaries, dto.EntryDevice{
				ID:          d.EntryID,
				IP:          d.IP,
				Port:        d.Port,
				Name:        d.Name,
				Description: d.Description,
				CreatedAt:   d.CreatedAt,
			})
		}
		entrySummaries = summaries
		return nil
	})

	// Discovered devices
	g.Go(func() error {
		discovered, err := s.discoveredDeviceRepo.List(ctx)
		if err != nil {
			return err
		}
		summaries := make([]dto.DiscoveryStatus, 0, len(discovered))
		for _, d := range discovered {
			summaries = append(summaries, dto.DiscoveryStatus{
				ID:         d.ID,
				MacAddress: d.MacAddress,
				Instance:   d.Instance,
				IPAddress:  d.IPv4,
				DeviceType: d.DeviceType,
				LastSeen:   d.LastSeen.Unix(),
			})
		}
		discoveredSummaries = summaries
		return nil
	})

	// wait for goroutines
	if err := g.Wait(); err != nil {
		return nil, err
	}

	resp := &dto.StatusResponse{
		EntryDevices:      entrySummaries,
		DiscoveredDevices: discoveredSummaries,
	}
	if cfg != nil {
		resp.SystemName = &cfg.SystemName
	}

	return resp, nil
}
