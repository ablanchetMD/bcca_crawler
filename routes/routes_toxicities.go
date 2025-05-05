package routes

import (	
	"bcca_crawler/api/protocols"
	"bcca_crawler/internal/config"
	"github.com/gorilla/mux"
	"net/http"	
)

func RegisterToxicitiesRoutes(prefix string, mux *mux.Router, s *config.Config) {
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

	

}