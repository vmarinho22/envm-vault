package main

import (
	"log"

	"github.com/envm-vault/envm-vault/api/internal/database"
	"github.com/envm-vault/envm-vault/api/internal/env"
)

func main() {
	if err := env.Load(); err != nil {
		log.Fatalf("load env: %v", err)
	}

	db, err := database.Open(env.Envs.Database.URL)
	if err != nil {
		log.Fatalf("migrate: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("get sql db: %v", err)
	}
	defer sqlDB.Close()

	log.Println("database migrated successfully")
}
