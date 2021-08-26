package models

import "gorm.io/gorm"

// Pet contains basic information of an animal that may or may not be available for adoption.
type Pet struct {
	gorm.Model
	Kind    string  `json:"kind" binding:"required"`
	Specie  string  `json:"specie" binding:"required"`
	Name    string  `json:"name" binding:"required"`
	Age     uint    `json:"age" binding:"required"`
	OwnerID uint    `json:"owner_id"`
	Owner   *Person `gorm:"foreignKey:OwnerID" json:"-"`
}
