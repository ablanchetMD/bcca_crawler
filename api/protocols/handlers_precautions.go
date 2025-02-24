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


type PrecautionReq struct {
	ID 			string `json:"id" validate:"omitempty,uuid"`
	Title 		string `json:"title" validate:"required,min=1,max=250"`
	Description string `json:"description" validate:"required,min=1,max=500"`
	ProtocolID	string `json:"protocol_id" validate:"omitempty,uuid"`
}


type PrecausionResp struct {
	ID 			string `json:"id"`
	Title 		string `json:"title"`
	Description string `json:"description"`	
	LinkedProtocols []api.LinkedProtocols `json:"linked_protocols"`
}



func HandleGetPrecautions(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	precautions := []PrecausionResp{}
	raw_cautions, err := c.Db.GetPrecautionWithProtocols(ctx)

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
		precautions = append(precautions, PrecausionResp{
			ID:          a.ID.String(),
			Title:       a.Title,			
			Description:     a.Description,			
			LinkedProtocols: linkedProtocols,		
		})
	}

	json_utils.RespondWithJSON(w, http.StatusOK, precautions)
}


func HandleGetPrecautionByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw_caution, err := c.Db.GetPrecautionByIDWithProtocols(ctx, parsed_id.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting cautions: %s", parsed_id.ID.String()))
		return
	}
	
	linkedProtocols, err := api.ConvertTuplesToStructs[api.LinkedProtocols](raw_caution.ProtocolIds)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting linked protocols: %s", parsed_id.ID.String()))
		return
	}

	Precaution := PrecausionResp{
		ID: raw_caution.ID.String(),
		Title: raw_caution.Title,
		Description: raw_caution.Description,
		LinkedProtocols: linkedProtocols,
	}
	
	json_utils.RespondWithJSON(w, http.StatusOK, Precaution)
}

func HandleDeletePrecautionByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteProtocolPrecaution(ctx, parsed_id.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting precaution: %s", parsed_id.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "precaution deleted"})
}

func HandleUpsertPrecaution(c *config.Config, w http.ResponseWriter, r *http.Request) {

	var req PrecautionReq	
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
	
	caution,err := c.Db.UpsertPrecaution(ctx,database.UpsertPrecautionParams{
		ID: pid,
		Title: req.Title,		
		Description: req.Description,		
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upserting precaution: %s", req.ID))
		return		
	}

	proto_id, err:= uuid.Parse(req.ProtocolID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding precaution to protocol (invalid UUID): %s", req.ProtocolID))		
	}else{
		err = c.Db.AddProtocolPrecautionToProtocol(ctx,database.AddProtocolPrecautionToProtocolParams{
			ProtocolID: proto_id,			
			PrecautionID: caution.ID,
		})
	
		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding precaution to protocol: %s", req.ProtocolID))			
		}	

	}

	json_utils.RespondWithJSON(w, http.StatusOK, api.MapPrecaution(caution))	
}

func HandleAddPrecautionToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
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
	
	err = c.Db.AddProtocolPrecautionToProtocol(ctx, database.AddProtocolPrecautionToProtocolParams{
		PrecautionID: parsed_id.ID,
		ProtocolID: parsed_pid,
	})	

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding precaution to protocol: %s", parsed_id.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "precaution added to protocol"})

}

func HandleRemovePrecautionFromProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
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
	
	err = c.Db.RemoveProtocolPrecautionFromProtocol(ctx, database.RemoveProtocolPrecautionFromProtocolParams{
		PrecautionID: parsed_id.ID,
		ProtocolID: parsed_pid,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing eligibility criteria from protocol: %s", parsed_id.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "caution removed from protocol"})

}
