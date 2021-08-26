package models

import "gorm.io/gorm"

// Person contains basic information of an individual.
type Person struct {
	gorm.Model
	Name     string `json:"name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	Age      uint   `json:"age" binding:"required"`
}
