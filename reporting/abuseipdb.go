package reporting

import (
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/hwalker928/minecraft-log4j-honeypot/database"
)

func AbuseIPDBReport(ipAddress string, server database.ServerConfig) {
	// Validate the IP address
	if !isValidIP(ipAddress) {
		log.Println("Invalid IP address!")
		return
	}

	// Prepare the HTTP POST request to abuseipdb
	url := "https://api.abuseipdb.com/api/v2/report"
	comment := server.Comment
	comment = strings.ReplaceAll(comment, "$DATETIME", time.Now().UTC().Format("2006-01-02 15:04:05"))
	comment = strings.ReplaceAll(comment, "$NAME", server.Name)
	comment = strings.ReplaceAll(comment, "$PORT", server.Port)

	data := strings.NewReader("ip=" + ipAddress + "&comment=" + comment + "&categories=14,15") // "Port scan" and "Hacking" categories

	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Key", database.GetConfig().AbuseIPDB.APIKey) // Replace with your actual API key

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		log.Printf("Reported IP address (%s) to AbuseIPDB!", ipAddress)
		SendWebhook("AbuseIPDB Report: Success", "Reported IP address ("+ipAddress+") to AbuseIPDB!", 0xe06666, server)
	} else {
		log.Printf("Failed to report IP address to AbuseIPDB! Error: %s", resp.Body)
	}
}

func isValidIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	return parsedIP.To4() != nil
}
