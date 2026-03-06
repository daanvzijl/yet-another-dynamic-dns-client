package yaddc

import (
	"context"
	"fmt"
	"testing"
)

type mockIPProvider struct {
	ip  string
	err error
}

func (m *mockIPProvider) GetCurrentIP(_ context.Context) (string, error) {
	return m.ip, m.err
}

type mockDNSProvider struct {
	recordIP  string
	updatedTo string
	err       error
}

func (m *mockDNSProvider) GetRecordIP(_ context.Context, _ string) (string, error) {
	return m.recordIP, m.err
}

func (m *mockDNSProvider) UpdateRecordIP(_ context.Context, _ string, ip string) error {
	m.updatedTo = ip
	return m.err
}

// Tests

func TestSyncRecords_UpToDate(t *testing.T) {
	ip := &mockIPProvider{ip: "1.2.3.4"}
	dns := &mockDNSProvider{recordIP: "1.2.3.4"}

	err := SyncRecords(context.Background(), ip, dns, []string{"test.example.com"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if dns.updatedTo != "" {
		t.Fatalf("expected no update, got %s", dns.updatedTo)
	}
}

func TestSyncRecords_UpdatesRecord(t *testing.T) {
	ip := &mockIPProvider{ip: "1.2.3.4"}
	dns := &mockDNSProvider{recordIP: "4.3.2.1"}

	err := SyncRecords(context.Background(), ip, dns, []string{"test.example.com"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if dns.updatedTo != "1.2.3.4" {
		t.Fatalf("expected update to 1.2.3.4, got %s", dns.updatedTo)
	}
}

func TestSyncRecords_IPProviderError(t *testing.T) {
	ip := &mockIPProvider{err: fmt.Errorf("network error")}
	dns := &mockDNSProvider{}

	err := SyncRecords(context.Background(), ip, dns, []string{"test.example.com"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
