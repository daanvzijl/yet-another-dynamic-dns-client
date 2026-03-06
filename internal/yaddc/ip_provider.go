package yaddc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type IpifyProvider struct {
	client *http.Client
}

type IpifyResponse struct {
	IP string `json:"ip"`
}

func NewIpifyProvider() *IpifyProvider {
	return &IpifyProvider{
		client: &http.Client{},
	}
}

func (p *IpifyProvider) GetCurrentIP(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.ipify.org?format=json",
		nil,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "yaddc/0.1 (+https://github.com/daanvzijl/yet-another-dynamic-dns-client)")

	resp, err := p.client.Do(req)
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

	var data IpifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if data.IP == "" {
		return "", errors.New("empty ip response")
	}

	return data.IP, nil
}
