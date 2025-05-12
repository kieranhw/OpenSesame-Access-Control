package service

// ValidatePin validates the pin
func ValidatePin(pin string) bool {
	// Example: Here we would have logic like checking the pin against a database
	if pin == "1234" {
		return true
	}
	return false
}
