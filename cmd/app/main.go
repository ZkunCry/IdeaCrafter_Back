package main

import (
	"log"
	"startup_back/internal/config"
	"startup_back/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("error loading config: %v", err) 
	}
	log.("Loaded config: %v",cfg)
	
	db, err := gorm.Open(postgres.Open(cfg.DBConnectionString()),&gorm.Config{})
	if err !=nil{
		log.Fatalf("Failed to connect to database: &v",err)

	}
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Startup{},
		&domain.Category{},
		&domain.Membership{},
		&domain.Favorite{},
		&domain.StartupFile{},
	)
	if err != nil{
		log.Fatalf("Failed to migrate database: %v",err)
	}
}