package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
)

type IPAddr struct {
	Number string `json:"ip"`
}

var dnsProvider = struct {
	Token     string
	ZoneID    string
	RecordIds string
}{
	Token:     os.Getenv("CLOUDFLARE_API_TOKEN"),
	ZoneID:    os.Getenv("CLOUDFLARE_ZONE_ID"),
	RecordIds: os.Getenv("CLOUDFLARE_DNS_RECORD_ID"),
}

func main() {
	if err := checkEnv(); err != nil {
		panic(err)
	}

	currentIP := getCurrentIP()
	recordData := getRecordData()

	if compareIPs(currentIP, recordData) {
		fmt.Println("No update required")
	} else {
		fmt.Println("Update required...")
		updateRecordData(currentIP)
	}
}

func compareIPs(ip1, ip2 string) bool {
	return ip1 == ip2
}

func isValidIPv4(ip string) bool {
	return net.ParseIP(ip) != nil && net.ParseIP(ip).To4() != nil
}

func checkEnv() error {
	missing := []string{}
	if dnsProvider.Token == "" {
		missing = append(missing, "CLOUDFLARE_API_TOKEN")
	}
	if dnsProvider.ZoneID == "" {
		missing = append(missing, "CLOUDFLARE_ZONE_ID")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}
	return nil
}

func updateRecordData(ip string) {
	client := cloudflare.NewClient(
		option.WithAPIToken(dnsProvider.Token),
	)
	updateParams := dns.RecordUpdateParams{
		ZoneID: cloudflare.F(dnsProvider.ZoneID),
		Body: dns.ARecordParam{
			Name:    cloudflare.F("vpn.vanzijl.io"),
			Type:    cloudflare.F(dns.ARecordTypeA),
			Content: cloudflare.F(ip),
			TTL:     cloudflare.F(dns.TTL1),
		},
	}
	recordResponse, err := client.DNS.Records.Update(
		context.Background(),
		dnsProvider.RecordIds,
		updateParams,
	)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Updated content to: %+v\n", recordResponse.Content)
}

func getRecordData() string {
	client := cloudflare.NewClient(
		option.WithAPIToken(dnsProvider.Token),
	)
	recordResponse, err := client.DNS.Records.Get(
		context.Background(),
		dnsProvider.RecordIds,
		dns.RecordGetParams{
			ZoneID: cloudflare.F(dnsProvider.ZoneID),
		},
	)
	if err != nil {
		panic(err.Error())
	}
	if recordResponse.Type != "A" {
		panic(err.Error())
	}
	return recordResponse.Content
}

func getCurrentIP() string {
	url := "https://api.ipify.org?format=json"

	yaddcClient := http.Client{
		Timeout: time.Second * 5, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "yaddc/0.1 (+https://github.com/daanvzijl/yet-another-dynamic-dns-client)")

	res, getErr := yaddcClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var ip IPAddr
	jsonErr := json.Unmarshal(body, &ip)
	if jsonErr != nil {
		log.Fatalf("unable to parse value: %q, error: %s", string(body), jsonErr.Error())
	}

	if !isValidIPv4(ip.Number) {
		log.Fatalf("invalid IPv4 address: %s", ip.Number)
	}

	return ip.Number
}
