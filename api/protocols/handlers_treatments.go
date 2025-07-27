package protocols

import (
	"bcca_crawler/api"
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"
	"context"
	"fmt"
	"net/http"
	"github.com/google/uuid"
)

func HandleGetTreatments(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGet(c, w, r, getTx)
}

func HandleGetTreatmentsByCycleID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGet(c, w, r, getTxsByCycleID)
}

func HandleGetTreatmentByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGet(c, w, r, getTxByID)
}

func HandleDeleteTreatmentByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c, w, r, deleteTxByID)
}

func HandleUpsertTreatment(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleUpsert(c, w, r, upsertTreatment)
}

func HandleAddTreatmentToCycle(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c, w, r, addTxToCycle)
}
func HandleRemoveTreatmentToCycle(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c, w, r, removeTxToCycle)
}

func HandleUpsertTreatmentCycle(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleUpsert(c, w, r, upsertCycle)
}

func HandleGetCycles(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGet(c, w, r, getCyclesByProtocolID)
}

func HandleDeleteCycleByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c, w, r, deleteCycleByID)
}


func getTx(c *config.Config, ctx context.Context, ids api.IDs) ([]api.Treatment, error) {
	items, err := c.Db.GetTreatments(ctx)

	if err != nil {
		return nil, err
	}

	response := api.MapAll(items, api.MapTreatment)

	return response, nil
}

func getTxsByCycleID(c *config.Config, ctx context.Context, ids api.IDs) ([]api.Treatment, error) {
	items, err := c.Db.GetTreatmentsByCycle(ctx, ids.ID)
	if err != nil {
		return nil, err
	}

	response := api.MapAll(items, api.MapTreatment)

	return response, nil
}

func getTxByID(c *config.Config, ctx context.Context, ids api.IDs) (api.Treatment, error) {
	item, err := c.Db.GetProtocolTreatmentByID(ctx, ids.ID)

	if err != nil {
		return api.Treatment{}, err
	}

	response := api.MapTreatment(item)

	return response, nil
}

func deleteTxByID(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.RemoveProtocolTreatment(ctx, ids.ID)

	if err != nil {
		return "", fmt.Errorf("error deleting tx: %s, with error: %v", ids.ID.String(), err)

	}
	return fmt.Sprintf("Tx %s deleted.", ids.ID.String()), nil
}

func upsertTreatment(c *config.Config, ctx context.Context, req TreatmentReq, ids api.IDs) error {
	id := api.ParseOrNilUUID(req.ID)

	mid, err := uuid.Parse(req.MedicationID)
	if err != nil {
		return fmt.Errorf("medication ID: %s is not a valid UUID", req.MedicationID)
	}

	_, err = c.Db.UpsertProtocolTreatment(ctx, database.UpsertProtocolTreatmentParams{
		ID:                  id,
		MedicationID:        mid,
		Dose:                req.Dose,
		Route:               req.ToRouteEnum(),
		Frequency:           req.Frequency,
		Duration:            req.Duration,
		AdministrationGuide: req.AdministrationGuide,
	})

	if err != nil {
		return fmt.Errorf("error upserting treatment: %s with error:%s", req.ID, err.Error())
	}

	return nil
}

func addTxToCycle(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.AddTreatmentToCycle(ctx, database.AddTreatmentToCycleParams{
		ProtocolTreatmentID: ids.ID,
		ProtocolCyclesID:    ids.CycleID,
	})
	if err != nil {
		return "", fmt.Errorf("error adding tx: %s to cycle %s, with error: %v", ids.ID.String(), ids.CycleID.String(), err)

	}
	return fmt.Sprintf("Treatment %s added to Cycle %s", ids.ID.String(), ids.CycleID.String()), nil
}

func removeTxToCycle(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.RemoveTreatmentFromCycle(ctx, database.RemoveTreatmentFromCycleParams{
		ProtocolTreatmentID: ids.ID,
		ProtocolCyclesID:    ids.CycleID,
	})
	if err != nil {
		return "", fmt.Errorf("error removing tx: %s to cycle %s, with error: %v", ids.ID.String(), ids.CycleID.String(), err)

	}
	return fmt.Sprintf("Treatment %s removed from Cycle %s", ids.ID.String(), ids.CycleID.String()), nil
}



func upsertCycle(c *config.Config, ctx context.Context, req CycleReq, ids api.IDs) error {
	id := api.ParseOrNilUUID(req.ID)

	_, err := c.Db.UpsertCycleToProtocol(ctx, database.UpsertCycleToProtocolParams{
		ID:            id,
		ProtocolID:    ids.ProtocolID,
		Cycle:         req.Cycle,
		CycleDuration: req.CycleDuration,
	})
	if err != nil {
		return err
	}

	return nil
}

func getCyclesByProtocolID(c *config.Config, ctx context.Context, ids api.IDs) ([]api.ProtocolCycle, error) {
	items, err := c.Db.GetProtocolCyclesWithTreatments(ctx, ids.ProtocolID)
	if err != nil {
		return nil, err
	}

	response, err := api.ToResponseData[[]api.ProtocolCycle](items)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func deleteCycleByID(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.RemoveCycleByID(ctx, ids.ID)

	if err != nil {
		return "", fmt.Errorf("error deleting cycle: %s, with error: %v", ids.ID.String(), err)

	}
	return fmt.Sprintf("Cycle %s deleted.", ids.ID.String()), nil
}


