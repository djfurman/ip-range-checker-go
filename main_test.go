package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"testing"
)

func TestProcessIpListAsData(t *testing.T) {
	cidr := "10.0.0.0/8"
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		t.Errorf("Error parsing CIDR range '%s': %v", cidr, err)
	}

	ipList := []net.IP{net.ParseIP("10.32.16.5"), net.ParseIP("192.168.0.1")}

	res := processIPListAsData(ipList, network)

	if len(res) != 1 {
		t.Errorf("Expected 1 IP address to match, got %d", len(res))
	}
}

func BenchmarkTestProcessIPListAsData(b *testing.B) {
	cidr := "10.32.16.0/20"
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Printf("Error parsing CIDR range '%s': %v\n", cidr, err)
	}

	file, err := os.Open("data.csv")
	if err != nil {
		panic(fmt.Errorf("Failed to open CSV file: %w", err))
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var ips []net.IP

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Errorf("Failed to read CSV file: %w", err))
		}

		ipStr := strings.TrimSpace(record[0])
		ip := net.ParseIP(ipStr)
		if ip == nil {
			panic(fmt.Errorf("Failed to parse IP address '%s'", record[0]))
		}

		ips = append(ips, ip)
	}

	for b.Loop() {
		res := processIPListAsData(ips, network)

		if len(res) != 2 {
			b.Errorf("Expected 2 IP addresses to match, got %d", len(res))
		}
	}
}

func BenchmarkTestProcessIPListAsRegex(b *testing.B) {
	file, err := os.Open("data.csv")
	if err != nil {
		panic(fmt.Errorf("Failed to open CSV file: %w", err))
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var ips []net.IP

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Errorf("Failed to read CSV file: %w", err))
		}

		ipStr := strings.TrimSpace(record[0])
		ip := net.ParseIP(ipStr)
		if ip == nil {
			panic(fmt.Errorf("Failed to parse IP address '%s'", record[0]))
		}

		ips = append(ips, ip)
	}

	for b.Loop() {
		res := processIPListAsRegex(ips)

		if len(res) != 1 {
			b.Errorf("Expected 1 IP addresses to match, got %d", len(res))
		}
	}
}
