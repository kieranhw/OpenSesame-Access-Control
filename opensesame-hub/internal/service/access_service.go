package service

func ValidatePin(pin string) bool {
	if pin == "1234" {
		return true
	}
	return false
}
