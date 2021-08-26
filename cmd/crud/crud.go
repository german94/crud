package main

import (
	"crud/controllers"
	"crud/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
)

func main() {
	models.ConnectDataBase(sqlite.Open("crud-db"))

	r := gin.Default()

	r.GET("/persons", controllers.ListPersons)
	r.GET("/persons/:id", controllers.FindPerson)
	r.POST("/persons", controllers.CreatePerson)

	r.GET("/pets", controllers.ListPets)
	r.GET("/pets/:id", controllers.FindPet)
	r.POST("/pets", controllers.CreatePet)
	r.DELETE("/pets/:id", controllers.CreatePet)

	r.PUT("/adopt", controllers.Adopt)

	r.Run()
}
