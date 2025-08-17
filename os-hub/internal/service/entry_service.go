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

func (s *EntryService) ListEntryDevices(ctx context.Context) ([]dto.EntryDeviceDTO, error) {
	devices, err := s.repo.ListEntryDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing entry devices: %w", err)
	}

	dtos := make([]dto.EntryDeviceDTO, len(devices))
	for i, d := range devices {
		dtos[i] = s.mapEntryDeviceToDTO(d)
	}

	return dtos, nil
}

func (s *EntryService) mapEntryDeviceToDTO(model *db.EntryDevice) dto.EntryDeviceDTO {
	cmds := make([]dto.EntryCommandDTO, len(model.Commands))
	for i, c := range model.Commands {
		cmds[i] = s.mapEntryCommandToDTO(&c)
	}

	return dto.EntryDeviceDTO{
		ID:          model.EntryID,
		Name:        model.Name,
		IP:          model.IP,
		Port:        model.Port,
		Description: model.Description,
		Protocol:    string(model.Protocol),
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		Commands:    cmds,
	}
}

func (s *EntryService) mapEntryCommandToDTO(model *db.EntryCommand) dto.EntryCommandDTO {
	var httpDTO *dto.HttpCommandDTO
	if model.HttpCommand != nil {
		httpDTO = &dto.HttpCommandDTO{
			URL:    model.HttpCommand.URL,
			Method: model.HttpCommand.Method,
			Body:   model.HttpCommand.Body,
		}

		if model.HttpCommand.Headers != "" {
			httpDTO.Headers = parseHeaders(model.HttpCommand.Headers)
		}
	}

	return dto.EntryCommandDTO{
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
