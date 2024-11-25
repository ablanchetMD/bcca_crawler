package main

import (	
	"net/http"
	"fmt"
	"bcca_crawler/internal/database"
	"github.com/google/uuid"
	"time"
	"encoding/json"
	"io"
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

type Protocols struct {
	Protocols []Protocol `json:"protocols"`
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

	var requestData map[string]interface{}
	err = json.Unmarshal(body, &requestData)
	if err != nil {		
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	fmt.Println("Request data: ", requestData)

	tumor_group, ok := requestData["tumor_group"]
	if !ok {		
		respondWithError(w, http.StatusBadRequest, "Missing tumor_group field")
		return
	}
	
	code, ok := requestData["code"]
	if !ok {		
		respondWithError(w, http.StatusBadRequest, "Missing code field")
		return
	}
	
	name, ok := requestData["name"]
	if !ok {		
		respondWithError(w, http.StatusBadRequest, "Missing name field")
		return
	}

	tags, ok := requestData["tags"]
	if !ok {		
		respondWithError(w, http.StatusBadRequest, "Missing tags field")
		return
	}

	notes, ok := requestData["notes"]
	if !ok {		
		respondWithError(w, http.StatusBadRequest, "Missing notes field")
		return
	}		
	
	protocol, err := c.Db.CreateProtocol(r.Context(), database.CreateProtocolParams{			
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		TumorGroup: tumor_group.(string),
		Code: code.(string),
		Name: name.(string),
		Tags: tags.([]string),
		Notes: notes.(string),		
	})
	if err != nil {
		
		fmt.Println("Error creating chirp: ", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp")
		return
	}
	
	// user.Password = nil
	respondWithJSON(w, http.StatusCreated, mapProtocolStruct(protocol))
}
