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



func HandleGetMeds(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	items, err := c.Db.GetMedications(ctx)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting Meds")
		return
	}
	
	returned_value := api.MapAll(items,MapMedication)

	json_utils.RespondWithJSON(w, http.StatusOK, returned_value)
}

func HandleGetMedByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	item, err := c.Db.GetMedicationByID(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting Med: %s", ids.ID.String()))
		return
	}
	
	returned_value := MapMedication(item)

	json_utils.RespondWithJSON(w, http.StatusOK, returned_value)
}

func HandleDeleteMedByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteMedication(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting Med: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "medication deleted"})	
}

func HandleUpsertMed(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	med := MedReq{}
	
	err := api.UnmarshalAndValidatePayload(c,r, &med)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}	

	pid, err:= uuid.Parse(med.ID)
	if err != nil {
		pid = uuid.Nil
	}

	medication := database.UpsertMedicationParams{
		ID: pid,		
		Name: med.Name,		
		Description: med.Description,		
		Category: med.Category,
		AlternateNames: med.AlternateNames,		
	}

	return_med, err := c.Db.UpsertMedication(ctx, medication)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating or updating med: %s", err.Error()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusCreated, MapMedication(return_med))
}

func HandleGetPrescriptions(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()	

	med_id := r.URL.Query().Get("medication_id")	

	switch med_id {
		case "":
			items, err := c.Db.GetPrescriptions(ctx)
			if err != nil {
				json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting Prescriptions")
				return
			}

			returned_value := api.MapAll(items,MapPrescription)			
			json_utils.RespondWithJSON(w, http.StatusOK, returned_value)
			
		default:
			pmed_id, err := uuid.Parse(med_id)
			if err != nil {
				json_utils.RespondWithError(w, http.StatusBadRequest, "Invalid Medication ID")
				return
			}

			items, err := c.Db.GetPrescriptionsByMed(ctx, pmed_id)
			if err != nil {
				json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting Prescriptions")
				return
			}
			returned_value := api.MapAll(items,MapPrescription)			
			json_utils.RespondWithJSON(w, http.StatusOK, returned_value)		
			
		}
	
}

func HandleGetPrescriptionByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	item, err := c.Db.GetPrescriptionByID(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting Prescription: %s", ids.ID.String()))
		return
	}

	returned_value := MapPrescription(item)

	json_utils.RespondWithJSON(w, http.StatusOK, returned_value)
}

func HandleDeletePrescriptionByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.RemovePrescription(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting Prescription: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Prescription deleted"})	
}

func HandleUpsertPrescription(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	med := PrescriptionReq{}
	
	err := api.UnmarshalAndValidatePayload(c,r, &med)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pid, err:= uuid.Parse(med.ID)
	if err != nil {
		pid = uuid.New()
	}
	
	med_id, err := uuid.Parse(med.MedicationID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, "Invalid Medication ID")
		return
	}
	

	medication := database.UpsertPrescriptionParams{
		ID: pid,
		MedicationID: med_id,
		Dose: med.Dose,
		Route: med.ToRouteEnum(),
		Frequency: med.Frequency,
		Duration: med.Duration,
		Instructions: med.Instructions,
		Renewals: med.Renewals,
	}		
	return_px, err := c.Db.UpsertPrescription(ctx, medication)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating or updating prescription: %s", err.Error()))
		return
	}

	return_med, err := c.Db.GetMedicationByID(ctx, return_px.MedicationID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting Med: %s", return_px.MedicationID.String()))
		return
	}

	formated_px := PrescriptionResp{
		ID: return_px.ID,
		MedicationName: return_med.Name,
		MedicationID: return_px.MedicationID,
		Dose: return_px.Dose,
		Route: string(return_px.Route),
		Frequency: return_px.Frequency,
		Duration: return_px.Duration,
		Instructions: return_px.Instructions,
		Renewals: return_px.Renewals,
	}
	
	json_utils.RespondWithJSON(w, http.StatusCreated, formated_px)
}

func HandleAddPrescriptionToProtocolByCategory(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	ids, err := api.ParseAndValidateID(r)
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
	
	category := r.URL.Query().Get("prescription_category")

	err = c.Validate.Var(category, "required, protocol_prescription_category")

	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid prescription_category: %s", category))
		return
	}

	err = c.Db.AddPrescriptionToProtocolByCategory(ctx, database.AddPrescriptionToProtocolByCategoryParams{
		ProtocolID: parsed_pid,
		PrescriptionID: ids.ID,
		Category: database.MedProtoCategoryEnum(strings.ToLower(category)),
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding prescription to protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "prescription added to protocol"})

}

func HandleRemovePrescriptionFromProtocolByCategory(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ids, err := api.ParseAndValidateID(r)
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
	
	category := r.URL.Query().Get("prescription_category")

	err = c.Validate.Var(category, "required, protocol_prescription_category")

	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid prescription_category: %s", category))
		return
	}

	err = c.Db.RemovePrescriptionFromProtocolByCategory(ctx, database.RemovePrescriptionFromProtocolByCategoryParams{
		ProtocolID: parsed_pid,
		PrescriptionID: ids.ID,
		Category: database.MedProtoCategoryEnum(strings.ToLower(category)),
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing prescription from protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "prescription removed from protocol"})

}

func HandleGetPrescriptionsByCategory(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()	
	
	proto_id := r.URL.Query().Get("protocol_id")

	parsed_pid, err := uuid.Parse(proto_id)
    if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest,"protocol_id is not a valid uuid")
		return       
    }
	
	category := r.URL.Query().Get("prescription_category")

	err = c.Validate.Var(category, "required, protocol_prescription_category")

	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid prescription_category: %s", category))
		return
	}

	items,err := c.Db.GetPrescriptionsByProtocolByCategory(ctx, database.GetPrescriptionsByProtocolByCategoryParams{
		ProtocolID: parsed_pid,
		Category: database.MedProtoCategoryEnum(strings.ToLower(category)),
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting prescriptions by category: %s", category))
		return
	}

	returned_value := api.MapAll(items,MapPrescription)

	json_utils.RespondWithJSON(w, http.StatusOK, returned_value)

}





