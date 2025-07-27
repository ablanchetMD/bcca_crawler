package protocols

import (
	"bcca_crawler/api"
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"	
	"context"	
	"fmt"
	"net/http"
	"net/url"
	
)


type LabReq struct {
	ID           string  `json:"id" validate:"omitempty,uuid"`
	Name         string  `json:"name" validate:"required,min=1,max=250"`
	Description  string  `json:"description" validate:"required,min=1,max=500"`
	FormUrl      string  `json:"form_url" validate:"omitempty,url"`
	Unit         string  `json:"unit" validate:"omitempty,min=1,max=50"`
	LowerLimit   float64 `json:"lower_limit" validate:"omitempty,min=1,max=50"`
	UpperLimit   float64 `json:"upper_limit" validate:"omitempty,min=1,max=50"`
	TestCategory string  `json:"test_category" validate:"omitempty,min=1,max=50"`
}

type LabCategoryReq struct {
	ID            string `json:"id" validate:"omitempty,uuid"`	
	Category      string `json:"category" validate:"required,min=1,max=500"`
	Comments	  string `json:"comments" validate:"omitempty"`
	Position      int32  `json:"position" validate:"gte=0"`
}

func getLabs(c *config.Config, ctx context.Context, ids api.IDs, query url.Values) ([]LabResp, error) {
	category := query.Get("test_category")
	var items []database.Test
	switch category {
	case "":
		i, err := c.Db.GetTests(ctx)
		if err != nil {
			return nil, fmt.Errorf("error getting labs: %s", err)
		}
		items = i
	default:
		i, err := c.Db.GetTestsByCategory(ctx, category)
		if err != nil {
			return nil, fmt.Errorf("error getting labs: %s", err)
		}
		items = i
	}
	
	return api.MapAll(items, MapLab), nil

}
func getLabByID(c *config.Config, ctx context.Context, ids api.IDs) (LabResp, error) {
	item, err := c.Db.GetTestByID(ctx, ids.ID)

	if err != nil {		
		return LabResp{}, fmt.Errorf("error getting test: %s", ids.ID.String())
	}	

	return MapLab(item), nil
}

func deleteLabByID(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.DeleteTest(ctx, ids.ID)
	if err != nil {		
		return "", fmt.Errorf("error deleting test: %s", ids.ID.String())
	}	

	return fmt.Sprintf("Test %s deleted.", ids.ID.String()), nil
}

func upsertLab(c *config.Config, ctx context.Context,req LabReq, ids api.IDs)  error {
	id := api.ParseOrNilUUID(req.ID)
	_, err := c.Db.UpsertTest(ctx, database.UpsertTestParams{
		ID:           id,
		Name:         req.Name,
		Description:  req.Description,
		FormUrl:      req.FormUrl,
		Unit:         req.Unit,
		LowerLimit:   req.LowerLimit,
		UpperLimit:   req.UpperLimit,
		TestCategory: req.TestCategory,
	})

	if err != nil {
		return fmt.Errorf("error upserting Test: %s with error:%s", req.ID, err.Error())
	}

	return nil

}

func upsertLabCategory(c *config.Config, ctx context.Context,req LabCategoryReq, ids api.IDs)  error {
	id := api.ParseOrNilUUID(req.ID)

	_, err := c.Db.UpsertProtoTestCategory(ctx, database.UpsertProtoTestCategoryParams{
		ID: id,
		ProtocolID: ids.ProtocolID,
		Category: req.Category,
		Comments: req.Comments,
		Position: req.Position,
	})

	if err != nil {
		return err
	}

	return nil

}

func getLabCategoryByProtocolID(c *config.Config, ctx context.Context, ids api.IDs) ([]ProtocolTestGroup, error) {
	items, err := c.Db.GetProtocolTests(ctx, ids.ProtocolID)
	if err != nil {
		return nil, err
	}

	response, err := api.ToResponseData[[]ProtocolTestGroup](items)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func getLabCategoryByID(c *config.Config, ctx context.Context, ids api.IDs) (ProtocolTestGroup, error) {
	item, err := c.Db.GetTestCategoryByID(ctx, ids.ID)
	if err != nil {
		return ProtocolTestGroup{}, err
	}

	response, err := api.ToResponseData[ProtocolTestGroup](item)
	if err != nil {
		return ProtocolTestGroup{}, err
	}	

	return response, nil
}

func deleteLabCategoryByID(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.RemoveTestCategoryByID(ctx, ids.ID)

	if err != nil {
		return "", fmt.Errorf("error deleting test category: %s, with error: %v", ids.ID.String(), err)

	}
	return fmt.Sprintf("Test Category %s deleted.", ids.ID.String()), nil
}

func addLabToLabCategory(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.AddTestToProtoTestCategory(ctx, database.AddTestToProtoTestCategoryParams{
		TestsID: ids.ID,
		ProtocolTestsID:    ids.LabCategoryID,
	})
	if err != nil {
		return "", fmt.Errorf("error adding test: %s to test category %s, with error: %v", ids.ID.String(), ids.LabCategoryID.String(), err)

	}
	return fmt.Sprintf("Test %s added to test category %s", ids.ID.String(), ids.LabCategoryID.String()), nil
}

func removeLabFromLabCategory(c *config.Config, ctx context.Context, ids api.IDs) (string, error) {
	err := c.Db.RemoveTestToProtoTestCategory(ctx, database.RemoveTestToProtoTestCategoryParams{
		TestsID: ids.ID,
		ProtocolTestsID:    ids.LabCategoryID,
	})
	if err != nil {
		return "", fmt.Errorf("error removing test: %s to test category %s, with error: %v", ids.ID.String(), ids.LabCategoryID.String(), err)

	}
	return fmt.Sprintf("Test %s removed from test category %s", ids.ID.String(), ids.LabCategoryID.String()), nil
}


func HandleGetLabs(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGetWithQ(c, w, r, getLabs)
}

func HandleGetLabByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGet(c, w, r, getLabByID)	
}

func HandleDeleteLabByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c,w,r, deleteLabByID)	
}

func HandleUpsertLab(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleUpsert(c,w,r, upsertLab)	
}

func HandleAddTestGroupToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleUpsert(c,w,r,upsertLabCategory)	
}

func HandleGetTestGroupByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGet(c, w, r, getLabCategoryByID)
}

func HandleDeleteTestGroupByID(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c, w, r, deleteLabCategoryByID)
}

func HandleGetLabsByProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleGet(c, w, r,getLabCategoryByProtocolID)
}

func HandleAddTestToTestGroup(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c, w, r, addLabToLabCategory)
}
func HandleRemoveTestFromTestGroup(c *config.Config, w http.ResponseWriter, r *http.Request) {
	api.HandleModify(c, w, r, removeLabFromLabCategory)
}

