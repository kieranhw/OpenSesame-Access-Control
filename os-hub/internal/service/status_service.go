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

	// otherwise, do the long-poll wait
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
		cfg                 *db.SystemConfig
		entrySummaries      []dto.EntryDevice
		discoveredSummaries []dto.DiscoveryStatus
	)

	g, ctx := errgroup.WithContext(ctx)

	// system config
	g.Go(func() error {
		c, err := s.configRepo.GetSystemConfig(ctx)
		if err != nil {
			return err
		}
		cfg = c
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
				IP:          d.IP,
				Port:        d.Port,
				Name:        d.Name,
				Description: d.Description,
				CreatedAt:   d.CreatedAt,
				UpdatedAt:   d.UpdatedAt,
			})
		}
		entrySummaries = summaries
		return nil
	})

	// discovered devices
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
		ETag:              etag.Current(),
		EntryDevices:      entrySummaries,
		DiscoveredDevices: discoveredSummaries,
	}
	if cfg != nil {
		resp.SystemName = cfg.SystemName
	}

	return resp, nil
}
