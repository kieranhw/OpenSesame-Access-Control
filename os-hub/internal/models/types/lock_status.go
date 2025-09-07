package types

type LockStatus string

const (
	LockStatusUnknown  LockStatus = "UNKNOWN"
	LockStatusLocked   LockStatus = "LOCKED"
	LockStatusUnlocked LockStatus = "UNLOCKED"
)
