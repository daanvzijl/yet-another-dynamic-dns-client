package yaddc

import (
	"context"
	"errors"
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
	records map[string]*mockDNSRecord
	errs    map[string]error
	calls   []string
}

func (m *mockDNSProvider) GetRecord(_ context.Context, domain string) (DNSRecord, error) {
	m.calls = append(m.calls, domain)
	if err := m.errs[domain]; err != nil {
		return nil, err
	}
	r, ok := m.records[domain]
	if !ok {
		return nil, errors.New("no such record")
	}
	return r, nil
}

func TestSyncRecords_UpToDate(t *testing.T) {
	ip := &mockIPProvider{ip: "1.2.3.4"}
	record := &mockDNSRecord{ip: "1.2.3.4"}
	dns := &mockDNSProvider{records: map[string]*mockDNSRecord{"a.example.com": record}}

	err := SyncRecords(t.Context(), ip, dns, []string{"a.example.com"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if record.updatedTo != "" {
		t.Errorf("expected no update, got %s", record.updatedTo)
	}
}

func TestSyncRecords_UpdatesRecord(t *testing.T) {
	ip := &mockIPProvider{ip: "1.2.3.4"}
	record := &mockDNSRecord{ip: "4.3.2.1"}
	dns := &mockDNSProvider{records: map[string]*mockDNSRecord{"a.example.com": record}}

	err := SyncRecords(t.Context(), ip, dns, []string{"a.example.com"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if record.updatedTo != "1.2.3.4" {
		t.Errorf("expected update to 1.2.3.4, got %s", record.updatedTo)
	}
}

func TestSyncRecords_MultipleRecords(t *testing.T) {
	ip := &mockIPProvider{ip: "1.2.3.4"}
	current := &mockDNSRecord{ip: "1.2.3.4"}
	stale := &mockDNSRecord{ip: "4.3.2.1"}
	dns := &mockDNSProvider{records: map[string]*mockDNSRecord{
		"current.example.com": current,
		"stale.example.com":   stale,
	}}

	err := SyncRecords(
		t.Context(),
		ip,
		dns,
		[]string{"current.example.com", "stale.example.com"},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if current.updatedTo != "" {
		t.Errorf("expected current record untouched, got update to %s", current.updatedTo)
	}
	if stale.updatedTo != "1.2.3.4" {
		t.Errorf("expected stale record updated to 1.2.3.4, got %s", stale.updatedTo)
	}
}

func TestSyncRecords_NoRecords(t *testing.T) {
	ip := &mockIPProvider{ip: "1.2.3.4"}
	dns := &mockDNSProvider{}

	err := SyncRecords(t.Context(), ip, dns, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(dns.calls) != 0 {
		t.Errorf("expected no DNS lookups, got %v", dns.calls)
	}
}

func TestSyncRecords_IPProviderError(t *testing.T) {
	netErr := errors.New("network error")
	ip := &mockIPProvider{err: netErr}
	dns := &mockDNSProvider{}

	err := SyncRecords(t.Context(), ip, dns, []string{"a.example.com"})
	if !errors.Is(err, netErr) {
		t.Fatalf("expected wrapped network error, got %v", err)
	}
	if len(dns.calls) != 0 {
		t.Errorf("expected no DNS lookups, got %v", dns.calls)
	}
}

func TestSyncRecords_GetRecordErrorContinues(t *testing.T) {
	lookupErr := errors.New("lookup failed")
	ip := &mockIPProvider{ip: "1.2.3.4"}
	stale := &mockDNSRecord{ip: "4.3.2.1"}
	dns := &mockDNSProvider{
		records: map[string]*mockDNSRecord{"b.example.com": stale},
		errs:    map[string]error{"a.example.com": lookupErr},
	}

	err := SyncRecords(t.Context(), ip, dns, []string{"a.example.com", "b.example.com"})
	if !errors.Is(err, lookupErr) {
		t.Fatalf("expected wrapped lookup error, got %v", err)
	}
	if stale.updatedTo != "1.2.3.4" {
		t.Errorf("expected b.example.com updated despite earlier failure, got %q", stale.updatedTo)
	}
}

func TestSyncRecords_UpdateErrorContinues(t *testing.T) {
	updateErr := errors.New("update failed")
	ip := &mockIPProvider{ip: "1.2.3.4"}
	failing := &mockDNSRecord{ip: "4.3.2.1", err: updateErr}
	stale := &mockDNSRecord{ip: "4.3.2.1"}
	dns := &mockDNSProvider{records: map[string]*mockDNSRecord{
		"a.example.com": failing,
		"b.example.com": stale,
	}}

	err := SyncRecords(t.Context(), ip, dns, []string{"a.example.com", "b.example.com"})
	if !errors.Is(err, updateErr) {
		t.Fatalf("expected wrapped update error, got %v", err)
	}
	if stale.updatedTo != "1.2.3.4" {
		t.Errorf("expected b.example.com updated despite earlier failure, got %q", stale.updatedTo)
	}
}
