package routes

import (
	"bcca_crawler/api" 
	"bcca_crawler/internal/config"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"strconv"
)

type QueryValidation struct {
	ValidSortBy  []string	
	MaxLimit	 int
	MinLimit	 int
}


func ParseQueryParams(r *http.Request, v QueryValidation) (*api.QueryParams, error) {
	q := r.URL.Query()
	params := &api.QueryParams{
		Sort:   strings.ToLower(q.Get("sort")),
		SortBy: strings.ToLower(q.Get("sort_by")),
		FilterBy: strings.ToLower(q.Get("filter_by")),
	}

	// Default validation and assignment
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	// Validate SortBy
	if !contains(v.ValidSortBy, params.SortBy) {
		params.SortBy = v.ValidSortBy[0] // Default to the first valid field
	}

	// Parse optional integers
	if page := q.Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			params.Page = p
		}
	}

	if limit := q.Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			params.Limit = l
		}
	}	

	if offset := q.Get("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			params.Offset = o
		}
	}

	if params.Limit < v.MinLimit {
		params.Limit = 10
	}

	if params.Limit > v.MaxLimit {
		params.Limit = 100
	}

	if params.Page <1 {
		params.Page = 1
	}

	if params.Offset < 1 {				
		params.Offset = (params.Page - 1) * params.Limit				
	}

	// Parse comma-separated fields for include/exclude
	if fields := q.Get("fields"); fields != "" {
		params.Fields = strings.Split(strings.ToLower(fields), ",")
	}
	if include := q.Get("include"); include != "" {
		params.Include = strings.Split(strings.ToLower(include), ",")
	}
	if exclude := q.Get("exclude"); exclude != "" {
		params.Exclude = strings.Split(strings.ToLower(exclude), ",")
	}

	return params, nil
}

func RegisterRoutes(router *mux.Router, s *config.Config) {
	pre := "/api/v1"
	// Register all routes
	RegisterProtocolRoutes(pre, router, s)
	RegisterUserRoutes(pre, router, s)
	RegisterCancerRoutes(pre, router, s)
	RegisterCriteriaRoutes(pre, router, s)
	RegisterLabRoutes(pre, router, s)
	RegisterMedRoutes(pre, router, s)
	RegisterReferencesRoute(pre, router, s)
	RegisterPhysicianRoutes(pre, router, s)
	RegisterToxicitiesRoutes(pre, router, s)
	RegisterTreatmentRoutes(pre, router, s)

}

// Helper function to check if a value exists in a list
func contains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}
