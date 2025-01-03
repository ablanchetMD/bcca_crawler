package routes

import (
	"bcca_crawler/internal/config"
	"net/http"
	"bcca_crawler/api"
)

func RegisterProtocolRoutes(prefix string, mux *http.ServeMux, s *config.Config) {
	mux.HandleFunc(prefix +"/protocols", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			//optional queries : sort, sort_by, page, limit, offset, filter, fields, include, exclude,
			v := QueryValidation{
				ValidSortBy:  []string{"name"},
				MaxLimit: 100,
				MinLimit: 1,								
			}
			params, err := ParseQueryParams(r,v)

			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			
			api.HandleGetProtocols(s, *params, w, r)
		case http.MethodPost:
			api.HandleCreateProtocol(s, w, r)				
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/protocols/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.HandleGetProtocolById(s, w, r)
		case http.MethodPut:
			api.HandleUpdateProtocol(s, w, r)
		case http.MethodDelete:
			api.HandleDeleteProtocol(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/protocols/summary/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.HandleGetProtocolSummary(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/protocols/summarycode/{code}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.HandleGetProtocolSummaryCode(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

}