package routes

import (
	"bcca_crawler/api/protocols"
	"bcca_crawler/internal/config"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterMedRoutes(prefix string, mux *mux.Router, s *config.Config) {
	uuidPattern := "[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}"

	mux.HandleFunc(prefix +"/medications", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetMeds(s, w, r)
		case http.MethodPut:
			protocols.HandleUpsertMed(s, w, r)				
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/medications/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {		
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetMedByID(s, w, r)
		case http.MethodDelete:
			protocols.HandleDeleteMedByID(s, w, r)						
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/medications/{id:"+uuidPattern+"}/modifications", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetMedModificationsByMedication(s, w, r)							
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/medications/modifications", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Med Modifications Endpoint")
		switch r.Method {		
		case http.MethodPut:
			protocols.HandlerUpsertMedMod(s, w, r)							
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/medications/modifications/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {		
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetMedModByID(s, w, r)
		case http.MethodDelete:
			protocols.HandleDeleteMedModByID(s, w, r)							
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/prescriptions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetPrescriptions(s, w, r)
		case http.MethodPut:
			protocols.HandleUpsertPrescription(s, w, r)				
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	
	mux.HandleFunc(prefix +"/prescriptions/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetPrescriptionByID(s, w, r)
		case http.MethodDelete:
			protocols.HandleDeletePrescriptionByID(s, w, r)						
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix + "/pxgroup/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {		
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetPrescriptionsByCategory(s, w, r)		
		case http.MethodDelete:
			protocols.HandleDeleteProtocolMedCategory(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Create a subrouter for labgroup routes
    labgroupRouter := mux.PathPrefix(prefix+"/pxgroup/{px_category_id:"+uuidPattern+"}").Subrouter()	

	labgroupRouter.HandleFunc("/prescriptions/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {
		//query = cycle_id
		switch r.Method {
		case http.MethodPost:
			protocols.HandleAddPrescriptionToProtocolByCategory(s, w, r)
		case http.MethodDelete:
			protocols.HandleRemovePrescriptionFromProtocolByCategory(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})	

}