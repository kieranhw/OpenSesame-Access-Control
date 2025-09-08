package types

type DeviceInfo struct {
	MacAddress   string `json:"mac_address"`
	DeviceType   string `json:"device_type"`   // e.g. "entry"/"access"
	InstanceName string `json:"instance_name"` // e.g. "OpenSesame Relay Lock"
	InstanceType string `json:"instance_type"` // e.g. "relay_lock"
}
