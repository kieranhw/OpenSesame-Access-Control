package dto

type DiscoveredDevice struct {
	ID         uint   `json:"id"`
	IPAddress  string `json:"ip_address"`
	MacAddress string `json:"mac_address"`
	Instance   string `json:"instance"`
	DeviceType string `json:"type"`
	LastSeen   int64  `json:"last_seen"`
}
