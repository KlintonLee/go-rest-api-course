package database

import (
	"github.com/KlintonLee/go-rest-api-course/internal/database/models"
	"github.com/jinzhu/gorm"
)

//MigrateDB - migrates our database and creates our comment table
func MigrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Comment{}).Error
	if err != nil {
		return (err)
	}

	return nil
}
