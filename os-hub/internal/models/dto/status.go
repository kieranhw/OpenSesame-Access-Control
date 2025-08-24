package dto

type StatusResponse struct {
	SystemName        *string           `json:"system_name,omitempty"`
	EntryDevices      []EntryStatus     `json:"entry_devices"`
	DiscoveredDevices []DiscoveryStatus `json:"discovery"`
	// AccessDevices     []AccessStatus    `json:"access_devices"`
}

type EntryStatus struct {
	ID        uint   `json:"id"`
	IPAddress string `json:"ip_address"`
}

type DiscoveryStatus struct {
	ID         uint   `json:"id"`
	IPAddress  string `json:"ip_address"`
	MacAddress string `json:"mac_address"`
	Instance   string `json:"instance"`
	DeviceType string `json:"type"`
	LastSeen   int64  `json:"last_seen"`
}
