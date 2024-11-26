package main

import (	
	"net/http"
	"fmt"
	"bcca_crawler/internal/database"
	"github.com/google/uuid"
	"time"
	"encoding/json"
	"io"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Protocol struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TumorGroup      string `json:"tumor_group"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	Tags    []string `json:"tags"`
	Notes   string `json:"notes"`
}

type ProtocolRequest struct {
	TumorGroup      string `json:"tumor_group" validate:"required,tumorgroup"`
	Code    string `json:"code" validate:"required,min=1,max=10"`
	Name    string `json:"name" validate:"required,min=1,max=50"`
	Tags    []string `json:"tags" validate:"omitempty,max=10,dive,min=1,max=50"`
	Notes   string `json:"notes" validate:"omitempty,max=500"`
}

type Protocols struct {
	Protocols []Protocol `json:"protocols"`
}

// Predefined list of valid tumor group codes
var validTumorGroups = map[string]bool{
	"lymphoma&myeloma": true,
	"leukemia&bmt": true,
	"breast": true,
	"gastrointestinal": true,
	"genitourinary": true,
	"gynecology": true,
	"head&neck": true,
	"lung": true,
	"melanoma": true,
	"neuro-oncology": true,
	"sarcoma": true,	
	// Add more as needed
}

// Custom validation function
func tumorGroupValidator(fl validator.FieldLevel) bool {
	tumorGroup := strings.ToLower(fl.Field().String()) // Ensure case-insensitivity
	return validTumorGroups[tumorGroup]
}


func mapProtocolStruct(src database.Protocol) Protocol {
	return Protocol{
		ID:        src.ID,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		TumorGroup:      src.TumorGroup,
		Code:    src.Code,
		Name:    src.Name,
		Tags:    src.Tags,
		Notes:   src.Notes,
	}
}

func handleGetProtocols(c *ApiConfig, w http.ResponseWriter, r *http.Request) {
	
	protocols, err := c.Db.GetProtocols(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching protocols")
		return
	}
	
	var response []Protocol
	for _, protocol := range protocols {
		response = append(response, mapProtocolStruct(protocol))
	}
	
	respondWithJSON(w, http.StatusOK, response)
}


func handleCreateProtocol(c *ApiConfig, w http.ResponseWriter, r *http.Request) {
	
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		respondWithError(w, http.StatusBadRequest, "No body in request")
		return
	}
	defer r.Body.Close()

	var req ProtocolRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate the request data
	err = validate.Struct(req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}		
	
	protocol, err := c.Db.CreateProtocol(r.Context(), database.CreateProtocolParams{			
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		TumorGroup: req.TumorGroup,
		Code: req.Code,
		Name: req.Name,
		Tags: req.Tags,
		Notes: req.Notes,		
	})
	if err != nil {
		
		fmt.Println("Error creating protocol: ", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating protocol")
		return
	}
	
	// user.Password = nil
	respondWithJSON(w, http.StatusCreated, mapProtocolStruct(protocol))
}
