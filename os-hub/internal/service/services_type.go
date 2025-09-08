package service

type ServicesType struct {
	Config    *ConfigService
	Auth      *AuthService
	Entry     *EntryService
	Discovery *DiscoveryService
	Status    *StatusService
}
