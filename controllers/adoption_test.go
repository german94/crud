package controllers

import (
	"net/http"
	"testing"
)

func TestAdopt_Ok(t *testing.T) {
	setUpDB()
	pepa := pepa(true)
	response := makeAdoptRequest(t, 2, pepa.ID, http.StatusCreated)
	pepa.Available = false
	assertPetFound(t, pepa, response.Data.ID)
}

func TestAdopt_Error(t *testing.T) {
	cases := []struct {
		testName       string
		ownerID        uint
		petID          uint
		expectedStatus int
		expectedErr    string
	}{
		{
			testName:       "OwnerDoesNotExist",
			ownerID:        10,
			petID:          pepa(true).ID,
			expectedErr:    personNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			testName:       "PetDoesNotExist",
			ownerID:        2,
			petID:          10,
			expectedErr:    petNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			testName:       "PetAlreadyHasOwner",
			ownerID:        2,
			petID:          1,
			expectedErr:    petNotAvailable,
			expectedStatus: http.StatusConflict,
		},
	}

	for _, testCase := range cases {
		setUpDB()
		response := makeAdoptRequest(t, testCase.ownerID, testCase.petID, testCase.expectedStatus)
		if testCase.expectedErr != response.Error {
			t.Errorf("Error for test %s: expected error %s, actual %s",
				testCase.testName, testCase.expectedErr, response.Error)
		}
	}
}
