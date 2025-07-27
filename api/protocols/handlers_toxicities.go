package protocols

import (
	"bcca_crawler/api"
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"
	"bcca_crawler/internal/json_utils"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func HandleGetToxicities(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	items, err := c.Db.GetToxicitiesWithGrades(ctx)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting toxicities:%s", err.Error()))
		return
	}

	returned_value, err := api.MapAllWithError(items, api.MapToxicityWithGrades)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error mapping toxicities")
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, returned_value)
}

func HandleGetToxicityByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw, err := c.Db.GetToxicityByID(ctx, ids.ID)

	if err != nil {

		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting toxicity by id with error: %s", err.Error()))
		return
	}

	fobj, err := api.MapToxicityWithGrades(raw)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error mapping toxicity by id")
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, fobj)
}

func HandleGetToxicitiesWithAdjustmentsByProtocolID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	items, err := c.Db.GetToxicitiesWithGradesAndAdjustmentsByProtocol(ctx, ids.ProtocolID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting toxicities for the protocol: %s", ids.ProtocolID.String()))
		return
	}

	returned_values, err := api.MapAllWithError(items, api.MapToxicityWithGradesWithAdjust)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error mapping toxicities")
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, returned_values)
}

// To do

func HandleDeleteToxicityByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.RemoveToxicity(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting toxicity: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "toxicity deleted"})
}

func HandlerUpsertToxicityWithGrades(c *config.Config, w http.ResponseWriter, r *http.Request) {

	var req ToxicityReq
	err := api.UnmarshalAndValidatePayload(c, r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()

	params := database.UpsertToxicityWithGradesParams{
		ID:               api.ParseOrGenerateUUID(req.ID),
		Title:            req.Title,
		Category:         req.Category,
		Description:      req.Description,
		GradeIds:         make([]uuid.UUID, len(req.Grades)),
		GradeNumber:      make([]database.GradeEnum, len(req.Grades)),
		GradeDescription: make([]string, len(req.Grades)),
	}

	for i, g := range req.Grades {
		params.GradeIds[i] = api.ParseOrGenerateUUID(g.ID)
		params.GradeNumber[i] = database.GradeEnum(g.Grade)
		params.GradeDescription[i] = g.Description
	}

	err = c.Db.UpsertToxicityWithGrades(ctx, params)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upserting toxicity: %s", req.ID))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, req)
}

func HandleUpsertAdjustmentsToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req ToxModificationsReq

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = api.UnmarshalAndValidatePayload(c, r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	params := database.UpsertToxicityModificationParams{
		ProtocolID: ids.ProtocolID,
		ID:         make([]uuid.UUID, len(req.Grades)),
		GradeIds:   make([]uuid.UUID, len(req.Grades)),
		Adjustment: make([]string, len(req.Grades)),
	}

	for i, g := range req.Grades {
		params.ID[i] = api.ParseOrGenerateUUID(g.ID)
		params.GradeIds[i] = api.ParseOrGenerateUUID(g.GradeID)
		params.Adjustment[i] = g.Adjustment
	}

	err = c.Db.UpsertToxicityModification(ctx, params)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding Toxicity Adjustment to Protocol: %s", ids.ProtocolID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "toxicity adjustments added to protocol"})

}

func HandleRemoveAdjustmentsToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteProtocolToxModificationsByProtocolAndToxicity(ctx, database.DeleteProtocolToxModificationsByProtocolAndToxicityParams{
		ProtocolID: id.ProtocolID,
		ToxicityID: id.ID,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing toxicity adjustments from protocol: %s", id.ProtocolID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "toxicity adjustments removed from protocol"})

}
