package routes

import (
	"bcca_crawler/internal/config"
	"github.com/gorilla/mux"
	"net/http"	
	"bcca_crawler/api/protocols"
)

func RegisterTreatmentRoutes(prefix string, mux *mux.Router, s *config.Config) {
	// Define UUID pattern for consistent use
    uuidPattern := "[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}"
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

	mux.HandleFunc(prefix +"/treatments/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetTreatmentByID(s, w, r)			
		case http.MethodDelete:
			protocols.HandleDeleteTreatmentByID(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix + "/cycles/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {		
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetTreatmentsByCycleID(s, w, r)		
		case http.MethodDelete:
			protocols.HandleDeleteCycleByID(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Create a subrouter for protocol_id routes
    cycleRouter := mux.PathPrefix(prefix+"/cycles/{cycle_id:"+uuidPattern+"}").Subrouter()	

	cycleRouter.HandleFunc("/treatments/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {
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

	
}