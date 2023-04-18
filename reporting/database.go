package reporting

import (
	"fmt"
	"log"
	"time"

	"github.com/hwalker928/minecraft-log4j-honeypot/database"
)

func UpdateIPValues(ip string, server database.ServerConfig) {
	dbConn := database.GetDB()

	// update the last attempt time
	_, err := dbConn.Exec("UPDATE attempts SET last_attempt = $1 WHERE ip = $2 AND reporting_server = $3", time.Now(), ip, server.Name)

	query := fmt.Sprintf("SELECT COUNT(*) FROM attempts WHERE ip = '%s' AND reporting_server = '%s'", ip, server.Name)
	var count int
	err = dbConn.QueryRow(query).Scan(&count)
	if err != nil {
		log.Println(err)
	}

	if count == 0 {
		log.Println("New IP detected, inserting into database")
		_, err = dbConn.Exec("INSERT INTO attempts (ip, last_attempt) VALUES ($1, $2, $3)", ip, time.Now(), server.Name)
		return
	}

	_, err = dbConn.Exec("UPDATE attempts SET attempts = attempts + 1 WHERE ip = $1 AND reporting_server = $2", ip, server.Name)

	// if the ip has already been reported, skip it
	var reported bool
	err = dbConn.QueryRow("SELECT abuseipdb_reported FROM attempts WHERE ip = $1 AND reporting_server = $2", ip, server.Name).Scan(&reported)
	if err != nil {
		log.Println("Error getting reported status")
	} else {
		if reported {
			log.Println("IP is already reported")
			return
		}
	}

	// if the last attempt was less than 5 minutes ago, block the ip
	var lastAttempt time.Time
	err = dbConn.QueryRow("SELECT last_attempt FROM attempts WHERE ip = $1 AND reporting_server = $2", ip, server.Name).Scan(&lastAttempt)
	if err != nil {
		log.Println("Error getting last attempt time")
	} else {
		if time.Since(lastAttempt) < 5*time.Minute && database.GetConfig().AbuseIPDB.Enabled {
			log.Println("IP is now being reported")
			AbuseIPDBReport(ip, server)
			_, err = dbConn.Exec("UPDATE attempts SET abuseipdb_reported = 1 WHERE ip = $1 AND reporting_server = $2", ip, server.Name)
			return
		}
	}
}
