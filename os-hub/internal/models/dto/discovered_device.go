package dto

type DiscoveredDevice struct {
	ID         uint   `json:"id"`
	MacAddress string `json:"mac_address"`
	IPAddress  string `json:"ip_address"`
	Port       int    `json:"port"`

	Instance   string `json:"instance"`
	DeviceType string `json:"type"`

	LastSeen int64 `json:"last_seen"`
}
