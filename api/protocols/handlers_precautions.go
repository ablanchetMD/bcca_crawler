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


type PrecautionReq struct {
	ID 			string `json:"id" validate:"omitempty,uuid"`
	Title 		string `json:"title" validate:"required,min=1,max=250"`
	Description string `json:"description" validate:"required,min=1,max=500"`
	ProtocolID	string `json:"protocol_id" validate:"omitempty,uuid"`
}


type PrecautionResp struct {
	ID 			string `json:"id"`
	Title 		string `json:"title"`
	Description string `json:"description"`
	CreateAt 	string `json:"created_at"`
	UpdateAt 	string `json:"updated_at"`	
	LinkedProtocols []api.LinkedProtocols `json:"linked_protocols"`
}

type PrecautionUpdateReq struct {
	SelectedProtocolIDs []string `json:"protocol_ids"`
}



func HandleGetPrecautions(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	precautions := []PrecautionResp{}
	raw_cautions, err := c.Db.GetPrecautionWithProtocols(ctx)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting cautions")
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
		precautions = append(precautions, PrecautionResp{
			ID:          a.ID.String(),
			Title:       a.Title,			
			Description:     a.Description,
			CreateAt:    a.CreatedAt.Format(`"2006-01-02 15:04:05 MST"`),
			UpdateAt:    a.UpdatedAt.Format(`"2006-01-02 15:04:05 MST"`),			
			LinkedProtocols: linkedProtocols,		
		})
	}

	json_utils.RespondWithJSON(w, http.StatusOK, precautions)
}


func HandleGetPrecautionByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw_caution, err := c.Db.GetPrecautionByIDWithProtocols(ctx, ids.ID)

	if err != nil {
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

	Precaution := PrecautionResp{
		ID: raw_caution.ID.String(),
		Title: raw_caution.Title,
		Description: raw_caution.Description,
		CreateAt: raw_caution.CreatedAt.Format(`"2006-01-02 15:04:05 MST"`),
		UpdateAt: raw_caution.UpdatedAt.Format(`"2006-01-02 15:04:05 MST"`),
		LinkedProtocols: linkedProtocols,
	}
	
	json_utils.RespondWithJSON(w, http.StatusOK, Precaution)
}

func HandleDeletePrecautionByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteProtocolPrecaution(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting precaution: %s", ids.ID.String()))
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
		Column1: pid,
		Column2: req.Title,		
		Column3: req.Description,		
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upserting precaution: %s", req.ID))
		return		
	}

	if req.ProtocolID == "" {
		json_utils.RespondWithJSON(w, http.StatusOK, api.MapPrecaution(caution))
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

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}	
	
	err = c.Db.AddProtocolPrecautionToProtocol(ctx, database.AddProtocolPrecautionToProtocolParams{
		PrecautionID: ids.ID,
		ProtocolID: ids.ProtocolID,
	})	

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding precaution to protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "precaution added to protocol"})

}

func HandleRemovePrecautionFromProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}	
	
	err = c.Db.RemoveProtocolPrecautionFromProtocol(ctx, database.RemoveProtocolPrecautionFromProtocolParams{
		PrecautionID: ids.ID,
		ProtocolID: ids.ProtocolID,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing eligibility criteria from protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "caution removed from protocol"})

}

func HandleUpdatePrecautionsToProtocols(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	var req PrecautionUpdateReq
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

	err = c.Db.UpdatePrecautionProtocols(ctx, database.UpdatePrecautionProtocolsParams{		
		PrecautionID: ids.ID,
		ProtocolIds: selectedUUIDs,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating precautions for protocol: %s", ids.ProtocolID))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "precautions updated for protocol"})
}
