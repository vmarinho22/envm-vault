package main

import (
	"log"

	"github.com/envm-vault/envm-vault/api/internal/database"
	"github.com/envm-vault/envm-vault/api/internal/env"
	"github.com/envm-vault/envm-vault/api/internal/secretcrypto"
)

func main() {
	if err := env.Load(); err != nil {
		log.Fatalf("load env: %v", err)
	}

	crypter, err := secretcrypto.NewAESGCM(env.Envs.Crypto.Key, env.Envs.Crypto.Version)
	if err != nil {
		log.Fatalf("init encryption: %v", err)
	}

	db, err := database.Open(env.Envs.Database.URL, crypter)
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
