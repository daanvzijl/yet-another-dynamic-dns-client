package yaddc

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	ARecords []string
}

func LoadConfig() (Config, error) {
	raw := os.Getenv("YADDC_A_RECORDS")
	if raw == "" {
		return Config{}, fmt.Errorf("YADDC_A_RECORDS is not set")
	}

	records := strings.Split(raw, ",")
	for i, r := range records {
		records[i] = strings.TrimSpace(r)
	}

	return Config{ARecords: records}, nil
}

func NewIPProvider() (IPProvider, error) {
	return NewIpifyProvider(), nil
}

func NewDNSProvider() (DNSProvider, error) {
	apiToken := os.Getenv("CF_API_TOKEN")
	zoneID := os.Getenv("CF_ZONE_ID")

	if apiToken == "" {
		return nil, fmt.Errorf("CF_API_TOKEN is not set")
	}
	if zoneID == "" {
		return nil, fmt.Errorf("CF_ZONE_ID is not set")
	}

	return NewCloudflareProvider(apiToken, zoneID), nil
}
