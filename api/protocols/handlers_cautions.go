package protocols

import (
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"
	"bcca_crawler/internal/json_utils"
	"bcca_crawler/api"
	"fmt"
	"net/http"
	"github.com/google/uuid"	

)


type CautionReq struct {
	ID 			string `json:"id" validate:"omitempty,uuid"`	
	Description string `json:"description" validate:"required,min=1,max=500"`
	ProtocolID	string `json:"protocol_id" validate:"omitempty,uuid"`
}


type CautionResp struct {
	ID 			string `json:"id"`	
	Description string `json:"description"`
	LinkedProtocols []api.LinkedProtocols `json:"linked_protocols"`
}



func HandleGetCautions(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cautions := []CautionResp{}
	raw_cautions, err := c.Db.GetCautionWithProtocols(ctx)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting cautions")
		return
	}
	
	for _, a := range raw_cautions {
		linkedProtocols, err := api.ConvertTuplesToStructs[api.LinkedProtocols](a.ProtocolIds)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}	
		cautions = append(cautions, CautionResp{
			ID:          a.ID.String(),			
			Description:     a.Description,			
			LinkedProtocols: linkedProtocols,		
		})
	}

	json_utils.RespondWithJSON(w, http.StatusOK, cautions)
}


func HandleGetCautionsByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw_caution, err := c.Db.GetCautionByIDWithProtocols(ctx, parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting cautions: %s", parsed_id.String()))
		return
	}
	
	linkedProtocols, err := api.ConvertTuplesToStructs[api.LinkedProtocols](raw_caution.ProtocolIds)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting linked protocols: %s", parsed_id.String()))
		return
	}

	Caution := CautionResp{
		ID: raw_caution.ID.String(),
		Description: raw_caution.Description,
		LinkedProtocols: linkedProtocols,
	}
	
	json_utils.RespondWithJSON(w, http.StatusOK, Caution)
}

func HandleDeleteCautionByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteProtocolCaution(ctx, parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting caution: %s", parsed_id.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "caution deleted"})
}

func HandleUpsertCaution(c *config.Config, w http.ResponseWriter, r *http.Request) {

	var req CautionReq	
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
	
	caution,err := c.Db.UpsertCaution(ctx,database.UpsertCautionParams{
		ID: pid,		
		Description: req.Description,		
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upserting caution: %s", req.ID))
		return		
	}

	proto_id, err:= uuid.Parse(req.ProtocolID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding precaution to protocol (invalid UUID): %s", req.ProtocolID))		
	}else{
		err = c.Db.AddProtocolCautionToProtocol(ctx,database.AddProtocolCautionToProtocolParams{
			ProtocolID: proto_id,			
			CautionID: caution.ID,
		})
	
		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding precaution to protocol: %s", req.ProtocolID))			
		}	

	}

	json_utils.RespondWithJSON(w, http.StatusOK, api.MapCaution(caution))	
}

func HandleAddCautionToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
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
	
	err = c.Db.AddProtocolCautionToProtocol(ctx, database.AddProtocolCautionToProtocolParams{
		CautionID: parsed_id,
		ProtocolID: parsed_pid,
	})	

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding caution to protocol: %s", parsed_id.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "caution added to protocol"})

}

func HandleRemoveCautionFromProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
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
	
	err = c.Db.RemoveProtocolCautionFromProtocol(ctx, database.RemoveProtocolCautionFromProtocolParams{
		CautionID: parsed_id,
		ProtocolID: parsed_pid,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing eligibility criteria from protocol: %s", parsed_id.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "caution removed from protocol"})

}
