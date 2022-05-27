package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

var migrationAction string

func init() {
	flag.StringVar(&migrationAction, "action", "up", "migration action (up/down)")
	flag.Parse()
}

func main() {
	godotenv.Load("../../.env")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	m, err := migrate.New(
		"file://..",
		dbURL,
	)

	if err != nil {
		log.Fatal(err)
	}

	if migrationAction == "up" {
		if err := m.Up(); err != nil {
			if err.Error() == "no change" {
				log.Println("No change")
			} else {
				log.Fatal(err)
			}
		}
	} else if migrationAction == "down" {
		if err := m.Down(); err != nil {
			if err.Error() == "no change" {
				log.Println("No change")
			} else {
				log.Fatal(err)
			}
		}
	} else {
		log.Fatal("invalid migration action")
	}
}
