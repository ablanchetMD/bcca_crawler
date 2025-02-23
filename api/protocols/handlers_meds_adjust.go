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

type MedModReq struct {
	ID 			string `json:"id" validate:"omitempty,uuid"`
	Category 		string `json:"category" validate:"required,min=1,max=250"`
	Subcategory string `json:"subcategory" validate:"omitempty,min=1,max=500"`
	Adjustment 	string `json:"adjustment" validate:"omitempty,min=1,max=500"`
	MedicationID string `json:"medication_id" validate:"omitempty,uuid"`

}

func HandleGetMedModificationsByProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()	

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw_med_modifications, err := c.Db.GetMedicationModificationsByProtocol(ctx, parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting Med Modifications")
		return
	}
	med_mods := api.MapToMedicationModifications(raw_med_modifications)
	
	json_utils.RespondWithJSON(w, http.StatusOK, med_mods)
}

func HandleGetMedModificationsByMedication(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw_med_modifications, err := c.Db.GetMedicationModificationsByMedication(ctx, parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting Med Modifications")
		return
	}
	medModReqs := []MedModReq{}

	for _, mod := range raw_med_modifications {
		medModReqs = append(medModReqs, MedModReq{
			ID: mod.ModificationID.String(),
			Category: mod.ModificationCategory,
			Subcategory: mod.ModificationSubcategory,
			Adjustment: mod.Adjustment,
			MedicationID: mod.MedicationID.String(),
		})
	}	
	json_utils.RespondWithJSON(w, http.StatusOK, medModReqs)
}

func HandleGetMedModByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw_med_mod, err := c.Db.GetMedicationModificationByID(ctx, parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting Med Mod: %s", parsed_id.String()))
		return
	}
	
	Med := api.MapMedModification(raw_med_mod)

	json_utils.RespondWithJSON(w, http.StatusOK, Med)
}

func HandleDeleteMedModByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	
	ctx := r.Context()

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.RemoveMedicationModification(ctx, parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting Med Mod: %s", parsed_id.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "medication modification deleted"})	
}

func HandlerUpsertMedMod(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	medmodreq := MedModReq{}
	
	err := api.UnmarshalAndValidatePayload(c,r, &medmodreq)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}	

	id, err:= uuid.Parse(medmodreq.ID)
	if err != nil {
		id = uuid.New()
	}
	
	mid, err := uuid.Parse(medmodreq.MedicationID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, "Invalid Medication ID")
		return
	}

	med,err := c.Db.GetMedicationByID(ctx,mid)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Invalid medification associated with modification: %s", err.Error()))
		return
	}


	medication_modification := database.UpsertMedicationModificationParams{
		ID: id,
		Category: medmodreq.Category,
		Subcategory: medmodreq.Subcategory,
		Adjustment: medmodreq.Adjustment,
		MedicationID: mid,
	}

	return_medmod, err := c.Db.UpsertMedicationModification(ctx, medication_modification)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating or updating modify med: %s", err.Error()))
		return
	}

	myReturn := &struct{
		ID string 				`json:"id"`
		MedicationID string 	`json:"medication_id"`
		MedicationName string 	`json:"medication_name"`
		Category string 		`json:"category"`
		Subcategory string 		`json:"subcategory"`
		Adjustment string		`json:"adjustment"`
	}{
		ID: return_medmod.ID.String(),
		MedicationID: return_medmod.MedicationID.String(),
		MedicationName: med.Name,
		Category: return_medmod.Category,
		Subcategory: return_medmod.Subcategory,
		Adjustment: return_medmod.Adjustment,
	}
	json_utils.RespondWithJSON(w, http.StatusCreated, myReturn)

}