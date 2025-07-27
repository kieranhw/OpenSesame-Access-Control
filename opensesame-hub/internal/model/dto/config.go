package dto

type CreateConfigPayload struct {
	SystemName    string `json:"system_name"`
	AdminPassword string `json:"admin_password"`
}

type ConfigResponse struct {
	Configured bool    `json:"configured"`
	SystemName *string `json:"system_name,omitempty"`
	BackupCode *string `json:"backup_code,omitempty"`
}

type UpdateConfigPayload struct {
	SystemName    *string `json:"system_name,omitempty"`
	AdminPassword *string `json:"admin_password,omitempty"`
}
