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

	var records []string
	for r := range strings.SplitSeq(raw, ",") {
		if r = strings.TrimSpace(r); r != "" {
			records = append(records, r)
		}
	}
	if len(records) == 0 {
		return Config{}, fmt.Errorf("YADDC_A_RECORDS contains no records")
	}

	return Config{ARecords: records}, nil
}

func NewIPProvider() (IPProvider, error) {
	return newIpifyProvider(), nil
}

func NewDNSProvider() (DNSProvider, error) {
	provider := os.Getenv("YADDC_DNS_PROVIDER")

	switch provider {
	case "cloudflare":
		apiToken := os.Getenv("CF_API_TOKEN")
		zoneID := os.Getenv("CF_ZONE_ID")
		if apiToken == "" {
			return nil, fmt.Errorf("CF_API_TOKEN is not set")
		}
		if zoneID == "" {
			return nil, fmt.Errorf("CF_ZONE_ID is not set")
		}
		return NewCloudflareProvider(apiToken, zoneID), nil
	case "":
		return nil, fmt.Errorf("YADDC_DNS_PROVIDER is not set")
	default:
		return nil, fmt.Errorf("unknown DNS provider: %s", provider)
	}
}
