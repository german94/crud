package controllers

import (
	"crud/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	petNotFound      = "Pet not found!"
	petIDNotProvided = "id of the pet must be provided"
)

type ListPetsResponse struct {
	Data  []Pet
	Error string
}

type FindPetResponse struct {
	Data  *Pet
	Error string
}

type CreatePetRequest struct {
	Kind   string `json:"kind" binding:"required"`
	Specie string `json:"specie" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Age    uint   `json:"age" binding:"required"`
}

type CreatePetResponse struct {
	Data  *Pet
	Error string
}

type Pet struct {
	ID        uint
	Kind      string
	Specie    string
	Name      string
	Age       uint
	Available bool
}

// ListPets writes a response containing a list of all the stored individuals.
func ListPets(c *gin.Context) {
	var pets []models.Pet
	if err := models.DB.Find(&pets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ListPetsResponse{Error: err.Error()})
	}

	outputPets := make([]Pet, 0, len(pets))
	for _, p := range pets {
		outputPets = append(outputPets, petFromModel(&p))
	}

	c.JSON(http.StatusOK, ListPetsResponse{Data: outputPets})
}

// FindPet writes a response with the pet that matches the requested ID.
func FindPet(c *gin.Context) {
	var pet models.Pet

	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, FindPetResponse{Error: petIDNotProvided})
		return
	}
	txResult := models.DB.Where("id = ?", c.Param("id")).First(&pet)
	if txResult.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, FindPetResponse{Error: petNotFound})
		return
	}
	if err := txResult.Error; err != nil {
		c.JSON(http.StatusInternalServerError, FindPetResponse{Error: err.Error()})
		return
	}

	outputPet := petFromModel(&pet)
	c.JSON(http.StatusOK, FindPetResponse{Data: &outputPet})
}

// CreatePet stores the information of a new pet.
// A new added pet will be automatically available for adoption.
func CreatePet(c *gin.Context) {
	var input CreatePetRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, CreatePetResponse{Error: err.Error()})
		return
	}

	pet := models.Pet{
		Name:   input.Name,
		Age:    input.Age,
		Specie: input.Specie,
		Kind:   input.Kind,
	}

	if err := models.DB.Create(&pet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, CreatePetResponse{Error: err.Error()})
		return
	}

	outputPet := petFromModel(&pet)
	c.JSON(http.StatusCreated, CreatePetResponse{Data: &outputPet})
}

// DeletePet deletes a pet.
func DeletePet(c *gin.Context) {
	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": petIDNotProvided})
		return
	}

	if err := models.DB.Delete(&models.Pet{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func petFromModel(pet *models.Pet) Pet {
	return Pet{
		ID:        pet.ID,
		Name:      pet.Name,
		Age:       pet.Age,
		Kind:      pet.Kind,
		Specie:    pet.Specie,
		Available: pet.OwnerID == 0,
	}
}
