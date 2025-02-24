package protocols

import (
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"
	"bcca_crawler/internal/json_utils"
	"bcca_crawler/api"
	"fmt"
	"net/http"
	"github.com/google/uuid"
	"encoding/json"	

)


type ToxicityReq struct {
	ID 						string `json:"id" validate:"omitempty,uuid"`	
	Title 					string `json:"title" validate:"required"`
	Category 				string `json:"category" validate:"required"`
	Description 			string `json:"description" validate:"omitempty,min=1,max=1000"`
	Grades 					[]ToxicityGradeReq `json:"grades" validate:"required"`
	
}

type ToxicityGradeReq struct {
	ID 						string `json:"id" validate:"omitempty,uuid"`	
	Grade 					string `json:"grade" validate:"required,grade"`
	Description				string `json:"description" validate:"min=1,max=1000"`
}

type ToxModReq struct {
	ID 						string `json:"id" validate:"omitempty,uuid"`
	ToxicityGradeID 		string `json:"toxicity_id" validate:"required,uuid"`
	ProtocolID 				string `json:"protocol_id" validate:"required,uuid"`
	Adjustment 				string `json:"adjustment" validate:"required"`
}

func HandleGetToxicities(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	format := []api.ToxicityWithGrades{}
	raw, err := c.Db.GetToxicitiesWithGrades(ctx)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting toxicities")
		return
	}
	
	for _, obj := range raw {
		fobj,err := api.MapToToxicityWithGrades(obj)
		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, "Error mapping toxicities")
			return
		}
		format = append(format, fobj)
	}

	json_utils.RespondWithJSON(w, http.StatusOK, format)
}


func HandleGetToxicityByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()	

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw, err := c.Db.GetToxicityByID(ctx, parsed_id.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting toxicity by id")
		return
	}
	
	fobj,err := api.MapToToxicityWithGradesOne(raw)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error mapping toxicity by id")
		return
	}	
	
	json_utils.RespondWithJSON(w, http.StatusOK, fobj)
}

func HandleGetToxicitiesWithAdjustmentsByProtocolID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	
	ctx := r.Context()

	protocol_id := r.URL.Query().Get("protocol_id")

	parsed_id, err := uuid.Parse(protocol_id)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest,"protocol_id is not a valid uuid")
		return       
	}	
	
	raw, err := c.Db.GetToxicitiesWithGradesAndAdjustments(ctx, parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting toxicities for the protocol: %s", parsed_id.String()))
		return
	}
	
	format := []api.ToxicityWithGradesAndAdjustments{}
	
	for _, obj := range raw {
		fobj,err := api.MapToToxicityWithGradesAndAdjustments(obj)
		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, "Error mapping toxicities")
			return
		}
		format = append(format, fobj)
	}

	json_utils.RespondWithJSON(w, http.StatusOK, format)
}

// To do

func HandleDeleteToxicityByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.RemoveToxicity(ctx, parsed_id.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting toxicity: %s", parsed_id.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "toxicity deleted"})
}

func HandlerUpsertToxicityWithGrades(c *config.Config, w http.ResponseWriter, r *http.Request) {

	var req ToxicityReq	
	err := api.UnmarshalAndValidatePayload(c,r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()	

	pid, err:= uuid.Parse(req.ID)
	if err != nil {
		pid = uuid.New()
	}

	raw, err := c.Db.UpsertToxicity(ctx, database.UpsertToxicityParams{
		ID: pid,
		Title: req.Title,
		Category: req.Category,
		Description: req.Description,
	})
	
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upserting toxicity: %s", req.ID))
		return
	}	

	gradesJSON, err := json.Marshal(req.Grades)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return		
	}
	
	raw_toxgrades,err := c.Db.UpsertToxicityGrades(ctx, database.UpsertToxicityGradesParams{
		ToxicityID: raw.ID,
		Column2: gradesJSON,
	})		

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upserting toxicity grades: %s", req.ID))
		return		
	}

	toxgrades := []api.ToxicityGrade{}

	for _, obj := range raw_toxgrades {
		fobj := api.MapToToxicityGrade(obj)		
		toxgrades = append(toxgrades, fobj)
	}

	Toxicity := api.ToxicityWithGrades{
		ID: raw.ID,
		Title: raw.Title,
		Category: raw.Category,
		Description: raw.Description,
		Grades: toxgrades,
	}

	json_utils.RespondWithJSON(w, http.StatusOK, Toxicity)	
}

func HandleUpsertAdjustmentsToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req ToxModReq
	err := api.UnmarshalAndValidatePayload(c,r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tox_id, err := uuid.Parse(req.ToxicityGradeID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest,"toxicity_grade_id is not a valid uuid")
		return
	}	

	protocol_id, err := uuid.Parse(req.ProtocolID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest,"protocol_id is not a valid uuid")
		return       
	}
	
	adj_id, err := uuid.Parse(req.ID)
	if err != nil {
		adj_id = uuid.New()
	}
	
	obj,err := c.Db.UpsertToxicityToProtocol(ctx, database.UpsertToxicityToProtocolParams{
		ID: adj_id,
		ToxicityGradeID: tox_id,
		ProtocolID: protocol_id,
		Adjustment: req.Adjustment,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding Toxicity Adjustment to Protocol: %s", adj_id.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK,obj)

}

func HandleRemoveAdjustmentsToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	err = c.Db.RemoveToxicityModification(ctx, id.ID)
		
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing toxicity adjustment from protocol: %s", id.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "toxicity adjustment removed from protocol"})

}