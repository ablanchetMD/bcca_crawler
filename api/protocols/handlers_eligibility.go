package protocols

import (
	"bcca_crawler/api"
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"
	"bcca_crawler/internal/json_utils"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

func HandleGetEligibilityCriteria(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	category := r.URL.Query().Get("category")

	switch category {
	case "":
		items, err := c.Db.GetElibilityCriteria(ctx)

		if err != nil {
			println(err.Error())
			json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting eligibility criteria with no requirements")
			return
		}
		returned_value, err := api.MapAllWithError(items, MapEligibilityWithProtocols)

		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, "Error mapping eligibility criteria")
			return
		}
		json_utils.RespondWithJSON(w, http.StatusOK, returned_value)
	default:

		err := c.Validate.Var(category, "required,eligibility_criteria")
		if err != nil {
			json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		items, err := c.Db.GetEligibilityCriteriaByType(ctx, category)

		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting eligibility criteria with requirements")
			return
		}
		returned_value, err := api.MapAllWithError(items, MapEligibilityWithProtocols)

		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, "Error mapping eligibility criteria")
			return
		}
		json_utils.RespondWithJSON(w, http.StatusOK, returned_value)
	}

}

func HandleGetEligibilityCriteriaByProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	items, err := c.Db.GetEligibilityByProtocol(ctx, ids.ProtocolID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	json_utils.RespondWithJSON(w, http.StatusOK, items)
}

func HandleGetEligibilityCriteriaByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	item, err := c.Db.GetEligibilityCriteriaByID(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting eligibility criteria: %s", ids.ID.String()))
		return
	}

	returned_value, err := MapEligibilityWithProtocols(item)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error mapping eligibility criteria")
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, returned_value)
}

func HandleDeleteEligibilityCriteriaByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteEligibilityCriteria(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting eligibility criteria: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "eligibility criteria deleted"})
}

func HandleUpsertEligibilityCriteria(c *config.Config, w http.ResponseWriter, r *http.Request) {

	var req EligibilityCriterionReq
	err := api.UnmarshalAndValidatePayload(c, r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	pid, err := uuid.Parse(req.ID)
	if err != nil || req.ID == "" {
		pid = uuid.Nil
	}

	elig, err := c.Db.UpsertEligibilityCriteria(ctx, database.UpsertEligibilityCriteriaParams{
		ID:          pid,
		Type:        req.ToTypeEnum(),
		Description: req.Description,
	})

	if err != nil {
		println(err.Error())
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upserting eligibility criteria: %s", elig.ID))
		return
	}

	if req.ProtocolID == "" {
		json_utils.RespondWithJSON(w, http.StatusOK, api.MapEligibilityCriterion(elig))
		return
	}

	proto_id, err := uuid.Parse(req.ProtocolID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding eligibility criteria to protocol (invalid UUID): %s", req.ProtocolID))
	} else {
		err = c.Db.AddEligibilityToProtocol(ctx, database.AddEligibilityToProtocolParams{
			ProtocolID: proto_id,
			CriteriaID: elig.ID,
		})

		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding eligibility criteria to protocol: %s", req.ProtocolID))
		}

	}

	json_utils.RespondWithJSON(w, http.StatusOK, api.MapEligibilityCriterion(elig))
}

func HandleAddEligibilityCriteriaToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.LinkEligibilityToProtocol(ctx, database.LinkEligibilityToProtocolParams{
		CriteriaID: ids.ID,
		ProtocolID: ids.ProtocolID,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding eligibility criteria to protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "eligibility criteria added to protocol"})

}

func HandleRemoveEligibilityFromProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.UnlinkEligibilityFromProtocol(ctx, database.UnlinkEligibilityFromProtocolParams{
		CriteriaID: ids.ID,
		ProtocolID: ids.ProtocolID,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing eligibility criteria from protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "eligibility criteria removed from protocol"})

}

func HandleUpdateEligibilityToProtocols(c *config.Config, w http.ResponseWriter, r *http.Request) {
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

	err = c.Db.UpdateEligibilityProtocols(ctx, database.UpdateEligibilityProtocolsParams{
		CriteriaID:  ids.ID,
		ProtocolIds: selectedUUIDs,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating criterias for protocol: %s", ids.ProtocolID))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "criterias updated for protocol"})
}
