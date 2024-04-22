package database

import (
	"fmt"

	"github.com/Chista-Framework/Chista/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(config *models.ConfigDB) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	dbConnection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return dbConnection, nil
}
