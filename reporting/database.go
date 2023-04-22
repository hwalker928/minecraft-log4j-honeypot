package reporting

import (
	"fmt"
	"log"
	"time"

	"github.com/hwalker928/minecraft-log4j-honeypot/database"
)

func UpdateIPValues(ip string, server database.ServerConfig) {
	dbConn := database.GetDB()

	query := fmt.Sprintf("SELECT COUNT(*) FROM attempts WHERE ip = '%s' AND reporting_server = '%s'", ip, server.Name)
	var count int
	err := dbConn.QueryRow(query).Scan(&count)
	if err != nil {
		log.Println(err)
	}

	if count == 0 {
		log.Println("New IP detected, inserting into database")
		_, err = dbConn.Exec("INSERT INTO attempts (ip, last_attempt, reporting_server) VALUES ($1, $2, $3)", ip, time.Now(), server.Name)
		if database.GetConfig().AbuseIPDB.Enabled {
			AbuseIPDBReport(ip, server)
		}
		return
	}

	_, err = dbConn.Exec("UPDATE attempts SET attempts = attempts + 1 WHERE ip = $1 AND reporting_server = $2", ip, server.Name)
}
