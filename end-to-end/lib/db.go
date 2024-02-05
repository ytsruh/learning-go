package lib

import (
	"errors"
	"fmt"
	"os"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"homethings.ytsruh.com/models"
)

var DBConn *gorm.DB

func GetDB() *gorm.DB {
	return DBConn
}

func InitDB() *gorm.DB {
	dburl := os.Getenv("DATABASE_URL")
	var err error
	DBConn, err = gorm.Open(postgres.Open(dburl))
	if err != nil {
		fmt.Println("Failed to connect to database")
		panic("Failed to connect to database")
	}

	// Enable uuid-ossp extension
	err = DBConn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		fmt.Println("Failed to enable uuid-ossp extension")
		panic(err)
	}

	DBConn.AutoMigrate(&models.Account{}, &models.User{}, &models.Feedback{}, &models.Book{}, &models.Document{})

	return DBConn
}

func SetupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, errors.New("an error occured when opening a stub database connection")
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, nil, errors.New("an error occured when opening gorm connection")
	}
	return db, mock, err
}
