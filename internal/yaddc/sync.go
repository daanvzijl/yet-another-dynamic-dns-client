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

	for _, domain := range records {
		record, err := dnsProvider.GetRecord(ctx, domain)
		if err != nil {
			return fmt.Errorf("failed to get record for %s: %w", domain, err)
		}

		if currentIP == record.IP() {
			continue
		}

		if err := record.Update(ctx, currentIP); err != nil {
			return fmt.Errorf("failed to update record for %s: %w", domain, err)
		}
		log.Printf("%s updated to %s", domain, currentIP)
	}

	return nil
}
