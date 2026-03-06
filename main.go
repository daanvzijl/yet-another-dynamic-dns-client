package main

import (
	"context"
	"log"

	"yaddc/internal/yaddc"
)

func main() {
	cfg, err := yaddc.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ipProvider, err := yaddc.NewIPProvider()
	if err != nil {
		log.Fatalf("Failed to create IP provider: %v", err)
	}

	dnsProvider, err := yaddc.NewDNSProvider()
	if err != nil {
		log.Fatalf("Failed to create DNS provider: %v", err)
	}

	if err := yaddc.SyncRecords(context.Background(), ipProvider, dnsProvider, cfg.ARecords); err != nil {
		log.Fatalf("Failed to sync records: %v", err)
	}
}
