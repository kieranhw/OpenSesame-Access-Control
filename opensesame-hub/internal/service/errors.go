package service

import "errors"

var ErrNotConfigured = errors.New("system not configured")
var ErrAlreadyConfigured = errors.New("system already configured")
