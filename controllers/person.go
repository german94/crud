package controllers

import (
	"crud/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	personNotFound      = "Person not found!"
	personIDNotProvided = "id of the person must be provided"
)

type Person struct {
	ID       uint   `json:"id"`
	Name     string `json:"name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	Age      uint   `json:"age" binding:"required"`
}

type ListPersonsResponse struct {
	Data  []Person
	Error string
}

type FindPersonResponse struct {
	Data  *Person
	Error string
}

type CreatePersonRequest struct {
	Name     string `json:"name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	Age      uint   `json:"age" binding:"required"`
}

type CreatePersonResponse struct {
	Data  *Person
	Error string
}

// ListPersons writes a response containing a list of all the stored individuals.
func ListPersons(c *gin.Context) {
	var persons []models.Person
	if err := models.DB.Find(&persons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ListPersonsResponse{Error: err.Error()})
	}
	outputPersons := make([]Person, 0, len(persons))
	for _, p := range persons {
		outputPersons = append(outputPersons, personFromModel(&p))
	}

	c.JSON(http.StatusOK, ListPersonsResponse{Data: outputPersons})
}

// FindPerson writes a response with the person that matches the requested ID.
func FindPerson(c *gin.Context) {
	var person models.Person

	if c.Param("id") == "" {
		c.JSON(http.StatusBadRequest, FindPersonResponse{Error: personIDNotProvided})
		return
	}
	txResult := models.DB.Where("id = ?", c.Param("id")).First(&person)
	if txResult.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, FindPersonResponse{Error: personNotFound})
		return
	}
	if err := txResult.Error; err != nil {
		c.JSON(http.StatusInternalServerError, FindPersonResponse{Error: err.Error()})
		return
	}

	outputPet := personFromModel(&person)
	c.JSON(http.StatusOK, FindPersonResponse{Data: &outputPet})
}

// CreatePerson stores the information of a new person.
func CreatePerson(c *gin.Context) {
	var input CreatePersonRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, CreatePersonResponse{Error: err.Error()})
		return
	}

	person := models.Person{
		Name:     input.Name,
		LastName: input.LastName,
		Age:      input.Age,
	}

	if err := models.DB.Create(&person).Error; err != nil {
		c.JSON(http.StatusInternalServerError, CreatePersonResponse{Error: err.Error()})
		return
	}

	outputPerson := personFromModel(&person)
	c.JSON(http.StatusCreated, CreatePersonResponse{Data: &outputPerson})
}

func personFromModel(person *models.Person) Person {
	return Person{
		ID:       person.ID,
		Name:     person.Name,
		LastName: person.LastName,
		Age:      person.Age,
	}
}
