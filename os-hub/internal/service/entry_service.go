package service

import (
	"context"
	"encoding/json"
	"fmt"

	"opensesame/internal/models/db"
	"opensesame/internal/models/dto"
	"opensesame/internal/repository"
)

type EntryService struct {
	repo repository.EntryRepository
}

func NewEntryService(repo repository.EntryRepository) *EntryService {
	return &EntryService{repo: repo}
}

func (s *EntryService) ListEntryDevices(ctx context.Context) ([]dto.EntryDevice, error) {
	devices, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing entry devices: %w", err)
	}

	dtos := make([]dto.EntryDevice, len(devices))
	for i, d := range devices {
		dtos[i] = s.mapEntryDeviceToDTO(d)
	}

	return dtos, nil
}

func (s *EntryService) CreateEntryDevice(ctx context.Context, req dto.CreateEntryDeviceRequest) (dto.EntryDevice, error) {
	model := &db.EntryDevice{
		Name:        req.Name,
		IP:          req.IP,
		Port:        req.Port,
		Description: req.Description,
	}

	if err := s.repo.CreateEntryDevice(ctx, model); err != nil {
		return dto.EntryDevice{}, fmt.Errorf("creating entry device: %w", err)
	}

	return s.mapEntryDeviceToDTO(model), nil
}

func (s *EntryService) mapEntryDeviceToDTO(model *db.EntryDevice) dto.EntryDevice {
	cmds := make([]dto.EntryCommand, len(model.Commands))
	for i, c := range model.Commands {
		cmds[i] = s.mapEntryCommandToDTO(&c)
	}

	return dto.EntryDevice{
		ID:          model.EntryID,
		Name:        model.Name,
		IP:          model.IP,
		Port:        model.Port,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		Commands:    cmds,
	}
}

func (s *EntryService) mapEntryCommandToDTO(model *db.EntryCommand) dto.EntryCommand {
	// Build HTTP DTO directly from EntryCommand fields
	httpDTO := &dto.HttpCommand{
		URL:    model.URL,
		Method: model.Method,
		Body:   model.Body,
	}

	if model.Headers != "" {
		httpDTO.Headers = parseHeaders(model.Headers)
	}

	return dto.EntryCommand{
		ID:        model.CommandID,
		Type:      string(model.CommandType),
		Status:    string(model.Status),
		CreatedAt: model.CreatedAt,
		Http:      httpDTO,
	}
}

func parseHeaders(raw string) map[string]string {
	if raw == "" {
		return nil
	}
	var headers map[string]string
	if err := json.Unmarshal([]byte(raw), &headers); err != nil {
		// If parsing fails, return nil instead of crashing
		return nil
	}
	return headers
}
