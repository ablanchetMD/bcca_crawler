package routes

import (	
	"bcca_crawler/api/protocols"
	"bcca_crawler/internal/config"
	"net/http"	
)

func RegisterToxicitiesRoutes(prefix string, mux *http.ServeMux, s *config.Config) {
	mux.HandleFunc(prefix +"/toxicities", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetToxicities(s, w, r)
		case http.MethodPut:
			protocols.HandlerUpsertToxicityWithGrades(s, w, r)				
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/toxicities/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetToxicityByID(s, w, r)
		case http.MethodDelete:
			protocols.HandleDeleteToxicityByID(s, w, r)						
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/protocols/toxicities/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			protocols.HandleRemoveAdjustmentsToProtocol(s, w, r)							
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/protocols/toxicities", func(w http.ResponseWriter, r *http.Request) {
		//protocol_id
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetToxicitiesWithAdjustmentsByProtocolID(s, w, r)
		case http.MethodPut:
			protocols.HandleUpsertAdjustmentsToProtocol(s, w, r)								
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})	

}