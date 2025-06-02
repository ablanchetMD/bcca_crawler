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


func HandleGetPrecautions(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	items, err := c.Db.GetPrecautionWithProtocols(ctx)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting precautions")
		return
	}
	
	returned_value, err := api.MapAllWithError(items,MapPrecautionWithProtocols)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting precautions")
		fmt.Println("Error:", err)		
		return
	}	

	json_utils.RespondWithJSON(w, http.StatusOK, returned_value)
}


func HandleGetPrecautionByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	item, err := c.Db.GetPrecautionByIDWithProtocols(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting precaution: %s", ids.ID.String()))
		return
	}

	returned_value,err := MapPrecautionWithProtocols(item)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting precautions")
		fmt.Println("Error:", err)		
		return
	}	
	
	json_utils.RespondWithJSON(w, http.StatusOK, returned_value)
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
		ID: pid,
		Title: req.Title,		
		Description: req.Description,		
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
	
	var req ChangeProtocolReq
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
