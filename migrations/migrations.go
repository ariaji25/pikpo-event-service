package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("mode") == "native" {
		godotenv.Load(".env")
	}
	// Read database config from env values
	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DB := os.Getenv("POSTGRES_DB")
	POSTGRES_HOST := os.Getenv("POSTGRES_HOST")

	// Pull it into dbUrl string
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_HOST, POSTGRES_DB)
	// Connect to database

	m, err := migrate.New(
		"file://migrations/migrations_files",
		dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
