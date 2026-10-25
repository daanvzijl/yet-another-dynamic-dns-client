package yaddc

import (
	"slices"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    []string
		wantErr bool
	}{
		{name: "single", raw: "a.example.com", want: []string{"a.example.com"}},
		{
			name: "trims whitespace",
			raw:  " a.example.com , b.example.com ",
			want: []string{"a.example.com", "b.example.com"},
		},
		{
			name: "skips empty entries",
			raw:  "a.example.com,, ,b.example.com,",
			want: []string{"a.example.com", "b.example.com"},
		},
		{name: "unset", raw: "", wantErr: true},
		{name: "only separators", raw: " , ,", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("YADDC_A_RECORDS", tt.raw)

			cfg, err := LoadConfig()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got config %v", cfg)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if !slices.Equal(cfg.ARecords, tt.want) {
				t.Errorf("expected %v, got %v", tt.want, cfg.ARecords)
			}
		})
	}
}

func TestNewDNSProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		token    string
		zone     string
		wantErr  bool
	}{
		{name: "cloudflare", provider: "cloudflare", token: "token", zone: "zone"},
		{name: "cloudflare missing token", provider: "cloudflare", zone: "zone", wantErr: true},
		{name: "cloudflare missing zone", provider: "cloudflare", token: "token", wantErr: true},
		{name: "unset", provider: "", wantErr: true},
		{name: "unknown", provider: "route53", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("YADDC_DNS_PROVIDER", tt.provider)
			t.Setenv("CF_API_TOKEN", tt.token)
			t.Setenv("CF_ZONE_ID", tt.zone)

			p, err := NewDNSProvider()
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if p == nil {
				t.Fatal("expected provider, got nil")
			}
		})
	}
}
