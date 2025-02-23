package api

import (
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"
	"bcca_crawler/internal/json_utils"	
	"fmt"	
	"net/http"		
	"github.com/google/uuid"
	"github.com/lib/pq"
)


type ArticleReferenceResponse struct {
	Id      uuid.UUID   `json:"Id"`
	Title   string `json:"Title"`
	Authors string `json:"Authors"`
	Journal string `json:"Journal"`
	Year    string `json:"Year"`
	Pmid    string `json:"Pmid"`
	Doi     string `json:"Doi"`
	LinkedProtocols []LinkedProtocols `json:"linked_protocols"`
}

type ArticleRefReq struct {
	ID 			string `json:"id" validate:"omitempty,uuid"`
	Title 		string `json:"title" validate:"required,min=1,max=250"`
	Authors 	string `json:"authors" validate:"required,min=1,max=250"`
	Journal 	string `json:"journal" validate:"required,min=1,max=250"`
	Year 		string `json:"year" validate:"required,min=4,max=4"`
	Pmid 		string `json:"pmid" validate:"omitempty,max=25"`
	Doi 		string `json:"doi" validate:"omitempty,max=25"`	
	ProtocolID	string `json:"protocol_id" validate:"omitempty,uuid"`
}

func HandleGetArticleReferences(c *config.Config, w http.ResponseWriter, r *http.Request){
	articles := []ArticleReferenceResponse{}
	articleReferences, err := c.Db.GetArticleReferencesWithProtocols(r.Context())
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting article references")
		return
	}
	
	for _, a := range articleReferences {
		linkedProtocols, err := ConvertTuplesToStructs[LinkedProtocols](a.ProtocolIds)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}	
		articles = append(articles, ArticleReferenceResponse{
			Id:          a.ID,
			Title:       a.Title,
			Authors:     a.Authors,
			Journal:     a.Journal,
			Year:        a.Year,
			Pmid:        a.Pmid,
			Doi:         a.Doi,
			LinkedProtocols: linkedProtocols,		
		})
	}

	json_utils.RespondWithJSON(w, http.StatusOK, articles)
	
}

func HandleUpsertReference(c *config.Config, w http.ResponseWriter, r *http.Request) {	
	
	var req ArticleRefReq	

	err := UnmarshalAndValidatePayload(c,r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	pid, err:= uuid.Parse(req.ID)
	if err != nil {
		pid = uuid.Nil
	}	
	
	article, err := c.Db.UpsertArticleReference(r.Context(), database.UpsertArticleReferenceParams{
		ID:      pid,
		Title:   req.Title,
		Authors: req.Authors,
		Journal: req.Journal,
		Year:    req.Year,
		Pmid:    req.Pmid,
		Doi:     req.Doi,		
	})

	if err != nil {		
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			// Duplicate key value violation
			json_utils.RespondWithError(w, http.StatusInternalServerError, "Record already exists")
			return			
		}
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error creating protocol")
		return
	}

	proto_id, err := uuid.Parse(req.ProtocolID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding article reference to protocol (invalid UUID): %s", req.ProtocolID))		
	}else{
		err = c.Db.AddArticleReferenceToProtocol(r.Context(), database.AddArticleReferenceToProtocolParams{
			ProtocolID: proto_id,
			ReferenceID:  article.ID,
		})
		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding article reference to protocol: %s", req.ProtocolID))
		}
	}

	req = ArticleRefReq{
		ID: article.ID.String(),
		Title: article.Title,
		Authors: article.Authors,
		Journal: article.Journal,
		Year: article.Year,
		Pmid: article.Pmid,
		Doi: article.Doi,
		ProtocolID: req.ProtocolID,
	}		
	
	json_utils.RespondWithJSON(w, http.StatusCreated, req)
}

func HandleGetArticleRefByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	parsed_id, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw_article, err := c.Db.GetArticleReferenceByIDWithProtocols(ctx, parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting article: %s", parsed_id.String()))
		return
	}
	
	linkedProtocols, err := ConvertTuplesToStructs[LinkedProtocols](raw_article.ProtocolIds)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting linked protocols: %s", parsed_id.String()))
		return
	}

	article := ArticleReferenceResponse{
		Id: raw_article.ID,
		Title: raw_article.Title,
		Authors: raw_article.Authors,
		Journal: raw_article.Journal,
		Year: raw_article.Year,
		Pmid: raw_article.Pmid,
		Doi: raw_article.Doi,
		LinkedProtocols: linkedProtocols,
	}

	json_utils.RespondWithJSON(w, http.StatusOK, article)
}

func HandleDeleteArticleRefByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	parsed_id, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteArticleReference(ctx, parsed_id)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting article reference: %s", parsed_id.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "article reference deleted"})
}

func HandleAddArticleToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsed_id, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	proto_id := r.URL.Query().Get("protocol_id")

	parsed_pid, err := uuid.Parse(proto_id)
    if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest,"protocol_id is not a valid uuid")
		return       
    }	

	err = c.Db.AddArticleReferenceToProtocol(ctx, database.AddArticleReferenceToProtocolParams{
		ReferenceID: parsed_id,
		ProtocolID: parsed_pid,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding reference to protocol: %s", parsed_id.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "article reference added to protocol"})

}

func HandleRemoveArticleFromProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsed_id, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	proto_id := r.URL.Query().Get("protocol_id")

	parsed_pid, err := uuid.Parse(proto_id)
    if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest,"protocol_id is not a valid uuid")
		return       
    }	

	err = c.Db.RemoveArticleReferenceFromProtocol(ctx, database.RemoveArticleReferenceFromProtocolParams{
		ReferenceID: parsed_id,
		ProtocolID: parsed_pid,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing reference from protocol: %s", parsed_id.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "article reference removed from protocol"})

}