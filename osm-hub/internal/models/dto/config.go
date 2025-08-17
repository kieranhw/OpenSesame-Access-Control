package dto

type ConfigResponse struct {
	Configured        bool    `json:"configured"`
	SystemName        *string `json:"system_name,omitempty"`
	BackupCode        *string `json:"backup_code,omitempty"`
	SessionTimeoutSec *int    `json:"session_timeout_sec,omitempty"`
}

type CreateConfigRequest struct {
	SystemName        string `json:"system_name"`
	AdminPassword     string `json:"admin_password"`
	SessionTimeoutSec int    `json:"session_timeout_sec"`
}

type UpdateConfigRequest struct {
	SystemName        *string `json:"system_name,omitempty"`
	AdminPassword     *string `json:"admin_password,omitempty"`
	SessionTimeoutSec *int    `json:"session_timeout_sec,omitempty"`
}
