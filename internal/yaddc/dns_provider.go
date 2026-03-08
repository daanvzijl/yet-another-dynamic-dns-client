package yaddc

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
)

type CloudflareDNSProvider struct {
	client *cloudflare.Client
	zoneID string
}

type cloudflareDNSRecord struct {
	client   *cloudflare.Client
	zoneID   string
	recordID string
	ip       string
	name     string
}

func NewCloudflareProvider(apiToken, zoneID string) *CloudflareDNSProvider {
	client := cloudflare.NewClient(
		option.WithAPIToken(apiToken),
	)

	return &CloudflareDNSProvider{client: client, zoneID: zoneID}
}

func (c *CloudflareDNSProvider) GetRecord(ctx context.Context, domain string) (DNSRecord, error) {
	page, err := c.client.DNS.Records.List(ctx, dns.RecordListParams{
		ZoneID: cloudflare.F(c.zoneID),
		Name:   cloudflare.F(dns.RecordListParamsName{Exact: cloudflare.F(domain)}),
		Type:   cloudflare.F(dns.RecordListParamsTypeA),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list dns records: %w", err)
	}

	if len(page.Result) == 0 {
		return nil, fmt.Errorf("no A record found for: %s", domain)
	}

	r := page.Result[0]
	return &cloudflareDNSRecord{
		client:   c.client,
		zoneID:   c.zoneID,
		recordID: r.ID,
		ip:       r.Content,
		name:     domain,
	}, nil
}

func (r *cloudflareDNSRecord) IP() string {
	return r.ip
}

func (r *cloudflareDNSRecord) Update(ctx context.Context, ip string) error {
	_, err := r.client.DNS.Records.Update(ctx, r.recordID, dns.RecordUpdateParams{
		ZoneID: cloudflare.F(r.zoneID),
		Body: dns.ARecordParam{
			Name:    cloudflare.F(r.name),
			Type:    cloudflare.F(dns.ARecordTypeA),
			TTL:     cloudflare.F(dns.TTL(1)),
			Content: cloudflare.F(ip),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update dns record for %s: %w", r.name, err)
	}
	return nil
}
