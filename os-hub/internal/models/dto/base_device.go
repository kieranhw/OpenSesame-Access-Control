package dto

type BaseDevice struct {
	ID           uint    `json:"id"`
	MacAddress   string  `json:"mac_address"`
	IPAddress    string  `json:"ip_address"`
	Port         int     `json:"port"`
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	DeviceType   string  `json:"device_type"`
	InstanceType string  `json:"instance_type"`
	InstanceName string  `json:"instance_name"`
	IsOnline     bool    `json:"is_online"`
	LastSeen     int64   `json:"last_seen"`
	CreatedAt    int64   `json:"created_at"`
	UpdatedAt    int64   `json:"updated_at"`
}
