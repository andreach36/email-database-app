package main

import (
	datasetindex "email-database-api/dataset-index"
	"log"

	"github.com/joho/godotenv"
)



func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnvVars()
	datasetindex.IndexAndCreateJson()
}

