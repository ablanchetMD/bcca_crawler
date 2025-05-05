package routes

import (
	"bcca_crawler/internal/config"
	"net/http"
	"bcca_crawler/api"
	"github.com/gorilla/mux"
)

func RegisterPhysicianRoutes(prefix string, mux *mux.Router, s *config.Config) {
	mux.HandleFunc(prefix +"/physicians", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.HandleGetPhysicians(s, w, r)
		case http.MethodPut:
			api.HandleUpsertPhysician(s, w, r)				
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/physicians/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.HandleGetPhysicianByID(s, w, r)
		case http.MethodDelete:
			api.HandleDeletePhysicianByID(s, w, r)						
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	
}