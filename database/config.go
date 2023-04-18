package database

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type ServerConfig struct {
	Port          string `json:"port"`
	Name          string `json:"name"`
	MinecraftOnly bool   `json:"minecraftOnly"`
	Comment       string `json:"comment"`
}

type Config struct {
	Servers   []ServerConfig `json:"servers"`
	AbuseIPDB struct {
		Enabled bool   `json:"enabled"`
		APIKey  string `json:"apiKey"`
	} `json:"AbuseIPDB"`
	Discord struct {
		Enabled    bool   `json:"enabled"`
		WebhookURL string `json:"webhookURL"`
	} `json:"discord"`
}

var ConfigFile *Config

func OpenConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bytes, &ConfigFile)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConfig() *Config {
	return ConfigFile
}
