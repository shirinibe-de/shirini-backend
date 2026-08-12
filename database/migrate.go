package database

import (
	"log"
	"os"
)

func RunMigrations() {

	sqlBytes, err := os.ReadFile("database/migrations/001_create_users.sql")

	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(string(sqlBytes))

	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migration completed successfully")
}