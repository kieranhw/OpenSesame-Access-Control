package dto

type LoginRequest struct {
	Password string `json:"password"`
}

type SessionResponse struct {
	Message       *string `json:"message,omitempty"`
	Authenticated bool    `json:"authenticated"`
	Configured    bool    `json:"configured"`
}

type LogoutResponse struct {
	Success bool `json:"success"`
}
