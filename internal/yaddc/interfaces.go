package yaddc

import "context"

type IPProvider interface {
	GetCurrentIP(ctx context.Context) (string, error)
}

type DNSProvider interface {
	GetRecordIP(ctx context.Context, domain string) (string, error)
	UpdateRecordIP(ctx context.Context, domain string, ip string) error
}
