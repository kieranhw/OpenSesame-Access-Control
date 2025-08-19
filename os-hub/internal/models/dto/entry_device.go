package dto

import "time"

type EntryDevice struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	IP          string         `json:"ip"`
	Port        int            `json:"port"`
	Description string         `json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Commands    []EntryCommand `json:"commands,omitempty"`
}

type CreateEntryDeviceRequest struct {
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	Description string `json:"description,omitempty"`
}

type UpdateEntryDeviceRequest struct {
	Name        *string `json:"name,omitempty"`
	IP          *string `json:"ip,omitempty"`
	Port        *int    `json:"port,omitempty"`
	Description *string `json:"description,omitempty"`
}
