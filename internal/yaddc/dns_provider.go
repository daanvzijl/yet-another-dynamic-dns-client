package yaddc

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
)

type CloudflareDNSProvider struct {
	client   *cloudflare.Client
	zoneID   string
	recordID string
}

func NewCloudflareProvider(apiToken, zoneID string) *CloudflareDNSProvider {
	client := cloudflare.NewClient(
		option.WithAPIToken(apiToken),
	)

	return &CloudflareDNSProvider{client: client, zoneID: zoneID}
}

func (c *CloudflareDNSProvider) GetRecordIP(ctx context.Context, domain string) (string, error) {
	page, err := c.client.DNS.Records.List(ctx, dns.RecordListParams{
		ZoneID: cloudflare.F(c.zoneID),
		Name:   cloudflare.F(dns.RecordListParamsName{Exact: cloudflare.F(domain)}),
		Type:   cloudflare.F(dns.RecordListParamsTypeA),
	})
	if err != nil {
		return "", fmt.Errorf("failed to list dns records: %w", err)
	}

	if len(page.Result) == 0 {
		return "", fmt.Errorf("no A record found for: %s", domain)
	}

	c.recordID = page.Result[0].ID

	return page.Result[0].Content, nil
}

func (c *CloudflareDNSProvider) UpdateRecordIP(ctx context.Context, name, ip string) error {
	_, err := c.client.DNS.Records.Update(ctx, c.recordID, dns.RecordUpdateParams{
		ZoneID: cloudflare.F(c.zoneID),
		Body: dns.ARecordParam{
			Name:    cloudflare.F(name),
			Type:    cloudflare.F(dns.ARecordTypeA),
			TTL:     cloudflare.F(dns.TTL(1)),
			Content: cloudflare.F(ip),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update dns record for %s: %w", name, err)
	}

	return nil
}
