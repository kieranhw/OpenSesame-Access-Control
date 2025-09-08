package types

import "errors"

var ErrNotConfigured = errors.New("system not configured")
var ErrAlreadyConfigured = errors.New("system already configured")
var ErrNoUpdateFields = errors.New("no fields provided for update")
var ErrPasswordHashingFailed = errors.New("failed to hash password")
var ErrBadRequest = errors.New("bad request")
var ErrUnreachableDevice = errors.New("unable to reach device at specified IP/port")
