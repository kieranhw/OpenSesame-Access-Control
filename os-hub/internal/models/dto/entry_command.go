package dto

import "time"

type EntryCommandDTO struct {
	ID        uint            `json:"id"`
	Type      string          `json:"type"`
	Status    string          `json:"status,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	Http      *HttpCommandDTO `json:"http,omitempty"`
	// Udp      *UdpCommandDTO  `json:"udp,omitempty"` // future extension
}

type HttpCommandDTO struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}
