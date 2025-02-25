package routes

import (
	"bcca_crawler/internal/config"
	"net/http"
	"bcca_crawler/api"
)

func RegisterReferencesRoute(prefix string, mux *http.ServeMux, s *config.Config) {
	mux.HandleFunc(prefix +"/references", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.HandleGetArticleReferences(s, w, r)
		case http.MethodPut:
			api.HandleUpsertReference(s, w, r)				
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/references/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.HandleGetArticleRefByID(s, w, r)
		case http.MethodDelete:
			api.HandleDeleteArticleRefByID(s, w, r)						
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	
}