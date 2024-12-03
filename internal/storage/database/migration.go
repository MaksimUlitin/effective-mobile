package database

import (
	"effectiveMobileTask/internal/models"
	"gorm.io/gorm"
	"log"
)

func Migrate(db *gorm.DB) {

	if err := db.AutoMigrate(&models.Song{}); err != nil {
		log.Fatal("migrate failed:  ", err)
	}
}
