package models

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDataBase(dialector gorm.Dialector) {
	database, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&Person{})
	database.AutoMigrate(&Pet{})

	DB = database
}
