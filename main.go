package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {
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

	res := processIPListAsData(ips, network)
	fmt.Printf("You matched %d ip addresses as data", len(res))
	res = processIPListAsRegex(ips)
	fmt.Printf("You matched %d ip addresses as regex", len(res))
}

func processIPListAsData(ipList []net.IP, cidrRange *net.IPNet) []net.IP {
	var results []net.IP

	for _, ip := range ipList {
		if cidrRange.Contains(ip) {
			results = append(results, ip)
		}
	}

	return results
}

func processIPListAsRegex(ipList []net.IP) []net.IP {
	var results []net.IP

	re := regexp.MustCompile(`^10\.32\.16\.\d{1,3}$`)

	for _, ip := range ipList {
		if re.MatchString(ip.String()) {
			results = append(results, ip)
		}
	}

	return results
}
