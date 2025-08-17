package service

import "errors"

var ErrNotConfigured = errors.New("system not configured")
var ErrAlreadyConfigured = errors.New("system already configured")
var ErrNoUpdateFields = errors.New("no fields provided for update")
var ErrPasswordHashingFailed = errors.New("failed to hash password")
