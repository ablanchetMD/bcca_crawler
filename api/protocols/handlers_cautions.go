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


func HandleGetCautions(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()		
	items, err := c.Db.GetCautionWithProtocols(ctx)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting cautions")
		fmt.Println("Error:", err)
		return
	}

	returned_value, err := api.MapAllWithError(items,MapCautionWithProtocols)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting cautions")
		fmt.Println("Error:", err)		
		return
	}	

	json_utils.RespondWithJSON(w, http.StatusOK, returned_value)
}


func HandleGetCautionsByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	item, err := c.Db.GetCautionByIDWithProtocols(ctx, ids.ID)

	if err != nil {
		println(err.Error())
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting cautions: %s", ids.ID.String()))
		return
	}	
	
	returned_value, err := MapCautionWithProtocols(item)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	json_utils.RespondWithJSON(w, http.StatusOK, returned_value)
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
		ID: pid,		
		Description: req.Description,		
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

	err = c.Db.UpdateCautionProtocols(ctx, database.UpdateCautionProtocolsParams{		
		CautionID: ids.ID,
		ProtocolIds: selectedUUIDs,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating cautions for protocol: %s", ids.ProtocolID))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "cautions updated for protocol"})
}
