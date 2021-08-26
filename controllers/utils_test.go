package controllers

import (
	"bytes"
	"crud/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
)

func setUpDB() {
	testDBName := "test-crud-db"
	os.Remove(testDBName)
	models.ConnectDataBase(sqlite.Open(testDBName))
	addTestPersons()
	addTestPets()
}

func addTestPersons() {
	models.DB.Save(&models.Person{
		Name:     "German",
		LastName: "Pinzon",
		Age:      27,
	})
	models.DB.Save(&models.Person{
		Name:     "Camila",
		LastName: "Nani",
		Age:      26,
	})
	models.DB.Save(&models.Person{
		Name:     "Oscar",
		LastName: "Pina",
		Age:      16,
	})
}

func addTestPets() {
	models.DB.Save(&models.Pet{
		Name:    "Piku",
		Age:     2,
		Kind:    "Dog",
		Specie:  "Poddle",
		OwnerID: 1,
	})
	models.DB.Save(&models.Pet{
		Name:    "Felipe",
		Age:     6,
		Kind:    "Dog",
		Specie:  "Labrador",
		OwnerID: 2,
	})
	models.DB.Save(&models.Pet{
		Name:   "Pepa",
		Age:    5,
		Kind:   "Cat",
		Specie: "Egyptian",
	})
}

func assertPetFound(t *testing.T, wantPet *Pet, id uint) {
	got := makeFindPetRequest(t, id, http.StatusOK)

	wantResponse := &FindPetResponse{Data: wantPet}

	assert.Equal(t, wantResponse, got)
}

func assertPetNotFound(t *testing.T, id uint) {
	got := makeFindPetRequest(t, id, http.StatusNotFound)

	wantResponse := &FindPetResponse{
		Error: petNotFound,
	}
	assert.Equal(t, wantResponse, got)
}

func makeFindPetRequest(t *testing.T, id uint, expectedStatusCode int) *FindPetResponse {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: fmt.Sprintf("%d", id)})
	FindPet(c)
	assert.Equal(t, expectedStatusCode, w.Code)

	var got FindPetResponse
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	return &got
}

func makeAdoptRequest(t *testing.T, ownerID uint, petID uint, expectedStatusCode int) *AdoptResponse {
	adoptReq := AdoptRequest{
		OwnerID: ownerID,
		PetID:   petID,
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	jsonData, err := json.Marshal(adoptReq)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PUT", "", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	c.Request = req
	Adopt(c)
	assert.Equal(t, expectedStatusCode, w.Code)
	var response AdoptResponse
	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatal(err)
	}

	return &response
}

func pepa(available bool) *Pet {
	return &Pet{
		ID:        3,
		Name:      "Pepa",
		Kind:      "Cat",
		Age:       uint(5),
		Specie:    "Egyptian",
		Available: available,
	}
}
