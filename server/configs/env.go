package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type LocalEnv struct {
	DatabaseHost		string
	DatabasePort		string
	DatabaseName		string
	DatabaseUser		string
	DatabasePassword	string
}

func NewLocalEnv() *LocalEnv {
	err := godotenv.Load()

	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	databaseHost, ok := os.LookupEnv("DATABASE_HOST")

	if !ok {
		log.Fatal("Failed to load DATABASE_HOST in .env file")
	}

	databasePort, ok := os.LookupEnv("DATABASE_PORT")

	if !ok {
		log.Fatal("Failed to load DATABASE_PORT in .env file")
	}

	databaseName, ok := os.LookupEnv("DATABASE_NAME")

	if !ok {
		log.Fatal("Failed to load DATABASE_NAME in .env file")
	}

	databaseUser, ok := os.LookupEnv("DATABASE_USER")

	if !ok {
		log.Fatal("Failed to load DATABASE_USER in .env file")
	}

	databasePassword, ok := os.LookupEnv("DATABASE_PASSWORD")

	if !ok {
		log.Fatal("Failed to load DATABASE_PASSWORD in .env file")
	}

	return &LocalEnv{
		DatabaseHost: databaseHost,
		DatabasePort: databasePort,
		DatabaseName: databaseName,
		DatabaseUser: databaseUser,
		DatabasePassword: databasePassword,
	}
}
