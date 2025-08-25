package dto

import "time"

type EntryDevice struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	MacAddress  string  `json:"mac_address"`
	IPAddress   string  `json:"ip_address"`
	Port        int     `json:"port"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Commands  []EntryCommand `json:"commands,omitempty"`
}

type CreateEntryDeviceRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	MacAddress  string  `json:"mac_address"`
	IPAddress   string  `json:"ip_address"`
	Port        int     `json:"port"`
}

type UpdateEntryDeviceRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	MacAddress  *string `json:"mac_address,omitempty"`
	IPAddress   *string `json:"ip_address,omitempty"`
	Port        *int    `json:"port,omitempty"`
}
