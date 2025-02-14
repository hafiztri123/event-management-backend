package postgres

import (
	"log"

	"github.com/hafiztri123/src/internal/model"
	"gorm.io/gorm"
)


func RunMigrations(db *gorm.DB)  {
	setExtension(db)
	migration(db)
}

func setExtension(db *gorm.DB) {
	result  := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if result.Error != nil {
		log.Fatal("[FAIL] fail to set extension: %v ", result.Error)
	}
}

func migration(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.Event{},
	)
	if err != nil {
		log.Fatal("[FAIL] fail to migrate")
	}
}
