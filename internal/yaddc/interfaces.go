package yaddc

import "context"

type IPProvider interface {
	GetCurrentIP(ctx context.Context) (string, error)
}

type DNSRecord interface {
	IP() string
	Update(ctx context.Context, ip string) error
}

type DNSProvider interface {
	GetRecord(ctx context.Context, domain string) (DNSRecord, error)
}
