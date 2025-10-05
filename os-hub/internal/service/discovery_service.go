package service

import (
	"context"
	"log"
	"strings"
	"time"

	"opensesame/internal/models/dto"
	"opensesame/internal/repository"
	"opensesame/internal/util"

	"github.com/grandcat/zeroconf"
)

type DiscoveryService struct {
	resolver *zeroconf.Resolver
	entries  chan *zeroconf.ServiceEntry
	cancel   context.CancelFunc
	repo     repository.DeviceRepository
	entrySvc *EntryService
}

func NewDiscoveryService(
	repo repository.DeviceRepository,
	entrySvc *EntryService,
) *DiscoveryService {
	return &DiscoveryService{
		entries:  make(chan *zeroconf.ServiceEntry),
		repo:     repo,
		entrySvc: entrySvc,
	}
}

func (d *DiscoveryService) Start(ctx context.Context) error {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return err
	}
	d.resolver = resolver

	ctx, cancel := context.WithCancel(ctx)
	d.cancel = cancel
	go d.run(ctx)
	return nil
}

func (d *DiscoveryService) run(ctx context.Context) {
	// browse for http services with mDNS
	err := d.resolver.Browse(ctx, "_http._tcp", "local.", d.entries)
	if err != nil {
		log.Printf("mDNS browse failed: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Discovery service stopped")
			return
		case entry := <-d.entries:
			if entry == nil || len(entry.AddrIPv4) == 0 {
				continue
			}
			// handle discovery on a separate goroutine
			go d.handleDiscoveredDevice(ctx, entry)
		}
	}
}

func (d *DiscoveryService) handleDiscoveredDevice(ctx context.Context, entry *zeroconf.ServiceEntry) {
	if entry.HostName == "" || !strings.Contains(entry.HostName, "opensesame-device") {
		return
	}
	if len(entry.AddrIPv4) == 0 {
		return
	}

	ip := entry.AddrIPv4[0].String()

	info, err := util.GetDeviceInfo(ctx, ip, entry.Port)
	if err != nil {
		log.Printf("failed to get device info from %s:%d: %v", ip, entry.Port, err)
		return
	} else {
		log.Printf("[DISCOVERY] Found device: %s (%s) at %s:%d", info.InstanceName, info.DeviceType, ip, entry.Port)
	}

	log.Printf("[DISCOVERY] Device info: %+v", info)

	switch info.DeviceType {
	case "entry":
		existingEntry, err := d.entrySvc.GetEntryDeviceByMac(ctx, info.MacAddress)
		if existingEntry != nil {
			log.Printf("[DISCOVERY] Updating entry device %s (%s) IP/Port -> %s:%d", existingEntry.Name, existingEntry.MacAddress, ip, entry.Port)

			updateReq := dto.UpdateEntryDeviceRequest{
				IPAddress:    &ip,
				Port:         &entry.Port,
				InstanceName: &info.InstanceName,
				InstanceType: &info.InstanceType,
				LastSeen:     func(t int64) *int64 { return &t }(time.Now().Unix()),
			}

			if _, err := d.entrySvc.UpdateEntryDevice(ctx, existingEntry.ID, updateReq); err != nil {
				log.Printf("failed to update entry device: %v", err)
			}
		} else if existingEntry == nil {
			// new device, create it
			log.Printf("[DISCOVERY] Creating new entry device at %s:%d", ip, entry.Port)

			createReq := dto.CreateEntryDeviceRequest{
				MacAddress:   info.MacAddress,
				IPAddress:    ip,
				Port:         entry.Port,
				Name:         info.InstanceName,
				Description:  nil,
				ServiceType:  &entry.ServiceRecord.Instance,
				DeviceType:   info.DeviceType,
				InstanceType: info.InstanceType,
				InstanceName: info.InstanceName,
			}

			if _, err := d.entrySvc.CreateEntryDevice(ctx, createReq); err != nil {
				log.Printf("failed to create entry device: %v", err)
			}
		} else if err != nil {
			log.Printf("failed to get entry device by mac: %v", err)
		}
	case "access":
		// TODO: implement updating access device IP/Port

	default:
		log.Printf("[DISCOVERY] Unknown device type: %s", info.DeviceType)
	}
}

func (d *DiscoveryService) Stop() {
	if d.cancel != nil {
		d.cancel()
	}
}
