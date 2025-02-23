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

type MedReq struct {
	ID 			string `json:"id" validate:"omitempty,uuid"`
	Name 		string `json:"name" validate:"required,min=1,max=250"`
	Description string `json:"description" validate:"omitempty,min=1,max=500"`
	Category 	string `json:"category" validate:"omitempty,min=1,max=50"`		
}

type PrescriptionReq struct {
	ID 				string `json:"id" validate:"omitempty,uuid"`	
	MedicationID 	string `json:"medication_id" validate:"required,uuid"`
	Dose 			string `json:"dose" validate:"required"`
	Route 			string `json:"route" validate:"required,prescription_route"`
	Frequency 		string `json:"frequency" validate:"required"`
	Duration 		string `json:"duration" validate:"omitempty"`
	Instructions 	string `json:"instructions" validate:"omitempty,min=1,max=1000"`
	Renewals 		int32 	`json:"renewals" validate:"omitempty,min=0,max=50"`
}

func (e *PrescriptionReq) ToRouteEnum() database.PrescriptionRouteEnum {
	return database.PrescriptionRouteEnum(strings.ToLower(e.Route))
}


func HandleGetMeds(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	meds := []api.MedicationResp{}
	raw_meds, err := c.Db.GetMedications(ctx)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting Meds")
		return
	}
	
	for _, a := range raw_meds {		
		meds = append(meds, api.MapMedication(a))		
	}

	json_utils.RespondWithJSON(w, http.StatusOK, meds)
}

func HandleGetMedByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw_med, err := c.Db.GetMedicationByID(ctx, parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting Med: %s", parsed_id.String()))
		return
	}
	
	Med := api.MapMedication(raw_med)

	json_utils.RespondWithJSON(w, http.StatusOK, Med)
}

func HandleDeleteMedByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	
	ctx := r.Context()

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteMedication(ctx, parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting Med: %s", parsed_id.String()))
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
		pid = uuid.New()
	}			

	medication := database.UpsertMedicationParams{
		ID: pid,
		Name: med.Name,
		Description: med.Description,
		Category: med.Category,
	}

	return_med, err := c.Db.UpsertMedication(ctx, medication)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating or updating med: %s", err.Error()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusCreated, api.MapMedication(return_med))
}

func HandleGetPrescriptions(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	prescriptions := []api.PrescriptionResp{}
	raw_prescriptions, err := c.Db.GetPrescriptions(ctx)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting Prescriptions")
		return
	}
	
	for _, a := range raw_prescriptions {		
		prescriptions = append(prescriptions, api.MapPrescription(a))		
	}

	json_utils.RespondWithJSON(w, http.StatusOK, prescriptions)
}

func HandleGetPrescriptionByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	
	ctx := r.Context()

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw_prescription, err := c.Db.GetPrescriptionByID(ctx, parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting Prescription: %s", parsed_id.String()))
		return
	}

	Prescription := api.MapPrescriptionsByID(raw_prescription)

	json_utils.RespondWithJSON(w, http.StatusOK, Prescription)
}

func HandleDeletePrescriptionByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	parsed_id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.RemovePrescription(ctx, parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting Prescription: %s", parsed_id.String()))
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
		Medication: med_id,
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

	return_med, err := c.Db.GetMedicationByID(ctx, return_px.Medication)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting Med: %s", return_px.Medication.String()))
		return
	}

	formated_px := api.PrescriptionResp{
		ID: return_px.ID.String(),
		MedicationName: return_med.Name,
		MedicationID: return_px.Medication.String(),
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
	
	category := r.URL.Query().Get("prescription_category")

	err = c.Validate.Var(category, "required, protocol_prescription_category")

	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid prescription_category: %s", category))
		return
	}

	err = c.Db.AddPrescriptionToProtocolByCategory(ctx, database.AddPrescriptionToProtocolByCategoryParams{
		ProtocolID: parsed_pid,
		PrescriptionID: parsed_id,
		Category: database.MedProtoCategoryEnum(strings.ToLower(category)),
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding prescription to protocol: %s", parsed_id.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "prescription added to protocol"})

}

func HandleRemovePrescriptionFromProtocolByCategory(c *config.Config, w http.ResponseWriter, r *http.Request) {
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
	
	category := r.URL.Query().Get("prescription_category")

	err = c.Validate.Var(category, "required, protocol_prescription_category")

	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid prescription_category: %s", category))
		return
	}

	err = c.Db.RemovePrescriptionFromProtocolByCategory(ctx, database.RemovePrescriptionFromProtocolByCategoryParams{
		ProtocolID: parsed_pid,
		PrescriptionID: parsed_id,
		Category: database.MedProtoCategoryEnum(strings.ToLower(category)),
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing prescription from protocol: %s", parsed_id.String()))
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

	px,err := c.Db.GetPrescriptionsByProtocolByCategory(ctx, database.GetPrescriptionsByProtocolByCategoryParams{
		ProtocolID: parsed_pid,
		Category: database.MedProtoCategoryEnum(strings.ToLower(category)),
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting prescriptions by category: %s", category))
		return
	}

	prescriptions := []api.PrescriptionResp{}

	for _, a := range px {		

		prescriptions = append(prescriptions, api.PrescriptionResp{
			ID: a.MedicationPrescriptionID.String(),
			MedicationName: a.Name,
			MedicationID: a.MedicationID.String(),
			Dose: a.Dose,
			Route: string(a.Route),
			Frequency: a.Frequency,
			Duration: a.Duration,
			Instructions: a.Instructions,
			Renewals: a.Renewals,
		})
	}

	json_utils.RespondWithJSON(w, http.StatusOK, prescriptions)

}





