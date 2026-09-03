package main

import (
	"log"

	"github.com/envm-vault/envm-vault/api/internal/config"
	"github.com/envm-vault/envm-vault/api/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.Open(cfg.DatabaseURL)
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
