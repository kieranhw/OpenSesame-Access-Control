package service

import (
	"context"
	"fmt"
	"strings"

	"opensesame/internal/models/db"
	"opensesame/internal/models/dto"
	"opensesame/internal/models/mappers"
	"opensesame/internal/models/types"
	"opensesame/internal/repository"
	"opensesame/internal/util"
)

type EntryService struct {
	repo repository.DeviceRepository
}

func NewEntryService(repo repository.DeviceRepository) *EntryService {
	return &EntryService{repo: repo}
}

func (s *EntryService) ListEntryDevices(ctx context.Context) ([]dto.EntryDevice, error) {
	devices, err := s.repo.ListEntryDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing entry devices: %w", err)
	}

	dtos := make([]dto.EntryDevice, len(devices))
	for i, d := range devices {
		dtos[i] = mappers.EntryDeviceToDTO(d)
	}

	return dtos, nil
}

func (s *EntryService) GetEntryDeviceByMac(ctx context.Context, mac string) (*dto.EntryDevice, error) {
	device, err := s.repo.GetEntryDeviceByMac(ctx, mac)
	if err != nil {
		return nil, fmt.Errorf("getting entry device by mac: %w", err)
	}
	if device == nil {
		return nil, nil
	}

	dto := mappers.EntryDeviceToDTO(device)
	return &dto, nil
}

func (s *EntryService) CreateEntryDevice(ctx context.Context, req dto.CreateEntryDeviceRequest) (dto.EntryDevice, error) {
	info, err := util.GetDeviceInfo(ctx, req.IPAddress, req.Port)
	if err != nil {
		return dto.EntryDevice{}, err
	}

	if info.MacAddress == "" {
		return dto.EntryDevice{}, fmt.Errorf("device at %s:%d did not return a MAC address", req.IPAddress, req.Port)
	}

	if info.DeviceType != "entry" {
		return dto.EntryDevice{}, fmt.Errorf("device at %s:%d is not an entry device (type=%s)", req.IPAddress, req.Port, info.DeviceType)
	}

	mac := strings.ToUpper(strings.ReplaceAll(info.MacAddress, ":", ""))

	device := db.Device{
		MacAddress:   mac,
		IPAddress:    req.IPAddress,
		Port:         req.Port,
		Name:         req.Name,
		Description:  req.Description,
		ServiceType:  req.ServiceType,
		DeviceType:   string(info.DeviceType),
		InstanceType: info.InstanceType,
		InstanceName: info.InstanceName,
	}

	entry := &db.EntryDevice{
		Device:     device,
		LockStatus: types.LockStatusUnknown,
	}

	if req.Description != nil && *req.Description != "" {
		entry.Device.Description = req.Description
	}
	if req.ServiceType != nil && *req.ServiceType != "" {
		entry.Device.ServiceType = req.ServiceType
	}

	if err := s.repo.CreateEntryDevice(ctx, entry); err != nil {
		return dto.EntryDevice{}, fmt.Errorf("creating entry device: %w", err)
	}

	return mappers.EntryDeviceToDTO(entry), nil
}

func (s *EntryService) UpdateEntryDevice(ctx context.Context, id uint, req dto.UpdateEntryDeviceRequest) (dto.EntryDevice, error) {
	existing, err := s.repo.GetEntryDeviceById(ctx, id)
	if err != nil {
		return dto.EntryDevice{}, err
	}
	if existing == nil {
		return dto.EntryDevice{}, fmt.Errorf("entry device %d not found", id)
	}

	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = req.Description
	}
	if req.MacAddress != nil {
		mac := strings.ToUpper(strings.ReplaceAll(*req.MacAddress, ":", ""))
		updates["mac_address"] = mac
	}
	if req.IPAddress != nil {
		updates["ip_address"] = *req.IPAddress
	}
	if req.Port != nil {
		updates["port"] = *req.Port
	}
	if req.ServiceType != nil {
		updates["service_type"] = req.ServiceType
	}
	if req.DeviceType != nil {
		updates["device_type"] = *req.DeviceType
	}
	if req.InstanceType != nil {
		updates["instance_type"] = *req.InstanceType
	}
	if req.InstanceName != nil {
		updates["instance_name"] = *req.InstanceName
	}
	if req.LastSeen != nil {
		updates["last_seen"] = *req.LastSeen
	}

	if len(updates) > 0 {
		if err := s.repo.UpdateEntryDevice(ctx, id, updates); err != nil {
			return dto.EntryDevice{}, err
		}
	}

	updated, err := s.repo.GetEntryDeviceById(ctx, id)
	if err != nil {
		return dto.EntryDevice{}, err
	}

	return mappers.EntryDeviceToDTO(updated), nil
}
