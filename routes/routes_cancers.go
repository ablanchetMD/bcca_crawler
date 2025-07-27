package routes

import (
	"bcca_crawler/internal/config"

	"bcca_crawler/api"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterCancerRoutes(prefix string, mux *mux.Router, s *config.Config) {
	mux.HandleFunc(prefix+"/cancers", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			//optional queries : sort, sort_by, page, limit, offset, filter, fields, include, exclude,
			v := QueryValidation{
				ValidSortBy: []string{"name"},
				MaxLimit:    100,
				MinLimit:    1,
			}
			params, err := ParseQueryParams(r, v)

			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			api.HandleGetCancers(s, *params, w, r)
		case http.MethodPost:
			api.HandleCreateCancer(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix+"/cancers/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			api.HandleGetCancerById(s, w, r)
		case http.MethodPut:
			api.HandleUpdateCancer(s, w, r)
		case http.MethodDelete:
			api.HandleDeleteCancer(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}
