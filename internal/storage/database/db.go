package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"sync"
	"time"
)

var (
	db   *gorm.DB
	once sync.Once
)

func DbConnect() *gorm.DB {
	once.Do(func() {
		dsn := os.Getenv("DATABASE_URL")
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect to the DB: ", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal("failed to get sqlDB from gormDB: ", err)

		}
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(50)
		sqlDB.SetConnMaxLifetime(time.Hour)

	})
	return db
}
