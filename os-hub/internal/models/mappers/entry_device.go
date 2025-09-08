package mappers

import (
	"time"

	"opensesame/internal/constants"
	"opensesame/internal/models/db"
	"opensesame/internal/models/dto"
)

func EntryDeviceToDTO(model *db.EntryDevice) dto.EntryDevice {
	// determine online status based on last seen timestamp
	cutoff := time.Now().Add(-time.Duration(constants.LastSeenThresholdSec) * time.Second)
	isOnline := model.Device.LastSeen.After(cutoff)

	return dto.EntryDevice{
		BaseDevice: dto.BaseDevice{
			ID:           model.DeviceID,
			MacAddress:   model.Device.MacAddress,
			Name:         model.Device.Name,
			IPAddress:    model.Device.IPAddress,
			Port:         model.Device.Port,
			Description:  model.Device.Description,
			DeviceType:   model.Device.DeviceType,
			InstanceType: model.Device.InstanceType,
			InstanceName: model.Device.InstanceName,
			IsOnline:     isOnline,
			LastSeen:     model.Device.LastSeen.Unix(),
			CreatedAt:    model.Device.CreatedAt.Unix(),
			UpdatedAt:    model.Device.UpdatedAt.Unix(),
		},
		LockStatus: model.LockStatus,
		// Commands:   cmds,
	}
}
