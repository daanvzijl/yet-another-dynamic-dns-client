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

type mockDNSRecord struct {
	ip        string
	updatedTo string
	err       error
}

func (m *mockDNSRecord) IP() string {
	return m.ip
}

func (m *mockDNSRecord) Update(_ context.Context, ip string) error {
	m.updatedTo = ip
	return m.err
}

type mockDNSProvider struct {
	record *mockDNSRecord
	err    error
}

func (m *mockDNSProvider) GetRecord(_ context.Context, _ string) (DNSRecord, error) {
	return m.record, m.err
}

func TestSyncRecords_UpToDate(t *testing.T) {
	ip := &mockIPProvider{ip: "1.2.3.4"}
	dns := &mockDNSProvider{record: &mockDNSRecord{ip: "1.2.3.4"}}

	err := SyncRecords(context.Background(), ip, dns, []string{"test.example.com"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if dns.record.updatedTo != "" {
		t.Fatalf("expected no update, got %s", dns.record.updatedTo)
	}
}

func TestSyncRecords_UpdatesRecord(t *testing.T) {
	ip := &mockIPProvider{ip: "1.2.3.4"}
	dns := &mockDNSProvider{record: &mockDNSRecord{ip: "4.3.2.1"}}

	err := SyncRecords(context.Background(), ip, dns, []string{"test.example.com"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if dns.record.updatedTo != "1.2.3.4" {
		t.Fatalf("expected update to 1.2.3.4, got %s", dns.record.updatedTo)
	}
}

func TestSyncRecords_IPProviderError(t *testing.T) {
	ip := &mockIPProvider{err: fmt.Errorf("network error")}
	dns := &mockDNSProvider{record: &mockDNSRecord{}}

	err := SyncRecords(context.Background(), ip, dns, []string{"test.example.com"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
