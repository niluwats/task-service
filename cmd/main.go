package main

import (
	"log"

	"github.com/niluwats/task-service/internal/db"
)

func main() {
	log.Println("Starting....")
	dbClient := db.ConnectDB()
	log.Println("connected to DB ", dbClient)
}
