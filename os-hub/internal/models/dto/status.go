package dto

type StatusResponse struct {
	ETag              uint64             `json:"etag"`
	SystemName        string             `json:"system_name"`
	EntryDevices      []EntryDevice      `json:"entry_devices"`
	DiscoveredDevices []DiscoveredDevice `json:"discovered_devices"`
	// AccessDevices     []AccessStatus    `json:"access_devices"`
}
