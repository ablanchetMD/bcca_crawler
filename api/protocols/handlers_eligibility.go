package protocols

import (
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"
	"bcca_crawler/internal/json_utils"
	"bcca_crawler/api"
	"fmt"
	"net/http"
	"github.com/google/uuid"
	"strings"
	"encoding/json"	

)


type EligibilityCriterionReq struct {
	ID 			string `json:"id" validate:"omitempty,uuid"`
	Type 		string `json:"type" validate:"required,eligibility_criteria"`
	Description string `json:"description" validate:"required,min=1,max=500"`
	ProtocolID	string `json:"protocol_id" validate:"omitempty,uuid"`
}

type EligibilityCriteriaReq struct {
	EligibilityCriteria []EligibilityCriterionReq `json:"eligibility_criteria"`
}

type EligibilityUpdateReq struct {
	SelectedProtocolIDs []string `json:"protocol_ids"`
}

type EligibilityCriterionResp struct {
	ID 			string `json:"id"`
	Type 		string `json:"type"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	LinkedProtocols []api.LinkedProtocols `json:"linked_protocols"`
}

func (e *EligibilityCriterionReq) ToTypeEnum() database.EligibilityEnum {
	return database.EligibilityEnum(strings.ToLower(e.Type))
}

type ErrorResponse struct {
	Field   string `json:"field"`   // Field that caused the error
	Message string `json:"message"` // Error message
}


func HandleGetEligibilityCriteria(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	category := r.URL.Query().Get("category")
	
	var eligibility  []database.GetEligibilityCriteriaByTypeRow

	switch category {
	case "":
		elig, err := c.Db.GetElibilityCriteria(ctx)
		
		if err != nil {
			println(err.Error())
			json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting eligibility criteria with no requirements")
			return
		}
		for _, e := range elig {
			eligibility = append(eligibility, database.GetEligibilityCriteriaByTypeRow(e))
		}
	default:		

		err := c.Validate.Var(category, "required,eligibility_criteria")
		if err != nil {
			json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		
		elig, err := c.Db.GetEligibilityCriteriaByType(ctx, category)
		
		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting eligibility criteria with requirements")
			return
		}
		eligibility = elig
	}
	var eligibilityCriteria []EligibilityCriterionResp

	for _, ec := range eligibility {
		
		var linkedProtocols []api.LinkedProtocols	
	
		protocolIdsBytes, ok := ec.ProtocolIds.([]byte)
		if !ok {
			json_utils.RespondWithError(w, http.StatusInternalServerError, "Error asserting protocol IDs to []byte")
			return
		}

		err := json.Unmarshal(protocolIdsBytes, &linkedProtocols)
		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, 
				fmt.Sprintf("Error parsing protocol data: %s", err.Error()))
			return
		}		
		
		eligibilityCriteria = append(eligibilityCriteria, EligibilityCriterionResp{
			ID: ec.ID.String(),
			Type: string(ec.Type),
			Description: ec.Description,
			CreatedAt: ec.CreatedAt.Format(`"2006-01-02 15:04:05 MST"`),
			UpdatedAt: ec.UpdatedAt.Format(`"2006-01-02 15:04:05 MST"`),
			LinkedProtocols: linkedProtocols,
		})
	}

	json_utils.RespondWithJSON(w, http.StatusOK, eligibilityCriteria)
}

func HandleGetEligibilityCriteriaByProtocol(c *config.Config, w http.ResponseWriter, r *http.Request){
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	elig_criterias,err := c.Db.GetEligibilityByProtocol(ctx,ids.ProtocolID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return 
	}

	r_elig_criterias := []EligibilityCriterionResp{}
	for _, ec := range elig_criterias {		
		r_elig_criterias = append(r_elig_criterias, EligibilityCriterionResp{
			ID: ec.ID.String(),
			Type: string(ec.Type),
			Description: ec.Description,
			CreatedAt: ec.CreatedAt.Format(`"2006-01-02 15:04:05 MST"`),
			UpdatedAt: ec.UpdatedAt.Format(`"2006-01-02 15:04:05 MST"`),			
		})
	}

	json_utils.RespondWithJSON(w, http.StatusOK, r_elig_criterias)
}


func HandleGetEligibilityCriteriaByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	elig, err := c.Db.GetEligibilityCriteriaByID(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting eligibility criteria: %s", ids.ID.String()))
		return
	}
	var eligibilityCriteria EligibilityCriterionResp

	var linkedProtocols []api.LinkedProtocols	
	
	protocolIdsBytes, ok := elig.ProtocolIds.([]byte)
	if !ok {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error asserting protocol IDs to []byte")
		return
	}

	err = json.Unmarshal(protocolIdsBytes, &linkedProtocols)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, 
			fmt.Sprintf("Error parsing protocol data: %s", err.Error()))
		return
	}

	eligibilityCriteria = EligibilityCriterionResp{
		ID: elig.ID.String(),
		Type: string(elig.Type),
		Description: elig.Description,
		LinkedProtocols: linkedProtocols,
	}

	json_utils.RespondWithJSON(w, http.StatusOK, eligibilityCriteria)
}

func HandleDeleteEligibilityCriteriaByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteEligibilityCriteria(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting eligibility criteria: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "eligibility criteria deleted"})
}

func HandleUpsertEligibilityCriteria(c *config.Config, w http.ResponseWriter, r *http.Request) {

	var req EligibilityCriterionReq	
	err := api.UnmarshalAndValidatePayload(c,r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()	

	pid, err:= uuid.Parse(req.ID)
	if err != nil || req.ID == "" {
		pid = uuid.Nil
	}		
	
	elig,err := c.Db.UpsertEligibilityCriteria(ctx,database.UpsertEligibilityCriteriaParams{
		Column1: pid,
		Column2: req.ToTypeEnum(),
		Column3: req.Description,
	})

	if err != nil {
		println(err.Error())
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upserting eligibility criteria: %s", elig.ID))
		return			
	}	

	if req.ProtocolID == "" {
		json_utils.RespondWithJSON(w, http.StatusOK, api.MapEligibilityCriterion(elig))
		return
	}
	
	proto_id, err:= uuid.Parse(req.ProtocolID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding eligibility criteria to protocol (invalid UUID): %s", req.ProtocolID))		
	}else{
		err = c.Db.AddEligibilityToProtocol(ctx,database.AddEligibilityToProtocolParams{
			ProtocolID: proto_id,			
			CriteriaID: elig.ID,
		})
	
		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding eligibility criteria to protocol: %s", req.ProtocolID))			
		}	

	}	
	
	json_utils.RespondWithJSON(w, http.StatusOK, api.MapEligibilityCriterion(elig))	
}

func HandleAddEligibilityCriteriaToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	
	err = c.Db.LinkEligibilityToProtocol(ctx, database.LinkEligibilityToProtocolParams{
		CriteriaID: ids.ID,
		ProtocolID: ids.ProtocolID,
	})
	

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding eligibility criteria to protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "eligibility criteria added to protocol"})

}

func HandleRemoveEligibilityFromProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}		
	
	err = c.Db.UnlinkEligibilityFromProtocol(ctx, database.UnlinkEligibilityFromProtocolParams{
		CriteriaID: ids.ID,
		ProtocolID: ids.ProtocolID,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing eligibility criteria from protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "eligibility criteria removed from protocol"})

}

func HandleUpdateEligibilityToProtocols(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	var req EligibilityUpdateReq
	err = api.UnmarshalAndValidatePayload(c, r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var selectedUUIDs []uuid.UUID
	for _, id := range req.SelectedProtocolIDs {
		uid, err := uuid.Parse(id)
		if err != nil {
			json_utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid UUID: %s", id))
			return
		}
		selectedUUIDs = append(selectedUUIDs, uid)
	}

	err = c.Db.UpdateEligibilityProtocols(ctx, database.UpdateEligibilityProtocolsParams{		
		CriteriaID: ids.ID,
		Column2: selectedUUIDs,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating criterias for protocol: %s", ids.ProtocolID))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "criterias updated for protocol"})
}


