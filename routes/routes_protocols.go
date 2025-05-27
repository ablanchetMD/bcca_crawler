package routes

import (
	"bcca_crawler/api"
	"bcca_crawler/api/protocols"
	"bcca_crawler/internal/config"
	"github.com/gorilla/mux"
	"net/http"	
)

func RegisterProtocolRoutes(prefix string, router *mux.Router, s *config.Config) {
    // Define UUID pattern for consistent use
    uuidPattern := "[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}"
    
    // Base protocols endpoints
    router.HandleFunc(prefix+"/protocols", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
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
            api.HandleGetProtocols(s, *params, w, r)
        case http.MethodPost:
            api.HandleCreateProtocol(s, w, r)
        case http.MethodPut:
            api.HandleUpsertProtocol(s, w, r)
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("GET", "POST", "PUT")

    // Protocol by ID
    router.HandleFunc(prefix+"/protocols/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            api.HandleGetProtocolById(s, w, r)        
        case http.MethodDelete:
            api.HandleDeleteProtocol(s, w, r)
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("GET", "PUT", "DELETE")
    
    // Special route for summary by code (doesn't follow the UUID pattern)
    router.HandleFunc(prefix+"/protocols/summarycode/{code}", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            api.HandleGetProtocolSummaryCode(s, w, r)
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("GET")
    
    

    // Create a subrouter for protocol_id routes
    protocolRouter := router.PathPrefix(prefix+"/protocols/{protocol_id:"+uuidPattern+"}").Subrouter()

    // Protocol summary
    protocolRouter.HandleFunc("/summary", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            api.HandleGetProtocolSummary(s, w, r)
        } else {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("GET")
    
    // Eligibility criteria routes
    protocolRouter.HandleFunc("/eligibility_criteria",func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            protocols.HandleGetEligibilityCriteriaByProtocol(s,w,r)        
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }}).Methods("GET")
    
    protocolRouter.HandleFunc("/eligibility_criteria/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodDelete:
            protocols.HandleRemoveEligibilityFromProtocol(s, w, r)
        case http.MethodPost:
            protocols.HandleAddEligibilityCriteriaToProtocol(s, w, r)
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("DELETE", "POST")
    
    // Physicians routes
    protocolRouter.HandleFunc("/physicians",func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            protocols.HandleGetEligibilityCriteriaByProtocol(s,w,r)        
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }}).Methods("GET")

    protocolRouter.HandleFunc("/physicians/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodDelete:
            api.HandleRemovePhysicianFromProtocol(s, w, r)
        case http.MethodPost:
            api.HandleAddPhysicianToProtocol(s, w, r)
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("DELETE", "POST")
    
    // Cautions routes
    protocolRouter.HandleFunc("/cautions",func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            protocols.HandleGetEligibilityCriteriaByProtocol(s,w,r)        
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }}).Methods("GET")

    protocolRouter.HandleFunc("/cautions/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodDelete:
            protocols.HandleRemoveCautionFromProtocol(s, w, r)
        case http.MethodPost:
            protocols.HandleAddCautionToProtocol(s, w, r)
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("DELETE", "POST")
    
    // References routes
    protocolRouter.HandleFunc("/references",func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            protocols.HandleGetEligibilityCriteriaByProtocol(s,w,r)        
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }}).Methods("GET")

    protocolRouter.HandleFunc("/references/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodDelete:
            api.HandleRemoveArticleFromProtocol(s, w, r)
        case http.MethodPost:
            api.HandleAddArticleToProtocol(s, w, r)
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("DELETE", "POST")
    
    // Precautions routes
    protocolRouter.HandleFunc("/precautions",func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            protocols.HandleGetEligibilityCriteriaByProtocol(s,w,r)        
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }}).Methods("GET")

    protocolRouter.HandleFunc("/precautions/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodDelete:
            protocols.HandleRemovePrecautionFromProtocol(s, w, r)
        case http.MethodPost:
            protocols.HandleAddPrecautionToProtocol(s, w, r)
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("DELETE", "POST")
    
    // Labs routes
    protocolRouter.HandleFunc("/labs", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            protocols.HandleGetLabsByProtocol(s, w, r)
        } else {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("GET")
    
    protocolRouter.HandleFunc("/labs/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodDelete:
            protocols.HandleRemoveLabFromProtocol(s, w, r)
        case http.MethodPost:
            protocols.HandleAddLabToProtocol(s, w, r)
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("DELETE", "POST")
    
    // Toxicity adjustments routes
    protocolRouter.HandleFunc("/toxicity_adjustments", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            protocols.HandleGetToxicitiesWithAdjustmentsByProtocolID(s, w, r)
        case http.MethodPut:
            protocols.HandleUpsertAdjustmentsToProtocol(s, w, r)
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("GET", "PUT")
    
    protocolRouter.HandleFunc("/toxicity_adjustments/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodDelete {
            protocols.HandleRemoveAdjustmentsToProtocol(s, w, r)
        } else {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("DELETE")
    
    // Prescriptions routes
    protocolRouter.HandleFunc("/prescriptions", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            protocols.HandleGetPrescriptionsByCategory(s, w, r)
        } else {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("GET")
    
    protocolRouter.HandleFunc("/prescriptions/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodDelete:
            protocols.HandleRemovePrescriptionFromProtocolByCategory(s, w, r)
        case http.MethodPost:
            protocols.HandleAddPrescriptionToProtocolByCategory(s, w, r)
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("DELETE", "POST")
    
    // Cycles routes
    protocolRouter.HandleFunc("/cycles", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            protocols.HandleGetCyclesByProtocolID(s, w, r)
        case http.MethodPost:
            protocols.HandleUpsertTreatmentCycle(s, w, r)
        default:
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("GET", "POST")
    
    // Medication modifications
    protocolRouter.HandleFunc("/medication_modifications", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            protocols.HandleGetMedModificationsByProtocol(s, w, r)
        } else {
            http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        }
    }).Methods("GET")
}
