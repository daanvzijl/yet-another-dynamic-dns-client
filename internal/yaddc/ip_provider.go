package yaddc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

const ipifyURL = "https://api.ipify.org?format=json"

type ipifyProvider struct {
	client *http.Client
	url    string
}

type ipifyResponse struct {
	IP string `json:"ip"`
}

func newIpifyProvider() *ipifyProvider {
	return &ipifyProvider{
		client: &http.Client{},
		url:    ipifyURL,
	}
}

func (p *ipifyProvider) GetCurrentIP(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		p.url,
		nil,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set(
		"User-Agent",
		"yaddc/"+Version+" (+https://github.com/daanvzijl/yet-another-dynamic-dns-client)",
	)

	resp, err := p.client.Do(req) // #nosec G704
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var data ipifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data.IP == "" {
		return "", errors.New("empty ip response")
	}

	if ip := net.ParseIP(data.IP); ip == nil || ip.To4() == nil {
		return "", fmt.Errorf("invalid IPv4 address in response: %q", data.IP)
	}

	return data.IP, nil
}
