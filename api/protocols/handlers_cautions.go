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


type CautionReq struct {
	ID 			string `json:"id" validate:"omitempty,uuid"`	
	Description string `json:"description" validate:"required,min=1,max=500"`
	ProtocolID	string `json:"protocol_id" validate:"omitempty,uuid"`
}

type CautionUpdateReq struct {
	SelectedProtocolIDs []string `json:"protocol_ids"`
}


type CautionResp struct {
	ID 			string `json:"id"`
	CreatedAt 	string `json:"created_at"`
	UpdatedAt 	string `json:"updated_at"`	
	Description string `json:"description"`
	LinkedProtocols []api.LinkedProtocols `json:"linked_protocols"`
}



func HandleGetCautions(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()	
	cautions := []CautionResp{}
	raw_cautions, err := c.Db.GetCautionWithProtocols(ctx)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting cautions")
		fmt.Println("Error:", err)
		return
	}
	
	for _, a := range raw_cautions {
		
		var linkedProtocols []api.LinkedProtocols	
	
		protocolIdsBytes, ok := a.ProtocolIds.([]byte)
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
		cautions = append(cautions, CautionResp{
			ID:          a.ID.String(),		
			Description:     a.Description,
			CreatedAt: a.CreatedAt.String(),
			UpdatedAt: a.UpdatedAt.String(),			
			LinkedProtocols: linkedProtocols,		
		})
	}

	json_utils.RespondWithJSON(w, http.StatusOK, cautions)
}


func HandleGetCautionsByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw_caution, err := c.Db.GetCautionByIDWithProtocols(ctx, ids.ID)

	if err != nil {
		println(err.Error())
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting cautions: %s", ids.ID.String()))
		return
	}	
	
	var linkedProtocols []api.LinkedProtocols	
	
	protocolIdsBytes, ok := raw_caution.ProtocolIds.([]byte)
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

	Caution := CautionResp{
		ID: raw_caution.ID.String(),
		Description: raw_caution.Description,
		CreatedAt: raw_caution.CreatedAt.String(),
		UpdatedAt: raw_caution.UpdatedAt.String(),
		LinkedProtocols: linkedProtocols,
	}
	
	json_utils.RespondWithJSON(w, http.StatusOK, Caution)
}

func HandleDeleteCautionByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteProtocolCaution(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting caution: %s", ids.ID.String()))
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
	if err != nil || req.ID == "" {
		pid = uuid.Nil
	}		
	
	caution,err := c.Db.UpsertCaution(ctx,database.UpsertCautionParams{
		Column1: pid,		
		Column2: req.Description,		
	})

	if err != nil {		
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upserting caution: %s", req.ID))
		return		
	}
	
	if req.ProtocolID == "" {
		json_utils.RespondWithJSON(w, http.StatusOK, api.MapCaution(caution))
		return
	}

	proto_id, err:= uuid.Parse(req.ProtocolID)
	
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding caution to protocol (invalid UUID): %s", req.ProtocolID))		
	}else{
		err = c.Db.AddProtocolCautionToProtocol(ctx,database.AddProtocolCautionToProtocolParams{
			ProtocolID: proto_id,			
			CautionID: caution.ID,
		})
	
		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding caution to protocol: %s", req.ProtocolID))			
		}	

	}

	json_utils.RespondWithJSON(w, http.StatusOK, api.MapCaution(caution))	
}

func HandleAddCautionToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}	
	
	err = c.Db.AddProtocolCautionToProtocol(ctx, database.AddProtocolCautionToProtocolParams{
		CautionID: ids.ID,
		ProtocolID: ids.ProtocolID,
	})	

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding caution to protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "caution added to protocol"})

}

func HandleRemoveCautionFromProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}	
	
	err = c.Db.RemoveProtocolCautionFromProtocol(ctx, database.RemoveProtocolCautionFromProtocolParams{
		CautionID: ids.ID,
		ProtocolID: ids.ProtocolID,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing eligibility criteria from protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "caution removed from protocol"})

}

func HandleUpdateCautionsToProtocols(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	var req CautionUpdateReq
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

	err = c.Db.UpdateCautionProtocols(ctx, database.UpdateCautionProtocolsParams{		
		CautionID: ids.ID,
		Column2: selectedUUIDs,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating cautions for protocol: %s", ids.ProtocolID))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "cautions updated for protocol"})
}
