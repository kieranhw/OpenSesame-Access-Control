package dto

type StatusResponse struct {
	SystemName        *string           `json:"system_name,omitempty"`
	Configured        bool              `json:"configured"`
	EntryDevices      []EntryStatus     `json:"entry_devices"`
	DiscoveredDevices []DiscoveryStatus `json:"discovered_devices"`
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
}
