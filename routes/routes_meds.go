package routes

import (	
	"bcca_crawler/api/protocols"
	"bcca_crawler/internal/config"
	"net/http"	
)

func RegisterMedRoutes(prefix string, mux *http.ServeMux, s *config.Config) {
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

	mux.HandleFunc(prefix +"/medications/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetMedByID(s, w, r)
		case http.MethodDelete:
			protocols.HandleDeleteMedByID(s, w, r)						
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/medications/{id}/modifications", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetMedModificationsByMedication(s, w, r)							
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/medications/modifications/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetMedModByID(s, w, r)
		case http.MethodDelete:
			protocols.HandleDeleteMedModByID(s, w, r)							
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/medications/modifications", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {		
		case http.MethodPut:
			protocols.HandlerUpsertMedMod(s, w, r)							
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
	
	mux.HandleFunc(prefix +"/prescriptions/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetPrescriptionByID(s, w, r)
		case http.MethodDelete:
			protocols.HandleDeletePrescriptionByID(s, w, r)						
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

}