package types

type DeviceType string

const (
	DeviceTypeEntry  DeviceType = "entry"
	DeviceTypeAccess DeviceType = "access"
)

type InstanceType string

const (
	InstanceTypeRelayLock InstanceType = "relay_lock"
)
