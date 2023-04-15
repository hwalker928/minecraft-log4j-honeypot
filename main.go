package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/Adikso/minecraft-log4j-honeypot/extractor"
	"github.com/Adikso/minecraft-log4j-honeypot/minecraft"
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
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var config Config

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal(err)
	}

	if config.AbuseIPDB.Enabled {
		log.Println("AbuseIPDB reporting is enabled")
	}

	server := minecraft.NewServer(":" + config.Server.Port)
	server.ChatMessageCallback = Analyse
	server.AcceptLoginCallback = Analyse

	if err := server.Run(); err != nil {
		log.Println(err)
	}
}
