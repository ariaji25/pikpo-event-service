package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type DBPool struct {
	DB *sql.DB
}

var Database = DBPool{}

func (dbPool *DBPool) Open() {
	// Read database config from env values
	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_DB := os.Getenv("POSTGRES_DB")
	POSTGRES_HOST := os.Getenv("POSTGRES_HOST")

	// Pull it into dbUrl string
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s/%s", POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_HOST, POSTGRES_DB)
	// Connect to database
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		log.Println("Error : ", err.Error())
		return
	}
	if err := db.Ping(); err != nil {
		log.Println("unable to reach database: ", err)
		return
	}
	log.Println("DB Connected to :: ", dbUrl)
	dbPool.DB = db
}
