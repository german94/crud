package controllers

import (
	"crud/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	cannotAdoptDueToAge = "Only adults can adopt pets!"
	petNotAvailable     = "The pet already has an owner!"
)

// AdoptRequest contains the information of a new adoption to be done.
type AdoptRequest struct {
	PetID   uint `json:"pet_id" binding:"required"`
	OwnerID uint `json:"owner_id" binding:"required"`
}

type AdoptResponse struct {
	Data  *Pet
	Error string
}

// Adopt handles an adoption request. If the pet is available for adoption,
// it will be assigned to the corresponding owner.
func Adopt(c *gin.Context) {
	var adopt AdoptRequest
	if err := c.ShouldBindJSON(&adopt); err != nil {
		c.JSON(http.StatusBadRequest, AdoptResponse{Error: err.Error()})
		return
	}

	var person models.Person
	if !validatePerson(c, &person, &adopt) {
		return
	}

	var pet models.Pet
	if !validatePet(c, &pet, &adopt) {
		return
	}

	if !performAdoption(c, &pet, &person) {
		return
	}

	outputPet := petFromModel(&pet)
	c.JSON(http.StatusCreated, AdoptResponse{Data: &outputPet})
}

func validatePerson(c *gin.Context, person *models.Person, adopt *AdoptRequest) bool {
	personTxResult := models.DB.Where("id = ?", adopt.OwnerID).First(&person)
	if personTxResult.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, AdoptResponse{Error: personNotFound})
		return false
	}
	if err := personTxResult.Error; err != nil {
		c.JSON(http.StatusInternalServerError, AdoptResponse{Error: err.Error()})
		return false
	}
	if person.Age < 18 {
		c.JSON(http.StatusConflict, AdoptResponse{Error: cannotAdoptDueToAge})
		return false
	}
	return true
}

func validatePet(c *gin.Context, pet *models.Pet, adopt *AdoptRequest) bool {
	petTxResult := models.DB.Where("id = ?", adopt.PetID).First(&pet)
	if petTxResult.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, AdoptResponse{Error: petNotFound})
		return false
	}
	if err := petTxResult.Error; err != nil {
		c.JSON(http.StatusInternalServerError, AdoptResponse{Error: err.Error()})
		return false
	}
	if pet.OwnerID != 0 {
		c.JSON(http.StatusConflict, AdoptResponse{Error: petNotAvailable})
		return false
	}
	return true
}

func performAdoption(c *gin.Context, pet *models.Pet, person *models.Person) bool {
	pet.OwnerID = person.ID

	if err := models.DB.Save(&pet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, AdoptResponse{Error: err.Error()})
		return false
	}

	return true
}
