package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func AbuseIPDBReport(ipAddress string) {
	// Validate the IP address
	if !isValidIP(ipAddress) {
		fmt.Println("Invalid IP address!")
		return
	}

	// Prepare the HTTP POST request to abuseipdb
	url := "https://api.abuseipdb.com/api/v2/report"
	data := strings.NewReader("ip=" + ipAddress + "&categories=14,15") // "Port scan" and "Hacking" categories

	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Key", "<your API key>") // Replace with your actual API key

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Printf("Reported IP address (%s) to AbuseIPDB!", ipAddress)
	} else {
		fmt.Printf("Failed to report IP address to AbuseIPDB! Error: %s", resp.Body)
	}
}

func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	return parsedIP.To4() != nil
}
