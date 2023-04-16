package main

import (
	"log"
	"regexp"

	"github.com/hwalker928/minecraft-log4j-honeypot/database"
	"github.com/hwalker928/minecraft-log4j-honeypot/extractor"
	"github.com/hwalker928/minecraft-log4j-honeypot/minecraft"
	"github.com/hwalker928/minecraft-log4j-honeypot/reporting"
)

type Config struct {
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
	AbuseIPDB struct {
		Enabled bool   `json:"enabled"`
		APIKey  string `json:"apiKey"`
	} `json:"AbuseIPDB"`
}

func Analyse(text string) {
	log.Printf("Testing text: %s\n", text)

	pattern := regexp.MustCompile(`\${jndi:(.*)}`)
	finder := extractor.NewFinder(pattern)

	injections := finder.FindInjections(text)
	for _, url := range injections {
		log.Printf("Fetching payload for: jndi:%s", url.String())

		files, err := extractor.FetchFromLdap(url)
		if err != nil {
			log.Printf("Failed to fetch class from %s", url)
			continue
		}

		for _, filename := range files {
			log.Printf("Saved payload to file %s\n", filename)
		}
	}
}

func main() {
	database.OpenConfig()
	database.OpenDatabase("./database.db")
	defer database.CloseDatabase()

	dbConn := database.GetDB()
	config := database.GetConfig()

	_, err := dbConn.Exec(`CREATE TABLE IF NOT EXISTS attempts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT NOT NULL,
		attempts INTEGER NOT NULL DEFAULT 1,
		last_attempt DATETIME NOT NULL,
		abuseipdb_reported INTEGER NOT NULL DEFAULT 0
	  )`)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database table 'attempts' created or already exists.")
	}

	if config.AbuseIPDB.Enabled {
		log.Println("AbuseIPDB reporting is enabled")
	}

	reporting.SendWebhook("Minecraft Honeypot: Started", "Service has started on port "+config.Server.Port, 0x98fb98)

	server := minecraft.NewServer(":" + config.Server.Port)
	server.ChatMessageCallback = Analyse
	server.AcceptLoginCallback = Analyse

	if err := server.Run(); err != nil {
		log.Println(err)
	}
}
