package discovery

import (
	"context"
	"fmt"
	"strings"
	"time"

	"bambu-farm/service"

	"github.com/grandcat/zeroconf"
	"go.uber.org/zap"
)

type DiscoveryEngine struct {
	logger         *zap.SugaredLogger
	printerService *service.PrinterService
}

func NewDiscoveryEngine(logger *zap.SugaredLogger, printerService *service.PrinterService) *DiscoveryEngine {
	return &DiscoveryEngine{
		logger:         logger,
		printerService: printerService,
	}
}

// Start begins background mDNS discovery
func (d *DiscoveryEngine) Start(ctx context.Context) {
	d.logger.Info("Starting BambuLab printer discovery engine...")

	go func() {
		// Run discovery periodically
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		d.scanNetwork()

		for {
			select {
			case <-ctx.Done():
				d.logger.Info("Stopping discovery engine")
				return
			case <-ticker.C:
				d.scanNetwork()
			}
		}
	}()
}

func (d *DiscoveryEngine) scanNetwork() {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		d.logger.Errorf("Failed to initialize zeroconf resolver: %v", err)
		return
	}

	entries := make(chan *zeroconf.ServiceEntry)
	// BambuLab printers typically broadcast on _bambu._tcp or similar mDNS services
	err = resolver.Browse(context.Background(), "_bambu._tcp", "local.", entries)
	if err != nil {
		d.logger.Errorf("mDNS browse failed: %v", err)
		return
	}

	for entry := range entries {
		d.logger.Infof("Discovered potential Bambu printer: %s (IP: %v)", entry.Instance, entry.AddrIPv4)
		d.detectBambuPrinter(entry)
	}
}

func (d *DiscoveryEngine) detectBambuPrinter(entry *zeroconf.ServiceEntry) {
	// Extract info from TXT records or instance name
	// Example entry.Text: ["model=X1C", "fw_ver=01.06.00.00"]
	
	ipAddress := ""
	if len(entry.AddrIPv4) > 0 {
		ipAddress = entry.AddrIPv4[0].String()
	}

	if ipAddress == "" {
		return
	}

	model := "Unknown"
	for _, txt := range entry.Text {
		if strings.HasPrefix(txt, "model=") {
			model = strings.TrimPrefix(txt, "model=")
		}
	}

	// Assuming instance name contains the serial or printer ID
	printerID := entry.Instance 

	d.logger.Infof("Detected Bambu Printer [%s] IP: %s, Model: %s", printerID, ipAddress, model)

	// Since we need an access code which cannot be known via mDNS, we cannot automatically 
	// fully register it into the database for full control yet. In a real system, we would:
	// 1. Add it to a "discovered_printers" table.
	// 2. Prompt the user on the UI to provide the Access Code.
	// For now, we log the discovery.
	
	fmt.Printf("TODO: Prompt user for access code for printer %s at %s\n", printerID, ipAddress)
}
