package routes

import (
	"bcca_crawler/internal/config"
	"net/http"	
	"bcca_crawler/api/protocols"
)

func RegisterTreatmentRoutes(prefix string, mux *http.ServeMux, s *config.Config) {
	mux.HandleFunc(prefix +"/treatments", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:			
			protocols.HandleGetTreatments(s, w, r)
		case http.MethodPut:
			protocols.HandleUpsertTreatment(s, w, r)			
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/treatments/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetTreatmentByID(s, w, r)			
		case http.MethodDelete:
			protocols.HandleDeleteTreatmentByID(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix + "/cycles/{cycle_id}/treatments", func(w http.ResponseWriter, r *http.Request) {
		//query = cycle_id
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetTreatmentsByCycleID(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix + "/cycles/{cycle_id}/treatments/{id}", func(w http.ResponseWriter, r *http.Request) {
		//query = cycle_id
		switch r.Method {
		case http.MethodPost:
			protocols.HandleAddTreatmentToCycle(s, w, r)
		case http.MethodDelete:
			protocols.HandleRemoveTreatmentToCycle(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})	

	mux.HandleFunc(prefix + "/cycles/{id}", func(w http.ResponseWriter, r *http.Request) {		
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetCycleByID(s, w, r)
		case http.MethodDelete:
			protocols.HandleDeleteCycleByID(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

}