package api

import (
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"
	"bcca_crawler/internal/json_utils"	
	"fmt"	
	"net/http"	
	"time"	
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	Name    string `json:"name" validate:"required,min=1,max=250"`
	Tags    []string `json:"tags" validate:"omitempty,max=10,dive,min=1,max=50"`
	Notes   string `json:"notes" validate:"omitempty,max=500"`
	ProtocolUrl string `json:"protocol_url" validate:"omitempty,max=250"`
	PatientHandoutUrl string `json:"patient_handout_url" validate:"omitempty,max=250"`
	RevisedOn string `json:"revised_on" validate:"omitempty,max=25"`
	ActivatedOn string `json:"activated_on" validate:"omitempty,max=25"`

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

func HandleGetProtocols(c *config.Config,q QueryParams, w http.ResponseWriter, r *http.Request) {
	var protocols []database.Protocol
	var err error
	params := database.GetProtocolsAscParams{
		Limit:  int32(q.Limit),
		Offset: int32(q.Offset),
	}

	//optional queries : sort, sort_by, page, limit, offset, filter, fields, include, exclude,
	switch {
	case len(q.FilterBy) > 0 && len(q.Include) > 0:
		protocols, err = c.Db.GetProtocolsOnlyTumorGroupAndTagsAsc(r.Context(), database.GetProtocolsOnlyTumorGroupAndTagsAscParams{
			TumorGroup: q.FilterBy,
			Tags:       q.Include,
			Limit:      params.Limit,
			Offset:     params.Offset,
		})
	case len(q.FilterBy) > 0 && len(q.Include) == 0:
		protocols, err = c.Db.GetProtocolsOnlyTumorGroupAsc(r.Context(), database.GetProtocolsOnlyTumorGroupAscParams{
			TumorGroup: q.FilterBy,
			Limit:      params.Limit,
			Offset:     params.Offset,
		})
	default:
		protocols, err = c.Db.GetProtocolsAsc(r.Context(), params)
	}
				
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching protocols")
		return
	}
	
	var response []Protocol
	for _, protocol := range protocols {
		response = append(response, mapProtocolStruct(protocol))
	}
	
	json_utils.RespondWithJSON(w, http.StatusOK, response)
}

func HandleDeleteProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ids, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteProtocol(r.Context(),ids.ID)
		
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting protocols")
		return
	}
	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Protocol %s deleted", ids.ID.String())})
}

func HandleGetProtocolById(c *config.Config, w http.ResponseWriter, r *http.Request) {
	
	ids, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	protocol,err := c.Db.GetProtocolByID(r.Context(),ids.ID)
		
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting protocol: %s", ids.ID.String()))
		return
	}
	json_utils.RespondWithJSON(w, http.StatusOK, mapProtocolStruct(protocol))
}

func HandleUpdateProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ids, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req ProtocolRequest
	err = UnmarshalAndValidatePayload(c,r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	protocol,err := c.Db.UpdateProtocol(r.Context(),database.UpdateProtocolParams{
		ID: ids.ID,		
		TumorGroup: req.TumorGroup,
		Code: req.Code,
		Name: req.Name,
		Tags: req.Tags,
		Notes: req.Notes,
		ProtocolUrl: req.ProtocolUrl,
		PatientHandoutUrl: req.PatientHandoutUrl,
		RevisedOn: req.RevisedOn,
		ActivatedOn: req.ActivatedOn,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating protocol: %s", ids.ID.String()))
		return
	}
	json_utils.RespondWithJSON(w, http.StatusOK, mapProtocolStruct(protocol))
}

func HandleGetProtocolSummary(c *config.Config, w http.ResponseWriter, r *http.Request) {
	
	ids, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response,err := CMD_GetProtocolBy(c,r.Context(),"id",ids.ID.String())
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting protocol: %s", ids.ID.String()))
		return
	}
	
	json_utils.RespondWithJSON(w, http.StatusOK, response)
}

func HandleGetProtocolSummaryCode(c *config.Config, w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
    if len(code) == 0 {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "No protocol code provided")
		return
    }	

	response,err := CMD_GetProtocolBy(c,r.Context(),"code",code)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting protocol: %s", code))
		return
	}
	
	json_utils.RespondWithJSON(w, http.StatusOK, response)
}

func HandleCreateProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {	
	
	var req ProtocolRequest
	err := UnmarshalAndValidatePayload(c,r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	protocol, err := c.Db.CreateProtocol(r.Context(), database.CreateProtocolParams{	
		
		TumorGroup: req.TumorGroup,
		Code: req.Code,
		Name: req.Name,
		Tags: req.Tags,
		Notes: req.Notes,		
	})
	if err != nil {		
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			// Duplicate key value violation
			json_utils.RespondWithError(w, http.StatusInternalServerError, "Record already exists")
			return			
		}
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error creating protocol")
		return
	}	
	
	json_utils.RespondWithJSON(w, http.StatusCreated, mapProtocolStruct(protocol))
}


