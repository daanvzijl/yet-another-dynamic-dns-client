package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"yaddc/internal/yaddc"
)

// timeout bounds the whole run so a stalled connection can't hang the process.
const timeout = 2 * time.Minute

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := yaddc.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ipProvider, err := yaddc.NewIPProvider()
	if err != nil {
		return fmt.Errorf("failed to create IP provider: %w", err)
	}

	dnsProvider, err := yaddc.NewDNSProvider()
	if err != nil {
		return fmt.Errorf("failed to create DNS provider: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := yaddc.SyncRecords(ctx, ipProvider, dnsProvider, cfg.ARecords); err != nil {
		return fmt.Errorf("failed to sync records: %w", err)
	}
	return nil
}
