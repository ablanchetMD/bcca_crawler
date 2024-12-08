package routes

import (
	"bcca_crawler/internal/config"
	"net/http"
	"bcca_crawler/api"
)

func RegisterUserRoutes(prefix string, mux *http.ServeMux, s *config.Config) {
	mux.HandleFunc(prefix +"/users", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:			
			v := QueryValidation{
				ValidSortBy:  []string{"email"},
				MaxLimit: 100,
				MinLimit: 1,								
			}
			params, err := ParseQueryParams(r,v)

			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			
			api.HandleGetUsers(s, *params, w, r)
		case http.MethodPost:
			api.HandleCreateUser(s, w, r)				
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/users/refresh", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			api.HandleRefresh(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/users/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			api.HandleLogin(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/users/revoke", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			api.HandleRevoke(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/users/reset", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			api.HandleReset(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.HandleGetUserById(s, w, r)
		case http.MethodPut:
			api.HandleUpdateUser(s, w, r)
		case http.MethodDelete:
			api.HandleDeleteUserById(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}