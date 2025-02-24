package routes

import (
	"bcca_crawler/internal/config"
	"net/http"

	"bcca_crawler/api/protocols"
)

func RegisterCriteriaRoutes(prefix string, mux *http.ServeMux, s *config.Config) {
	mux.HandleFunc(prefix +"/eligibility_criteria", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetEligibilityCriteria(s, w, r)
		case http.MethodPut:
			protocols.HandleUpsertEligibilityCriteria(s, w, r)				
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/eligibility_criteria/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetEligibilityCriteriaByID(s, w, r)
		case http.MethodDelete:
			protocols.HandleDeleteEligibilityCriteriaByID(s, w, r)						
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/caution", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetCautions(s, w, r)
		case http.MethodPut:
			protocols.HandleUpsertCaution(s, w, r)				
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	
	mux.HandleFunc(prefix +"/caution/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetCautionsByID(s, w, r)
		case http.MethodDelete:
			protocols.HandleDeleteCautionByID(s, w, r)						
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/precaution", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetPrecautions(s, w, r)
		case http.MethodPut:
			protocols.HandleUpsertPrecaution(s, w, r)				
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})	

	mux.HandleFunc(prefix +"/precaution/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetPrecautionByID(s, w, r)
		case http.MethodDelete:
			protocols.HandleDeletePrecautionByID(s, w, r)						
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}