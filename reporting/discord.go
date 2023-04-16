package reporting

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/hwalker928/minecraft-log4j-honeypot/database"
)

func SendWebhook(title string, description string, color int) {
	if !database.GetConfig().Discord.Enabled {
		return
	}

	webhookURL := database.GetConfig().Discord.WebhookURL

	// Define the embed payload as a struct
	type Embed struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Color       int    `json:"color"`
		Footer      struct {
			Text string `json:"text"`
		} `json:"footer"`
	}

	// Define the main payload as a struct
	type Payload struct {
		Content string  `json:"content"`
		Embeds  []Embed `json:"embeds"`
	}

	// Create a new embed
	embed := Embed{
		Title:       title,
		Description: description,
		Color:       color,
		Footer: struct {
			Text string `json:"text"`
		}{
			Text: "Sent from " + database.GetConfig().Server.Name,
		},
	}

	// Create a new payload with the embed
	payload := Payload{
		Content: "",
		Embeds:  []Embed{embed},
	}

	// Marshal the payload into a JSON string
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// Create a new HTTP POST request with the JSON payload
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}

	// Set the content type header to application/json
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
}
