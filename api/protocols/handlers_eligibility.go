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

type EligibilityCriterionResp struct {
	ID 			string `json:"id"`
	Type 		string `json:"type"`
	Description string `json:"description"`
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
			json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting eligibility criteria")
			return
		}
		for _, e := range elig {
			eligibility = append(eligibility, database.GetEligibilityCriteriaByTypeRow{
				ID:          e.ID,
				Type:        e.Type,
				Description: e.Description,
				ProtocolIds: e.ProtocolIds,
			})
		}
	default:		

		err := c.Validate.Var(category, "required,eligibility_criteria")
		if err != nil {
			json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		
		elig, err := c.Db.GetEligibilityCriteriaByType(ctx, category)
		
		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting eligibility criteria")
			return
		}
		eligibility = elig
	}
	var eligibilityCriteria []EligibilityCriterionResp

	for _, ec := range eligibility {
		
		linkedProtocols, err := api.ConvertTuplesToStructs[api.LinkedProtocols](ec.ProtocolIds)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}		
		
		eligibilityCriteria = append(eligibilityCriteria, EligibilityCriterionResp{
			ID: ec.ID.String(),
			Type: string(ec.Type),
			Description: ec.Description,
			LinkedProtocols: linkedProtocols,
		})
	}

	json_utils.RespondWithJSON(w, http.StatusOK, eligibilityCriteria)
}


func HandleGetEligibilityCriteriaByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	elig, err := c.Db.GetEligibilityCriteriaByID(ctx, parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting eligibility criteria: %s", parsed_id.String()))
		return
	}
	var eligibilityCriteria EligibilityCriterionResp

	linkedProtocols, err := api.ConvertTuplesToStructs[api.LinkedProtocols](elig.ProtocolIds)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting linked protocols: %s", parsed_id.String()))
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

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteEligibilityCriteria(ctx, parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting eligibility criteria: %s", parsed_id.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "eligibility criteria deleted"})
}

func HandleUpsertEligibilityCriteria(c *config.Config, w http.ResponseWriter, r *http.Request) {

	var req EligibilityCriteriaReq
	var EligibilityCriteria []api.ProtocolEligibilityCriterion
	err := api.UnmarshalAndValidatePayload(c,r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	if len(req.EligibilityCriteria) == 0 {
		json_utils.RespondWithError(w, http.StatusBadRequest, "no eligibility criteria provided")
		return
	}

	for _, ec := range req.EligibilityCriteria {

		pid, err:= uuid.Parse(ec.ID)
		if err != nil {
			pid = uuid.New()
		}		
		
		elig,err := c.Db.UpsertEligibilityCriteria(ctx,database.UpsertEligibilityCriteriaParams{
			ID: pid,
			Type: ec.ToTypeEnum(),
			Description: ec.Description,
		})

		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upserting eligibility criteria: %s", ec.ID))
			continue
		}

		EligibilityCriteria = append(EligibilityCriteria, api.MapEligibilityCriterion(elig))

		proto_id, err:= uuid.Parse(ec.ProtocolID)
		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding eligibility criteria to protocol (invalid UUID): %s", ec.ProtocolID))
			continue
		}	

		err = c.Db.AddEligibilityToProtocol(ctx,database.AddEligibilityToProtocolParams{
			ProtocolID: proto_id,			
			CriteriaID: elig.ID,
		})

		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding eligibility criteria to protocol: %s", ec.ProtocolID))
			continue
		}
	}	

	json_utils.RespondWithJSON(w, http.StatusOK, EligibilityCriteria)	
}

func HandleAddEligibilityCriteriaToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	proto_id := r.URL.Query().Get("protocol_id")

	parsed_pid, err := uuid.Parse(proto_id)
    if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest,"protocol_id is not a valid uuid")
		return       
    }	
	
	err = c.Db.LinkEligibilityToProtocol(ctx, database.LinkEligibilityToProtocolParams{
		CriteriaID: parsed_id,
		ProtocolID: parsed_pid,
	})
	

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding eligibility criteria to protocol: %s", parsed_id.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "eligibility criteria added to protocol"})

}

func HandleRemoveEligibilityFromProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	proto_id := r.URL.Query().Get("protocol_id")

	parsed_pid, err := uuid.Parse(proto_id)
    if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest,"protocol_id is not a valid uuid")
		return       
    }	
	
	err = c.Db.UnlinkEligibilityFromProtocol(ctx, database.UnlinkEligibilityFromProtocolParams{
		CriteriaID: parsed_id,
		ProtocolID: parsed_pid,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing eligibility criteria from protocol: %s", parsed_id.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "eligibility criteria removed from protocol"})

}
