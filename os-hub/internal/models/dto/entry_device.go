package dto

import "time"

type EntryDeviceDTO struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	IP          string            `json:"ip"`
	Port        int               `json:"port"`
	Description string            `json:"description,omitempty"`
	Protocol    string            `json:"protocol"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Commands    []EntryCommandDTO `json:"commands,omitempty"`
}
