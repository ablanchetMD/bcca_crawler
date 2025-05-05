package routes

import (
	"bcca_crawler/internal/config"
	"net/http"	
	"bcca_crawler/api/protocols"
	"github.com/gorilla/mux"
)

func RegisterLabRoutes(prefix string, mux *mux.Router, s *config.Config) {
	mux.HandleFunc(prefix +"/labs", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:			
			protocols.HandleGetLabs(s, w, r)
		case http.MethodPut:
			protocols.HandleUpsertLab(s, w, r)			
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/labs/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetLabByID(s, w, r)			
		case http.MethodDelete:
			protocols.HandleDeleteLabByID(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}