package dto

type EntryDevice struct {
	ID         uint   `json:"id"`
	MacAddress string `json:"mac_address"`
	IPAddress  string `json:"ip_address"`
	Port       int    `json:"port"`

	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`

	LockStatus string `json:"lock_status"`
	IsOnline   bool   `json:"is_online"`

	LastSeen  int64 `json:"last_seen"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`

	Commands []EntryCommand `json:"commands,omitempty"`
}

type CreateEntryDeviceRequest struct {
	MacAddress string `json:"mac_address"`
	IPAddress  string `json:"ip_address"`
	Port       int    `json:"port"`

	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type UpdateEntryDeviceRequest struct {
	MacAddress *string `json:"mac_address,omitempty"`
	IPAddress  *string `json:"ip_address,omitempty"`
	Port       *int    `json:"port,omitempty"`

	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`

	LastSeen *int64 `json:"last_seen,omitempty"`
}
