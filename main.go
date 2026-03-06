package main

import (
	"context"
	"fmt"
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

	currentIP, err := ipProvider.GetCurrentIP(context.Background())
	if err != nil {
		log.Fatalf("Failed to get current IP: %v", err)
	}

	fmt.Printf("Current IP: %s\n", currentIP)

	for _, record := range cfg.ARecords {
		recordIP, err := dnsProvider.GetRecordIP(context.Background(), record)
		if err != nil {
			log.Fatalf("Failed to get record IP for %s: %v", record, err)
		}

		if currentIP == recordIP {
			fmt.Printf("%s is up to date\n", record)
			continue
		}

		if err := dnsProvider.UpdateRecordIP(context.Background(), record, currentIP); err != nil {
			log.Fatalf("Failed to update record IP for %s: %v", record, err)
		}

		fmt.Printf("%s updated to %s\n", record, currentIP)
	}
}
