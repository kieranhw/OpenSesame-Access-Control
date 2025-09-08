package dto

import "opensesame/internal/models/types"

type EntryDevice struct {
	BaseDevice
	LockStatus types.LockStatus `json:"lock_status"`
	Commands   []EntryCommand   `json:"commands,omitempty"`
}

type CreateEntryDeviceRequest struct {
	MacAddress   string  `json:"mac_address"`
	IPAddress    string  `json:"ip_address"`
	Port         int     `json:"port"`
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	ServiceType  *string `json:"service_type,omitempty"`
	DeviceType   string  `json:"device_type"`
	InstanceType string  `json:"instance_type"`
	InstanceName string  `json:"instance_name"`
}

type UpdateEntryDeviceRequest struct {
	MacAddress   *string `json:"mac_address,omitempty"`
	IPAddress    *string `json:"ip_address,omitempty"`
	Port         *int    `json:"port,omitempty"`
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	ServiceType  *string `json:"service_type,omitempty"`
	DeviceType   *string `json:"device_type,omitempty"`
	InstanceType *string `json:"instance_type,omitempty"`
	InstanceName *string `json:"instance_name,omitempty"`
	LastSeen     *int64  `json:"last_seen,omitempty"`
}
