package api

import (
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"
	"bcca_crawler/internal/json_utils"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type PhysicianReq struct {
	ID        string `json:"id" validate:"omitempty,uuid"`
	FirstName string `json:"first_name" validate:"required,min=1,max=500"`
	LastName  string `json:"last_name" validate:"required,min=1,max=500"`
	Email     string `json:"email" validate:"omitempty,email"`
	Site      string `json:"site" validate:"physician_site"`
}

func HandleGetPhysicians(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	raw, err := c.Db.GetPhysicians(ctx)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting physicians")
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, raw)
}

func HandleGetPhysiciansByProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	items, err := c.Db.GetPhysicianByProtocol(ctx, ids.ProtocolID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	json_utils.RespondWithJSON(w, http.StatusOK, items)
}

func HandleGetPhysicianByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw, err := c.Db.GetPhysicianByID(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting physician: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, raw)
}

func HandleDeletePhysicianByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeletePhysician(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting physician: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "physician deleted"})
}

func HandleUpsertPhysician(c *config.Config, w http.ResponseWriter, r *http.Request) {

	var req PhysicianReq
	err := UnmarshalAndValidatePayload(c, r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	pid, err := uuid.Parse(req.ID)
	if err != nil {
		pid = uuid.New()
	}

	raw, err := c.Db.UpsertPhysician(ctx, database.UpsertPhysicianParams{
		ID:        pid,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Site:      database.PhysicianSiteEnum(req.Site),
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upserting physician: %s", pid.String()))
		return
	}
	json_utils.RespondWithJSON(w, http.StatusOK, raw)
}

func HandleAddPhysicianToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.AddPhysicianToProtocol(ctx, database.AddPhysicianToProtocolParams{
		PhysicianID: ids.ID,
		ProtocolID:  ids.ProtocolID,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding physician to protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "physician added to protocol"})

}

func HandleRemovePhysicianFromProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.RemovePhysicianFromProtocol(ctx, database.RemovePhysicianFromProtocolParams{
		PhysicianID: ids.ID,
		ProtocolID:  ids.ProtocolID,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing physician from protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "physician removed from protocol"})

}
