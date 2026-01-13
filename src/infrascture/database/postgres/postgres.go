package database_postgres

import (
	"app/conf"
	entity_accounts "app/entity/accounts"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var err error
	config := conf.LoadConfig()
	host := config.DBHost
	user := config.DBUser
	password := config.DBPassword
	dbname := config.DBName
	port := config.DBPort
	DB, err = gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	return DB
}

func RunMigrations(DB *gorm.DB) {
	if DB == nil {
		log.Fatal("Database connection not initialized")
	}

	DB.AutoMigrate(&entity_accounts.User{})

}
