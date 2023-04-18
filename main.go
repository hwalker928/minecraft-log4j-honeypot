package main

import (
	"log"
	"regexp"
	"time"

	"github.com/hwalker928/minecraft-log4j-honeypot/database"
	"github.com/hwalker928/minecraft-log4j-honeypot/extractor"
	"github.com/hwalker928/minecraft-log4j-honeypot/minecraft"
	"github.com/hwalker928/minecraft-log4j-honeypot/reporting"
)

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
		abuseipdb_reported INTEGER NOT NULL DEFAULT 0,
		reporting_server TEXT NOT NULL
	  )`)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database table 'attempts' created or already exists.")
	}

	if config.AbuseIPDB.Enabled {
		log.Println("AbuseIPDB reporting is enabled")
	}

	// get number of rows in db
	var count int
	err = dbConn.QueryRow("SELECT COUNT(*) FROM attempts").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Loaded", count, "IPs from database")

	for _, server := range config.Servers {
		log.Println("Starting server", server.Name)

		mcServer := minecraft.NewServer(":" + server.Port)

		mcServer.ChatMessageCallback = Analyse
		mcServer.AcceptLoginCallback = Analyse
		mcServer.Server = server

		go mcServer.Run()
		reporting.SendWebhook("Minecraft Honeypot: Started", "Service has started on port "+server.Port, 0x98fb98, server)
	}

	for {
		time.Sleep(time.Duration(1<<63 - 1))
	}
}
