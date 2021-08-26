package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestListPets(t *testing.T) {
	setUpDB()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	ListPets(c)
	assert.Equal(t, 200, w.Code)

	var got ListPetsResponse
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := ListPetsResponse{
		Data: []Pet{
			{
				ID:        1,
				Name:      "Piku",
				Age:       2,
				Kind:      "Dog",
				Specie:    "Poddle",
				Available: false,
			},
			{
				ID:        2,
				Name:      "Felipe",
				Age:       6,
				Kind:      "Dog",
				Specie:    "Labrador",
				Available: false,
			},
			{
				ID:        3,
				Name:      "Pepa",
				Age:       5,
				Kind:      "Cat",
				Specie:    "Egyptian",
				Available: true,
			},
		},
	}
	assert.Equal(t, want, got)
}

func TestFindPet(t *testing.T) {
	setUpDB()

	want := &Pet{
		ID:        1,
		Name:      "Piku",
		Kind:      "Dog",
		Specie:    "Poddle",
		Age:       2,
		Available: false,
	}
	assertPetFound(t, want, 1)

}

func TestFindPet_NotFound(t *testing.T) {
	setUpDB()

	assertPetNotFound(t, 50)
}

func TestCreatePet(t *testing.T) {
	setUpDB()

	w := httptest.NewRecorder()
	newPet := CreatePetRequest{
		Name:   "Chicho",
		Kind:   "Dog",
		Specie: "Streeter",
		Age:    18,
	}
	c, _ := gin.CreateTestContext(w)
	jsonData, err := json.Marshal(newPet)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	c.Request = req
	CreatePet(c)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response CreatePetResponse
	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Chicho", response.Data.Name)
	assert.Equal(t, "Dog", response.Data.Kind)
	assert.Equal(t, uint(18), response.Data.Age)
	assert.Equal(t, "Streeter", response.Data.Specie)
	assert.Equal(t, true, response.Data.Available)
	assertPetFound(t, &Pet{
		ID:        response.Data.ID,
		Name:      "Chicho",
		Kind:      "Dog",
		Age:       uint(18),
		Specie:    "Streeter",
		Available: true,
	}, response.Data.ID)
}

func TestDeletePet(t *testing.T) {
	setUpDB()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})
	DeletePet(c)
	assert.Equal(t, 200, w.Code)

	var got FindPetResponse
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assertPetNotFound(t, 1)
}
