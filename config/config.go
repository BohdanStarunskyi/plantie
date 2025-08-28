package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		panic("db url not found")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		panic("Couldn't init database " + err.Error())
	}
	DB = db

	fmt.Println("Connected to Postgres and migrated models")
}
