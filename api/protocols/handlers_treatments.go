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


type TreatmentReq struct {
	ID 						string `json:"id" validate:"omitempty,uuid"`	
	MedicationID 			string `json:"medication_id" validate:"required,uuid"`
	Dose 					string `json:"dose" validate:"required"`
	Route 					string `json:"route" validate:"required,prescription_route"`
	Frequency 				string `json:"frequency" validate:"required"`
	Duration 				string `json:"duration" validate:"required"`
	AdministrationGuide 	string `json:"administration_guide" validate:"omitempty,min=1,max=1000"`	
}

type CycleReq struct {
	ID 						string `json:"id" validate:"omitempty,uuid"`
	CycleID 				string `json:"cycle_id" validate:"required,uuid"`
	Cycle 					string `json:"cycle" validate:"required"`
	CycleDuration 			string `json:"cycle_duration" validate:"omitempty"`
}


func HandleGetTreatments(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	treatments := []api.Treatment{}
	raw_tx, err := c.Db.GetTreatments(ctx)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting treatments")
		return
	}
	
	for _, tx := range raw_tx {
		treatments = append(treatments, api.MapTreatment(tx))		
	}

	json_utils.RespondWithJSON(w, http.StatusOK, treatments)
}


func HandleGetTreatmentByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw_treatment, err := c.Db.GetProtocolTreatmentByID(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting treatment: %s", ids.ID.String()))
		return
	}	

	Treatment := api.MapTreatment(raw_treatment)
	
	json_utils.RespondWithJSON(w, http.StatusOK, Treatment)
}

func HandleGetTreatmentsByCycleID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}	
	
	raw_tx, err := c.Db.GetTreatmentsByCycle(ctx, ids.CycleID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting treatments for cycle: %s", ids.CycleID.String()))
		return
	}
	
	treatments := []api.Treatment{}
	
	for _, tx := range raw_tx {
		treatments = append(treatments, api.MapTreatment(tx))		
	}

	json_utils.RespondWithJSON(w, http.StatusOK, treatments)
}

func HandleDeleteTreatmentByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.RemoveProtocolTreatment(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting treatments: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "protocol treatment deleted"})
}

func HandleUpsertTreatment(c *config.Config, w http.ResponseWriter, r *http.Request) {

	var req TreatmentReq	
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
	
	mid, err := uuid.Parse(req.MedicationID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, "medication_id is not a valid uuid")
		return
	}
	
	raw_tx,err := c.Db.UpsertProtocolTreatment(ctx,database.UpsertProtocolTreatmentParams{
		ID: pid,
		Medication: mid,
		Dose: req.Dose,
		Route: req.Route,
		Frequency: req.Frequency,
		Duration: req.Duration,
		AdministrationGuide: req.AdministrationGuide,
		
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upserting treatment: %s", req.ID))
		return		
	}

	Treatment := api.MapTreatment(raw_tx)

	json_utils.RespondWithJSON(w, http.StatusOK, Treatment)	
}

func HandleAddTreatmentToCycle(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}	
	
	err = c.Db.AddTreatmentToCycle(ctx, database.AddTreatmentToCycleParams{
		ProtocolTreatmentID: ids.ID,
		ProtocolCyclesID: ids.CycleID,
	})	

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding treatment to cycle: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "treatment added to cycle"})

}

func HandleRemoveTreatmentToCycle(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}	

	err = c.Db.RemoveTreatmentFromCycle(ctx, database.RemoveTreatmentFromCycleParams{
		ProtocolTreatmentID: ids.ID,
		ProtocolCyclesID: ids.CycleID,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing treatment from cycle: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "treatment removed from cycle"})

}

func HandleUpsertTreatmentCycle(c *config.Config, w http.ResponseWriter, r *http.Request) {

	var req CycleReq	
	err := api.UnmarshalAndValidatePayload(c,r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()	

	id, err:= uuid.Parse(req.ID)
	if err != nil {
		id = uuid.New()
	}

	protocol_id := r.URL.Query().Get("protocol_id")

	
	pid, err := uuid.Parse(protocol_id)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, "protocol_id is not a valid uuid")
		return
	}
	
	raw_cyc,err := c.Db.UpsertCycleToProtocol(ctx,database.UpsertCycleToProtocolParams{
		ID: id,
		ProtocolID: pid,
		Cycle: req.Cycle,
		CycleDuration: req.CycleDuration,		
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upserting treatment cycle: %s", req.ID))
		return		
	}

	Cycle := api.MapCycle(raw_cyc)

	json_utils.RespondWithJSON(w, http.StatusOK, Cycle)	
}

func HandleGetCycles(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cycles := []api.ProtocolCycle{}
	raw_cyc, err := c.Db.GetCycles(ctx)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting cycles")
		return
	}
	
	for _, cyc := range raw_cyc {
		cycles = append(cycles, api.MapCycle(cyc))		
	}

	json_utils.RespondWithJSON(w, http.StatusOK, cycles)
}

func HandleGetCyclesByProtocolID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	raw_cyc, err := c.Db.GetCyclesByProtocol(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting cycles for protocol: %s", ids.ID.String()))
		return
	}
	
	cycles := []api.ProtocolCycle{}
	
	for _, cyc := range raw_cyc {
		cycles = append(cycles, api.MapCycle(cyc))		
	}

	json_utils.RespondWithJSON(w, http.StatusOK, cycles)
}


func HandleGetCycleByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw_cyc, err := c.Db.GetCycleByID(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting cycle: %s", ids.ID.String()))
		return
	}	

	Cycle := api.MapCycle(raw_cyc)
	
	json_utils.RespondWithJSON(w, http.StatusOK, Cycle)
}

func HandleDeleteCycleByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.RemoveCycleByID(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting cycle: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "protocol cycle deleted"})
}

