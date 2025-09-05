package service

import (
	"context"
	"opensesame/internal/etag"
	"opensesame/internal/models/db"
	"opensesame/internal/models/dto"
	"opensesame/internal/repository"
	"time"

	"golang.org/x/sync/errgroup"
)

type StatusService struct {
	configRepo           repository.ConfigRepository
	entryRepo            repository.EntryRepository
	discoveredDeviceRepo repository.DiscoveredDeviceRepository
	// accessRepo repository.AccessRepository
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

func (s *StatusService) WaitForStatus(
	ctx context.Context,
	clientETag uint64,
	timeout time.Duration,
) (*dto.StatusResponse, uint64, bool, error) {
	// if no etag or no timeout, return immediately
	if clientETag == 0 || timeout == 0 {
		status, err := s.GetStatus(ctx)
		if err != nil {
			return nil, etag.Current(), false, err
		}
		return status, etag.Current(), true, nil
	}

	changed := etag.Wait(clientETag, timeout)
	if !changed {
		return nil, etag.Current(), false, nil
	}

	status, err := s.GetStatus(ctx)
	if err != nil {
		return nil, etag.Current(), false, err
	}

	return status, etag.Current(), true, nil
}

func (s *StatusService) GetStatus(ctx context.Context) (*dto.StatusResponse, error) {
	var (
		systemConfig      *db.SystemConfig
		entryDevices      []dto.EntryDevice
		discoveredDevices []dto.DiscoveredDevice
	)

	g, ctx := errgroup.WithContext(ctx)

	// system config
	g.Go(func() error {
		cfg, err := s.configRepo.GetSystemConfig(ctx)
		if err != nil {
			return err
		}
		systemConfig = cfg
		return nil
	})

	// entry devices
	g.Go(func() error {
		devices, err := s.entryRepo.List(ctx)
		if err != nil {
			return err
		}
		summaries := make([]dto.EntryDevice, 0, len(devices))
		for _, d := range devices {
			summaries = append(summaries, dto.EntryDevice{
				ID:          d.EntryID,
				IPAddress:   d.IPAddress,
				MacAddress:  d.MacAddress,
				Port:        d.Port,
				Name:        d.Name,
				Description: d.Description,
				LockStatus:  string(d.LockStatus),
				IsOnline:    d.LastSeen.After(time.Now().Add(-5 * time.Minute)),
				LastSeen:    d.LastSeen.Unix(),
				CreatedAt:   d.CreatedAt.Unix(),
				UpdatedAt:   d.UpdatedAt.Unix(),
			})
		}
		entryDevices = summaries
		return nil
	})

	// discovered devices
	g.Go(func() error {
		discovered, err := s.discoveredDeviceRepo.List(ctx)
		if err != nil {
			return err
		}
		summaries := make([]dto.DiscoveredDevice, 0, len(discovered))
		for _, d := range discovered {
			summaries = append(summaries, dto.DiscoveredDevice{
				ID:         d.ID,
				IPAddress:  d.IPAddress,
				MacAddress: d.MacAddress,
				Instance:   d.Instance,
				DeviceType: d.DeviceType,
				LastSeen:   d.LastSeen.Unix(),
			})
		}
		discoveredDevices = summaries
		return nil
	})

	// wait for goroutines
	if err := g.Wait(); err != nil {
		return nil, err
	}

	// remove discovered devices that already exist as entry devices
	entryMacs := make(map[string]struct{}, len(entryDevices))
	for _, e := range entryDevices {
		entryMacs[e.MacAddress] = struct{}{}
	}

	filteredDiscovered := make([]dto.DiscoveredDevice, 0, len(discoveredDevices))
	for _, d := range discoveredDevices {
		if _, exists := entryMacs[d.MacAddress]; !exists {
			filteredDiscovered = append(filteredDiscovered, d)
		}
	}

	resp := &dto.StatusResponse{
		ETag:              etag.Current(),
		EntryDevices:      entryDevices,
		DiscoveredDevices: filteredDiscovered,
	}
	if systemConfig != nil {
		resp.SystemName = systemConfig.SystemName
	}

	return resp, nil
}
