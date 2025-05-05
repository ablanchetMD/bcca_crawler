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

type LabReq struct {
	ID 			string `json:"id" validate:"omitempty,uuid"`
	Name 		string `json:"name" validate:"required,min=1,max=250"`
	Description string `json:"description" validate:"required,min=1,max=500"`
	FormUrl 	string `json:"form_url" validate:"omitempty,url"`
	Unit 		string `json:"unit" validate:"omitempty,min=1,max=50"`
	LowerLimit 	float64 `json:"lower_limit" validate:"omitempty,min=1,max=50"`
	UpperLimit 	float64 `json:"upper_limit" validate:"omitempty,min=1,max=50"`
	TestCategory string `json:"test_category" validate:"omitempty,min=1,max=50"`	
}


func HandleGetLabs(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	labs := []api.LabResp{}

	category := r.URL.Query().Get("test_category")

	var test_labs []database.Test

	switch category {
		case "":
			raw_labs, err := c.Db.GetTests(ctx)
			if err != nil {
				json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting labs")
				return
			}
			test_labs = raw_labs
		default:
			raw_labs, err := c.Db.GetTestsByCategory(ctx, category)
			if err != nil {
				json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting labs")
				return
			}
			test_labs = raw_labs
		}	
	
	for _, a := range test_labs {		
		labs = append(labs, api.MapLab(a))	
	}

	json_utils.RespondWithJSON(w, http.StatusOK, labs)
}
	

func HandleGetLabByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw_test, err := c.Db.GetTestByID(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting test: %s", ids.ID.String()))
		return
	}		

	Test := api.MapLab(raw_test)
	
	json_utils.RespondWithJSON(w, http.StatusOK, Test)
}

func HandleDeleteLabByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteTest(ctx, ids.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting test: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "test deleted"})
}

func HandleUpsertLab(c *config.Config, w http.ResponseWriter, r *http.Request) {

	var req LabReq	
	err := api.UnmarshalAndValidatePayload(c,r, &req)
	if err != nil {
		println(err.Error())
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())		
		return
	}

	ctx := r.Context()	

	pid, err:= uuid.Parse(req.ID)
	if err != nil {
		pid = uuid.Nil
	}		
	
	test,err := c.Db.UpsertTest(ctx,database.UpsertTestParams{
		Column1: pid,
		Column2: req.Name,
		Column3: req.Description,
		Column4: req.FormUrl,
		Column5: req.Unit,
		Column6: req.LowerLimit,
		Column7: req.UpperLimit,
		Column8: req.TestCategory,
	})	

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error upserting test: %s", req.ID))
		return		
	}
	
	return_test := api.MapLab(test)

	json_utils.RespondWithJSON(w, http.StatusOK, return_test)	
}


func HandleAddLabToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	category := r.URL.Query().Get("test_category")

	err = c.Validate.Var(category, "required, test_protocol_category")

	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid test_category: %s", category))
		return
	}

	urgency := r.URL.Query().Get("test_urgency")

	err = c.Validate.Var(urgency, "required, test_protocol_urgency")

	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid test_urgency: %s", urgency))
		return
	}
	
	_,err = c.Db.AddTestToProtocolByCategoryAndUrgency(ctx, database.AddTestToProtocolByCategoryAndUrgencyParams{
		ProtocolID: ids.ProtocolID,
		TestID: ids.ID,
		Category: database.CategoryEnum(category),
		Urgency: database.UrgencyEnum(urgency),
	})	

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding test to protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "test added to protocol"})

}

func HandleGetLabsByProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	labs := []api.LabResp{}
	
	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}	
	
	category := r.URL.Query().Get("test_category")

	err = c.Validate.Var(category, "required, test_protocol_category")

	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid test_category: %s", category))
		return
	}

	urgency := r.URL.Query().Get("test_urgency")

	err = c.Validate.Var(urgency, "required, test_protocol_urgency")

	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid test_urgency: %s", urgency))
		return
	}

	
	raw_labs, err := c.Db.GetTestsByProtocolByCategoryAndUrgency(ctx, database.GetTestsByProtocolByCategoryAndUrgencyParams{
		ProtocolID: ids.ProtocolID,
		Category: database.CategoryEnum(category),
		Urgency: database.UrgencyEnum(urgency),
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting labs")
		return
	}
	
	for _, a := range raw_labs {		
		labs = append(labs, api.MapLab(a))		
	}

	json_utils.RespondWithJSON(w, http.StatusOK, labs)
}

func HandleRemoveLabFromProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := api.ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	category := r.URL.Query().Get("test_category")

	err = c.Validate.Var(category, "required, test_protocol_category")

	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid test_category: %s", category))
		return
	}

	urgency := r.URL.Query().Get("test_urgency")

	err = c.Validate.Var(urgency, "required, test_protocol_urgency")

	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid test_urgency: %s", urgency))
		return
	}
	
	
	err = c.Db.RemoveTestFromProtocolByCategoryAndUrgency(ctx, database.RemoveTestFromProtocolByCategoryAndUrgencyParams{
		ProtocolID: ids.ProtocolID,
		TestID: ids.ID,
		Category: database.CategoryEnum(category),
		Urgency: database.UrgencyEnum(urgency),
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing lab test from protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "lab test removed from protocol"})

}