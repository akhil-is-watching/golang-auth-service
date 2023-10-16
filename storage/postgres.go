package storage

import (
	"github.com/akhil-is-watching/authservice/config"
	"github.com/akhil-is-watching/authservice/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var defaultDB *gorm.DB

func ConnectDB(config *config.Config) {
	db, err := gorm.Open(postgres.Open(config.DBUrl), &gorm.Config{})
	if err != nil {
		panic("DB Connection failed")
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	db.Logger = logger.Default.LogMode(logger.Info)

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic("DB Migrations Failed")
	}

	defaultDB = db
}

func GetDB() *gorm.DB {
	return defaultDB
}
