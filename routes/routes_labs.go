package routes

import (
	"bcca_crawler/internal/config"
	"net/http"	
	"bcca_crawler/api/protocols"
	"github.com/gorilla/mux"
)

func RegisterLabRoutes(prefix string, mux *mux.Router, s *config.Config) {
	uuidPattern := "[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}"
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

	mux.HandleFunc(prefix + "/labgroup/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {		
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetTestGroupByID(s, w, r)		
		case http.MethodDelete:
			protocols.HandleDeleteTestGroupByID(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	// Create a subrouter for labgroup routes
    labgroupRouter := mux.PathPrefix(prefix+"/labgroup/{lab_category_id:"+uuidPattern+"}").Subrouter()	

	labgroupRouter.HandleFunc("/labs/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {
		
		switch r.Method {
		case http.MethodPost:
			protocols.HandleAddTestToTestGroup(s, w, r)
		case http.MethodDelete:
			protocols.HandleRemoveTestFromTestGroup(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})	





}