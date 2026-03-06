package yaddc

import (
	"context"
	"fmt"
	"log"
)

func SyncRecords(ctx context.Context, ipProvider IPProvider, dnsProvider DNSProvider, records []string) error {
	currentIP, err := ipProvider.GetCurrentIP(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current IP: %w", err)
	}

	for _, record := range records {
		recordIP, err := dnsProvider.GetRecordIP(ctx, record)
		if err != nil {
			return fmt.Errorf("failed to get record IP for %s: %w", record, err)
		}

		if currentIP == recordIP {
			continue
		}

		if err := dnsProvider.UpdateRecordIP(ctx, record, currentIP); err != nil {
			return fmt.Errorf("failed to update record IP for %s: %w", record, err)
		}
		log.Printf("%s updated to %s", record, currentIP)
	}

	return nil
}
