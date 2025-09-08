package util

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"opensesame/internal/models/types"
)

func GetDeviceInfo(ctx context.Context, ip string, port int) (*types.DeviceInfo, error) {
	// Validate IP
	if net.ParseIP(ip) == nil {
		return nil, fmt.Errorf("%w: invalid IP address %q", types.ErrBadRequest, ip)
	}

	// Validate port
	if port <= 0 || port > 65535 {
		return nil, fmt.Errorf("%w: invalid port %d", types.ErrBadRequest, port)
	}

	url := fmt.Sprintf("http://%s:%d/info", ip, port)

	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: building request: %v", types.ErrBadRequest, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", types.ErrUnreachableDevice, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: unexpected status %d", types.ErrUnreachableDevice, resp.StatusCode)
	}

	var info types.DeviceInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("%w: decode error: %v", types.ErrUnreachableDevice, err)
	}

	return &info, nil
}
