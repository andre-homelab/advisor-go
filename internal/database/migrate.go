package database

import (
	"github.com/andre-felipe-wonsik-alves/internal/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&models.Task{})
}
