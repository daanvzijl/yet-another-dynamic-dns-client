package yaddc

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIpifyProvider_GetCurrentIP(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		want    string
		wantErr bool
	}{
		{name: "ok", status: http.StatusOK, body: `{"ip":"1.2.3.4"}`, want: "1.2.3.4"},
		{
			name:    "non-200",
			status:  http.StatusServiceUnavailable,
			body:    `{"ip":"1.2.3.4"}`,
			wantErr: true,
		},
		{name: "bad json", status: http.StatusOK, body: `not json`, wantErr: true},
		{name: "empty ip", status: http.StatusOK, body: `{"ip":""}`, wantErr: true},
		{name: "not an ip", status: http.StatusOK, body: `{"ip":"<html>"}`, wantErr: true},
		{name: "ipv6", status: http.StatusOK, body: `{"ip":"2001:db8::1"}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if ua := r.Header.Get("User-Agent"); !strings.HasPrefix(ua, "yaddc/") {
						t.Errorf("unexpected User-Agent %q", ua)
					}
					w.WriteHeader(tt.status)
					_, _ = w.Write([]byte(tt.body))
				}),
			)
			defer srv.Close()

			p := newIpifyProvider()
			p.url = srv.URL

			got, err := p.GetCurrentIP(t.Context())
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
