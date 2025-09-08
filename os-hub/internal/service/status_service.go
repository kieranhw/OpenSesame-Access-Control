package service

import (
	"context"
	"opensesame/internal/etag"
	"opensesame/internal/models/db"
	"opensesame/internal/models/dto"
	"opensesame/internal/models/mappers"
	"opensesame/internal/repository"
	"time"

	"golang.org/x/sync/errgroup"
)

type StatusService struct {
	configRepo repository.ConfigRepository
	deviceRepo repository.DeviceRepository
}

func NewStatusService(
	configRepo repository.ConfigRepository,
	deviceRepo repository.DeviceRepository,
) *StatusService {
	return &StatusService{
		configRepo: configRepo,
		deviceRepo: deviceRepo,
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
		systemConfig *db.SystemConfig
		entryDevices []dto.EntryDevice
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
		devices, err := s.deviceRepo.ListEntryDevices(ctx)
		if err != nil {
			return err
		}

		entryDevices = make([]dto.EntryDevice, 0, len(devices))
		for _, d := range devices {
			entryDevices = append(entryDevices, mappers.EntryDeviceToDTO(d))
		}

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

	resp := &dto.StatusResponse{
		ETag:         etag.Current(),
		EntryDevices: entryDevices,
	}
	if systemConfig != nil {
		resp.SystemName = systemConfig.SystemName
	}

	return resp, nil
}
