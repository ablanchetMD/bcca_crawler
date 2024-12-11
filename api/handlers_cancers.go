package api

import (
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"
	"bcca_crawler/internal/json_utils"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Cancer struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	TumorGroup string    `json:"tumor_group"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	Tags       []string  `json:"tags"`
	Notes      string    `json:"notes"`
}

type CancerRequest struct {
	Code       string   `json:"code" validate:"omitempty,min=1,max=10"`
	Name       string   `json:"name" validate:"required,min=1,max=50"`
	TumorGroup string   `json:"tumor_group" validate:"omitempty,tumorgroup"`
	Tags       []string `json:"tags" validate:"omitempty,max=10,dive,min=1,max=50"`
	Notes      string   `json:"notes" validate:"omitempty,max=500"`
}

type Cancers struct {
	Cancers []Cancer `json:"cancers"`
}

func mapCancerStruct(src database.Cancer) Cancer {
	return Cancer{
		ID:         src.ID,
		CreatedAt:  src.CreatedAt,
		UpdatedAt:  src.UpdatedAt,
		TumorGroup: src.TumorGroup,
		Code:       src.Code.String,
		Name:       src.Name.String,
		Tags:       src.Tags,
		Notes:      src.Notes,
	}
}

func HandleGetCancers(c *config.Config, q QueryParams, w http.ResponseWriter, r *http.Request) {
	var payload []database.Cancer
	var err error
	params := database.GetCancersParams{
		Limit:  int32(q.Limit),
		Offset: int32(q.Offset),
	}

	//optional queries : sort, sort_by, page, limit, offset, filter, fields, include, exclude,
	switch {
	case len(q.FilterBy) > 0 && len(q.Include) > 0:
		payload, err = c.Db.GetCancersOnlyTumorGroupAndTagsAsc(r.Context(), database.GetCancersOnlyTumorGroupAndTagsAscParams{
			TumorGroup: q.FilterBy,
			Tags:       q.Include,
			Limit:      params.Limit,
			Offset:     params.Offset,
		})
	case len(q.FilterBy) > 0 && len(q.Include) == 0:
		payload, err = c.Db.GetCancersOnlyTumorGroupAsc(r.Context(), database.GetCancersOnlyTumorGroupAscParams{
			TumorGroup: q.FilterBy,
			Limit:      params.Limit,
			Offset:     params.Offset,
		})
	default:
		payload, err = c.Db.GetCancers(r.Context(), params)
	}

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching protocols")
		return
	}

	var response []Cancer
	for _, cancer := range payload {
		response = append(response, mapCancerStruct(cancer))
	}

	json_utils.RespondWithJSON(w, http.StatusOK, response)
}

func HandleGetCancerById(c *config.Config, w http.ResponseWriter, r *http.Request) {

	parsed_id, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cancer, err := c.Db.GetCancerByID(r.Context(), parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting cancer: %s", parsed_id.String()))
		return
	}
	json_utils.RespondWithJSON(w, http.StatusOK, mapCancerStruct(cancer))
}

func HandleUpdateCancer(c *config.Config, w http.ResponseWriter, r *http.Request) {

	parsed_id, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req CancerRequest
	err = UnmarshalAndValidatePayload(c, r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cancer, err := c.Db.UpdateCancer(r.Context(), database.UpdateCancerParams{
		ID:         parsed_id,
		TumorGroup: req.TumorGroup,
		Code:       ToNullString(req.Code),
		Name:       ToNullString(req.Name),
		Tags:       req.Tags,
		Notes:      req.Notes,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating protocol: %s", parsed_id.String()))
		return
	}
	json_utils.RespondWithJSON(w, http.StatusOK, mapCancerStruct(cancer))
}

func HandleCreateCancer(c *config.Config, w http.ResponseWriter, r *http.Request) {

	var req CancerRequest
	err := UnmarshalAndValidatePayload(c, r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cancer, err := c.Db.CreateCancer(r.Context(), database.CreateCancerParams{
		TumorGroup: req.TumorGroup,
		Code:       ToNullString(req.Code),
		Name:       ToNullString(req.Name),
		Tags:       req.Tags,
		Notes:      req.Notes,
	})
	if err != nil {

		fmt.Println("Error creating cancer: ", err)
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error creating cancer")
		return
	}

	json_utils.RespondWithJSON(w, http.StatusCreated, mapCancerStruct(cancer))
}

func HandleDeleteCancer(c *config.Config, w http.ResponseWriter, r *http.Request) {

	parsed_id, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteCancer(r.Context(), parsed_id)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting cancer: %s", parsed_id.String()))
		return
	}
	json_utils.RespondWithJSON(w, http.StatusOK, nil)
}

func HandleGetProtocolsByCancerId(c *config.Config, w http.ResponseWriter, r *http.Request) {

	parsed_id, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	protocols, err := c.Db.GetProtocolsForCancer(r.Context(), parsed_id)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting protocols for cancer: %s", parsed_id.String()))
		return
	}

	var response []Protocol
	for _, protocol := range protocols {
		response = append(response, mapProtocolStruct(protocol))
	}
	json_utils.RespondWithJSON(w, http.StatusOK, response)
}
