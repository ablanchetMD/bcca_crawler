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

func getCautions(c *config.Config, ctx context.Context, ids api.IDs) ([]CautionResp, error) {
	items, err := c.Db.GetCautionWithProtocols(ctx)

	if err != nil {
		return nil, err
	}

	response, err := api.MapAllWithError(items, MapCautionWithProtocols)

	if err != nil {
		return nil, fmt.Errorf("error getting cautions: %s", err)
	}

	return response, nil
}

func getCautionsByProtocol(c *config.Config, ctx context.Context, ids api.IDs) ([]database.ProtocolCaution, error) {
	items, err := c.Db.GetProtocolCautionsByProtocol(ctx, ids.ProtocolID)

	if err != nil {
		return nil, err
	}	

	return items, nil
}

func getCautionsByID(c *config.Config, ctx context.Context, ids api.IDs) (CautionResp, error) {
	items, err := c.Db.GetCautionByIDWithProtocols(ctx, ids.ID)

	if err != nil {
		return CautionResp{}, err
	}

	response, err := MapCautionWithProtocols(items)

	if err != nil {
		return CautionResp{}, fmt.Errorf("error getting cautions: %s", err)
	}

	return response, nil
}

func deleteCaution(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.DeleteProtocolCaution(ctx, ids.ID)

	if err != nil {
		return "", fmt.Errorf("error deleting caution: %s, with error: %v", ids.ID.String(), err)

	}
	return fmt.Sprintf("Caution %s deleted.", ids.ID.String()), nil
}

func upsertCaution(c *config.Config, ctx context.Context, req CautionReq, ids api.IDs) error {

	_, err := c.Db.UpsertCaution(ctx, database.UpsertCautionParams{
		ID:          req.ID,
		Description: req.Description,
	})

	if err != nil {
		return fmt.Errorf("error upserting caution: %s with error:%s", req.ID, err.Error())
	}

	return nil
}

func addCautionToProtocol(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.AddProtocolCautionToProtocol(ctx, database.AddProtocolCautionToProtocolParams{
		CautionID:  ids.ID,
		ProtocolID: ids.ProtocolID,
	})

	if err != nil {
		return "", fmt.Errorf("error adding caution: %s to protocol %s, with error: %v", ids.ID.String(), ids.ProtocolID.String(), err)

	}
	return fmt.Sprintf("Caution %s added to protocol %s", ids.ID.String(), ids.ProtocolID.String()), nil

}

func removeCautionToProtocol(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.RemoveProtocolCautionFromProtocol(ctx, database.RemoveProtocolCautionFromProtocolParams{
		CautionID:  ids.ID,
		ProtocolID: ids.ProtocolID,
	})

	if err != nil {
		return "", fmt.Errorf("error removing caution: %s to protocol %s, with error: %v", ids.ID.String(), ids.ProtocolID.String(), err)

	}
	return fmt.Sprintf("Caution %s removed from protocol %s", ids.ID.String(), ids.ProtocolID.String()), nil
}

func updateCautionToProtocol(c *config.Config, ctx context.Context, req ChangeProtocolReq, ids api.IDs) error {
	var selectedUUIDs []uuid.UUID
	for _, id := range req.SelectedProtocolIDs {
		uid, err := uuid.Parse(id)
		if err != nil {
			return fmt.Errorf("invalid UUID:%s", id)
		}
		selectedUUIDs = append(selectedUUIDs, uid)
	}

	err := c.Db.UpdateCautionProtocols(ctx, database.UpdateCautionProtocolsParams{
		CautionID:   ids.ID,
		ProtocolIds: selectedUUIDs,
	})

	if err != nil {
		return fmt.Errorf("error updating cautions (%s) to protocols with error:%s", ids.ID.String(), err.Error())
	}

	return nil
}

func HandleGetCautions(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGet(c, w, r, getCautions)
}

func HandleGetCautionsByProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGet(c, w, r, getCautionsByProtocol)
}

func HandleGetCautionsByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGet(c, w, r, getCautionsByID)
}

func HandleDeleteCautionByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c, w, r, deleteCaution)
}

func HandleUpsertCaution(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleUpsert(c, w, r, upsertCaution)
}

func HandleAddCautionToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c, w, r, addCautionToProtocol)
}

func HandleRemoveCautionFromProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c, w, r, removeCautionToProtocol)
}

func HandleUpdateCautionsToProtocols(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleUpsert(c, w, r, updateCautionToProtocol)
}
