package yaddc

import (
	"context"
	"errors"
	"fmt"
	"log"
)

// SyncRecords points every record at the current public IP. A failure on one
// record doesn't stop the others; all failures are returned joined together.
func SyncRecords(
	ctx context.Context,
	ipProvider IPProvider,
	dnsProvider DNSProvider,
	records []string,
) error {
	currentIP, err := ipProvider.GetCurrentIP(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current IP: %w", err)
	}

	var errs []error
	for _, domain := range records {
		if err := syncRecord(ctx, dnsProvider, domain, currentIP); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func syncRecord(ctx context.Context, dnsProvider DNSProvider, domain, currentIP string) error {
	record, err := dnsProvider.GetRecord(ctx, domain)
	if err != nil {
		return fmt.Errorf("failed to get record for %s: %w", domain, err)
	}

	if currentIP == record.IP() {
		return nil
	}

	if err := record.Update(ctx, currentIP); err != nil {
		return fmt.Errorf("failed to update record for %s: %w", domain, err)
	}
	log.Printf("%s updated to %s", domain, currentIP)
	return nil
}
