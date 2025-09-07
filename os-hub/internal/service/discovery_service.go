package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"opensesame/internal/etag"
	"opensesame/internal/models/db"
	"opensesame/internal/models/dto"
	"opensesame/internal/repository"

	"github.com/grandcat/zeroconf"
)

type DiscoveryService struct {
	resolver *zeroconf.Resolver
	entries  chan *zeroconf.ServiceEntry
	cancel   context.CancelFunc
	repo     repository.DiscoveredDeviceRepository
	entrySvc *EntryService
}

type deviceInfo struct {
	MacAddress   string `json:"mac_address"`
	DeviceType   string `json:"type"`          // e.g. "entry", "access"
	InstanceName string `json:"instance_name"` // e.g. "OpenSesame Relay Lock"
	InstanceType string `json:"instance_type"` // e.g. "relay_lock"
}

func NewDiscoveryService(
	repo repository.DiscoveredDeviceRepository,
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

	ip := entry.AddrIPv4[0].String()
	url := fmt.Sprintf("http://%s:%d/info", ip, entry.Port)

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("failed to fetch /info from %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	var info deviceInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Printf("failed to decode /info from %s: %v", url, err)
		return
	}

	device := &db.DiscoveredDevice{
		MacAddress:   info.MacAddress,
		DeviceType:   info.DeviceType,
		InstanceType: info.InstanceType,
		InstanceName: info.InstanceName,
		IPAddress:    ip,
		Port:         entry.Port,
		ServiceType:  entry.Service,
		LastSeen:     time.Now(),
	}

	if err := d.repo.Upsert(ctx, device); err != nil {
		log.Printf("failed to save discovered device: %v", err)
		return
	}

	log.Printf("[DISCOVERY] Saved device %s (%s) at %s:%d", info.InstanceName, info.MacAddress, ip, entry.Port)

	switch info.DeviceType {
	case "entry":
		// update entry device if we discover a MAC address that's changed IPs
		entryDevice, err := d.entrySvc.GetEntryDeviceByMac(ctx, info.MacAddress)
		if err == nil && (entryDevice.IPAddress != ip || entryDevice.Port != entry.Port) {
			log.Printf("[DISCOVERY] Updating entry device %s (%s) IP/Port -> %s:%d", entryDevice.Name, entryDevice.MacAddress, ip, entry.Port)

			updateReq := dto.UpdateEntryDeviceRequest{
				IPAddress: &ip,
				Port:      &entry.Port,
				LastSeen:  func(t int64) *int64 { return &t }(time.Now().Unix()),
			}

			if _, err := d.entrySvc.UpdateEntryDeviceInfo(ctx, entryDevice.ID, updateReq); err != nil {
				log.Printf("failed to update entry device: %v", err)
			}
		}
	case "access":
		// TODO: implement updating access device IP/Port if needed
	default:
		log.Printf("[DISCOVERY] Unknown device type: %s", info.DeviceType)
	}

	etag.Bump()
}

func (d *DiscoveryService) Stop() {
	if d.cancel != nil {
		d.cancel()
	}
}
