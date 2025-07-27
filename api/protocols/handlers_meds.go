package protocols

import (
	"bcca_crawler/api"
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

type ProtocolMedCategoryReq struct {
	ID       string `json:"id" validate:"omitempty,uuid"`
	Category string `json:"category" validate:"required,min=1,max=500"`
	Comments string `json:"comments" validate:"omitempty"`	
}

func getMeds(c *config.Config, ctx context.Context, ids api.IDs) ([]MedicationResp, error) {
	items, err := c.Db.GetMedications(ctx)

	if err != nil {
		return nil, err
	}

	return api.MapAll(items, MapMedication), nil
}

func getMedByID(c *config.Config, ctx context.Context, ids api.IDs) (MedicationResp, error) {
	item, err := c.Db.GetMedicationByID(ctx, ids.ID)

	if err != nil {
		return MedicationResp{}, err
	}

	response := MapMedication(item)

	return response, nil
}

func deleteMedByID(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.DeleteMedication(ctx, ids.ID)

	if err != nil {
		return "", fmt.Errorf("error deleting med: %s, with error: %v", ids.ID.String(), err)

	}
	return fmt.Sprintf("Med %s deleted.", ids.ID.String()), nil
}

func upsertMed(c *config.Config, ctx context.Context, req MedReq, ids api.IDs) error {
	id := api.ParseOrNilUUID(req.ID)

	if req.AlternateNames == nil {
		req.AlternateNames = []string{}
	}

	medication := database.UpsertMedicationParams{
		ID:             id,
		Name:           req.Name,
		Description:    req.Description,
		Category:       req.Category,
		AlternateNames: req.AlternateNames,
	}

	_, err := c.Db.UpsertMedication(ctx, medication)

	if err != nil {
		return fmt.Errorf("error upserting medication: %s with error:%s", req.ID, err.Error())
	}

	return nil
}

func getPrescriptionsWithQuery(c *config.Config, ctx context.Context, ids api.IDs, query url.Values) ([]PrescriptionResp, error) {

	med_id := query.Get("medication_id")

	switch med_id {
	case "":
		items, err := c.Db.GetPrescriptions(ctx)
		if err != nil {
			return nil, fmt.Errorf("error getting Prescriptions:%s", err.Error())
		}

		return api.MapAll(items, MapPrescription), nil

	default:
		pmed_id, err := uuid.Parse(med_id)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID for medication_id")
		}

		items, err := c.Db.GetPrescriptionsByMed(ctx, pmed_id)
		if err != nil {
			return nil, fmt.Errorf("error getting Prescriptions associated with Medication ID:%s", med_id)
		}
		return api.MapAll(items, MapPrescription), nil

	}

}

func getPxByID(c *config.Config, ctx context.Context, ids api.IDs) (PrescriptionResp, error) {
	item, err := c.Db.GetPrescriptionByID(ctx, ids.ID)

	if err != nil {
		return PrescriptionResp{}, err
	}

	response := MapPrescription(item)

	return response, nil
}

func deletePxByID(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.RemovePrescription(ctx, ids.ID)

	if err != nil {
		return "", fmt.Errorf("error deleting px: %s, with error: %v", ids.ID.String(), err)

	}
	return fmt.Sprintf("Px %s deleted.", ids.ID.String()), nil
}

func upsertPrescription(c *config.Config, ctx context.Context, req PrescriptionReq, ids api.IDs) error {
	id := api.ParseOrNilUUID(req.ID)

	mid, err := uuid.Parse(req.MedicationID)
	if err != nil {
		return fmt.Errorf("medication ID: %s is not a valid UUID", req.MedicationID)
	}

	px := database.UpsertPrescriptionParams{
		ID:           id,
		MedicationID: mid,
		Dose:         req.Dose,
		Route:        req.ToRouteEnum(),
		Frequency:    req.Frequency,
		Duration:     req.Duration,
		Instructions: req.Instructions,
		Renewals:     req.Renewals,
	}

	_, err = c.Db.UpsertPrescription(ctx, px)

	if err != nil {
		return fmt.Errorf("error upserting Px: %s with error:%s", req.ID, err.Error())
	}

	return nil
}

func upsertProtocolMedCategory(c *config.Config, ctx context.Context, req ProtocolMedCategoryReq, ids api.IDs) error {
	
	id := api.ParseOrNilUUID(req.ID)

	_, err := c.Db.UpsertProtoMedCategory(ctx, database.UpsertProtoMedCategoryParams{
		ID:         id,
		ProtocolID: ids.ProtocolID,
		Category:   req.Category,
		Comments:   req.Comments,
	})

	if err != nil {
		return err
	}

	return nil

}

func getPxByProtocolID(c *config.Config, ctx context.Context, ids api.IDs) ([]ProtocolMedGroup, error) {
	items, err := c.Db.GetProtocolPrescriptions(ctx, ids.ProtocolID)
	if err != nil {
		return nil, err
	}

	response, err := api.ToResponseData[[]ProtocolMedGroup](items)
	if err != nil {
		return nil, err
	}
	return response, nil
}


func getProtocolMedCategoryByID(c *config.Config, ctx context.Context, ids api.IDs) (ProtocolMedGroup, error) {
	item, err := c.Db.GetMedCategoryByID(ctx, ids.ID)
	if err != nil {
		return ProtocolMedGroup{}, err
	}

	response, err := api.ToResponseData[ProtocolMedGroup](item)
	if err != nil {
		return ProtocolMedGroup{}, err
	}

	return response, nil
}

func deleteProtocolMedCategoryByID(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.RemoveMedCategoryByID(ctx, ids.ID)

	if err != nil {
		return "", fmt.Errorf("error deleting medication category: %s, with error: %v", ids.ID.String(), err)

	}
	return fmt.Sprintf("Medication Category %s deleted.", ids.ID.String()), nil
}

func addPxToProtocolMedCategory(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.AddPrescriptionToProtocolCategory(ctx, database.AddPrescriptionToProtocolCategoryParams{
		MedicationPrescriptionID: ids.ID,
		ProtocolMedsID:           ids.PxCategoryID,
	})
	if err != nil {
		return "", fmt.Errorf("error adding Px: %s to medication category %s, with error: %v", ids.ID.String(), ids.PxCategoryID.String(), err)

	}
	return fmt.Sprintf("Px %s added to medication category %s", ids.ID.String(), ids.PxCategoryID.String()), nil
}

func removePxFromProtocolMedCategory(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.RemovePrescriptionFromProtocolCategory(ctx, database.RemovePrescriptionFromProtocolCategoryParams{
		MedicationPrescriptionID: ids.ID,
		ProtocolMedsID:           ids.PxCategoryID,
	})
	if err != nil {
		return "", fmt.Errorf("error removing Px: %s to medication category %s, with error: %v", ids.ID.String(), ids.PxCategoryID.String(), err)

	}
	return fmt.Sprintf("Px %s removed from medication category %s", ids.ID.String(), ids.PxCategoryID.String()), nil
}

func HandleGetMeds(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGet(c, w, r, getMeds)
}

func HandleGetMedByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGet(c, w, r, getMedByID)

}

func HandleDeleteMedByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c, w, r, deleteMedByID)
}

func HandleUpsertMed(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleUpsert(c, w, r, upsertMed)

}

func HandleGetPrescriptions(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGetWithQ(c, w, r, getPrescriptionsWithQuery)

}

func HandleGetPrescriptionByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGet(c,w,r,getPxByID)
}

func HandleDeletePrescriptionByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c,w,r,deletePxByID)

}

func HandleUpsertPrescription(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleUpsert(c,w,r,upsertPrescription)

}

func HandleAddPrescriptionToProtocolByCategory(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c,w,r, addPxToProtocolMedCategory)

}

func HandleRemovePrescriptionFromProtocolByCategory(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c,w,r,removePxFromProtocolMedCategory)
}


func HandleGetPrescriptionsByProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGet(c,w,r,getPxByProtocolID)

}

func HandleGetPrescriptionsByCategory(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGet(c,w,r, getProtocolMedCategoryByID)

}
func HandleUpsertProtocolMedCategory(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleUpsert(c, w, r, upsertProtocolMedCategory)
}

func HandleDeleteProtocolMedCategory(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c, w, r, deleteProtocolMedCategoryByID)
}

