package routes

import (
	"bcca_crawler/api"
	"bcca_crawler/api/protocols"
	"bcca_crawler/internal/config"
	"net/http"	
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

	mux.HandleFunc(prefix +"/protocols/{protocol_id}/eligibility_criteria/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			protocols.HandleRemoveEligibilityFromProtocol(s, w, r)
		case http.MethodPost:
			protocols.HandleAddEligibilityCriteriaToProtocol(s, w, r)			
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})		

	mux.HandleFunc(prefix +"/protocols/{protocol_id}/cautions/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			protocols.HandleRemoveCautionFromProtocol(s, w, r)
		case http.MethodPost:
			protocols.HandleAddCautionToProtocol(s, w, r)			
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})	
	
	mux.HandleFunc(prefix +"/protocols/{protocol_id}/precautions/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			protocols.HandleRemovePrecautionFromProtocol(s, w, r)
		case http.MethodPost:
			protocols.HandleAddPrecautionToProtocol(s, w, r)			
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/protocols/{protocol_id}/labs", func(w http.ResponseWriter, r *http.Request) {
		//queries : test_category, test_urgency
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetLabsByProtocol(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	
	mux.HandleFunc(prefix +"/protocols/{protocol_id}/labs/{id}", func(w http.ResponseWriter, r *http.Request) {
		//queries : test_category, test_urgency, protocol_id
		switch r.Method {
		case http.MethodDelete:
			protocols.HandleRemoveLabFromProtocol(s, w, r)
		case http.MethodPost:
			protocols.HandleAddLabToProtocol(s, w, r)			
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/protocols/toxicity_adjustments/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			protocols.HandleRemoveAdjustmentsToProtocol(s, w, r)							
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/protocols/{protocol_id}/toxicity_adjustments", func(w http.ResponseWriter, r *http.Request) {
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

	mux.HandleFunc(prefix +"/protocols/{protocol_id}/prescriptions", func(w http.ResponseWriter, r *http.Request) {
		//queries : prescription_category, protocol_id
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetPrescriptionsByCategory(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	
	mux.HandleFunc(prefix +"/protocols/{protocol_id}/prescriptions/{id}", func(w http.ResponseWriter, r *http.Request) {
		//queries : prescription_category, protocol_id
		switch r.Method {
		case http.MethodDelete:
			protocols.HandleRemovePrescriptionFromProtocolByCategory(s, w, r)
		case http.MethodPost:
			protocols.HandleAddPrescriptionToProtocolByCategory(s, w, r)			
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix + "protocols/{id}/cycles", func(w http.ResponseWriter, r *http.Request) {		
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetCyclesByProtocolID(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/protocols/{id}/medication_modifications", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetMedModificationsByProtocol(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix + "/protocols/{id}/cycles", func(w http.ResponseWriter, r *http.Request) {
		//query = protocol_id
		switch r.Method {
		case http.MethodGet:
			protocols.HandleGetCycles(s, w, r)
		case http.MethodPost:
			protocols.HandleUpsertTreatmentCycle(s, w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc(prefix +"/protocols/{id}/summary", func(w http.ResponseWriter, r *http.Request) {
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