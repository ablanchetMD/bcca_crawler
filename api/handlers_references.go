package api

import (
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"
	"bcca_crawler/internal/json_utils"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"net/http"
)

type ArticleReferenceResponse struct {
	Id              uuid.UUID         `json:"id"`
	CreatedAt       string            `json:"created_at"`
	UpdatedAt       string            `json:"updated_at"`
	Title           string            `json:"title"`
	Authors         string            `json:"authors"`
	Journal         string            `json:"journal"`
	Year            string            `json:"year"`
	Pmid            string            `json:"pmid"`
	Doi             string            `json:"doi"`
	LinkedProtocols []LinkedProtocols `json:"linked_protocols"`
}

type ArticleRefReq struct {
	ID         string `json:"id" validate:"omitempty,uuid"`
	Title      string `json:"title" validate:"required,min=1,max=250"`
	Authors    string `json:"authors" validate:"required,min=1,max=250"`
	Journal    string `json:"journal" validate:"required,min=1,max=250"`
	Year       string `json:"year" validate:"required,min=4,max=4"`
	Pmid       string `json:"pmid" validate:"omitempty,max=25"`
	Doi        string `json:"doi" validate:"omitempty,max=25"`
	ProtocolID string `json:"protocol_id" validate:"omitempty,uuid"`
}

func HandleGetArticleReferences(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	articles := []ArticleReferenceResponse{}
	articleReferences, err := c.Db.GetArticleReferencesWithProtocols(ctx)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting article references")
		return
	}

	for _, a := range articleReferences {

		var linkedProtocols []LinkedProtocols

		protocolIdsBytes, ok := a.ProtocolIds.([]byte)
		if !ok {
			json_utils.RespondWithError(w, http.StatusInternalServerError, "Error asserting protocol IDs to []byte")
			return
		}

		err = json.Unmarshal(protocolIdsBytes, &linkedProtocols)
		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError,
				fmt.Sprintf("Error parsing protocol data: %s", err.Error()))
			return
		}
		articles = append(articles, ArticleReferenceResponse{
			Id:              a.ID,
			CreatedAt:       a.CreatedAt.String(),
			UpdatedAt:       a.UpdatedAt.String(),
			Title:           a.Title,
			Authors:         a.Authors,
			Journal:         a.Journal,
			Year:            a.Year,
			Pmid:            a.Pmid,
			Doi:             a.Doi,
			LinkedProtocols: linkedProtocols,
		})
	}

	json_utils.RespondWithJSON(w, http.StatusOK, articles)

}

func HandleGetArticleRefByProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	items, err := c.Db.GetArticleReferencesByProtocol(ctx, ids.ProtocolID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	json_utils.RespondWithJSON(w, http.StatusOK, items)
}

func HandleUpsertReference(c *config.Config, w http.ResponseWriter, r *http.Request) {

	var req ArticleRefReq

	err := UnmarshalAndValidatePayload(c, r, &req)
	if err != nil {
		println("Error:", err.Error())
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	pid, err := uuid.Parse(req.ID)
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
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error creating reference")
		return
	}
	return_article := ArticleRefReq{
		ID:      article.ID.String(),
		Title:   article.Title,
		Authors: article.Authors,
		Journal: article.Journal,
		Year:    article.Year,
		Pmid:    article.Pmid,
		Doi:     article.Doi,
	}
	if req.ProtocolID == "" {
		json_utils.RespondWithJSON(w, http.StatusOK, return_article)
		return
	}
	proto_id, err := uuid.Parse(req.ProtocolID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding article reference to protocol (invalid UUID): %s", req.ProtocolID))
	} else {
		err = c.Db.AddArticleReferenceToProtocol(r.Context(), database.AddArticleReferenceToProtocolParams{
			ProtocolID:  proto_id,
			ReferenceID: article.ID,
		})
		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding article reference to protocol: %s", req.ProtocolID))
		}
	}
	return_article.ProtocolID = req.ProtocolID

	json_utils.RespondWithJSON(w, http.StatusCreated, return_article)
}

func HandleGetArticleRefByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	parsed_id, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	raw_article, err := c.Db.GetArticleReferenceByIDWithProtocols(ctx, parsed_id.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting article: %s", parsed_id.ID.String()))
		return
	}

	var linkedProtocols []LinkedProtocols

	protocolIdsBytes, ok := raw_article.ProtocolIds.([]byte)
	if !ok {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error asserting protocol IDs to []byte")
		return
	}

	err = json.Unmarshal(protocolIdsBytes, &linkedProtocols)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError,
			fmt.Sprintf("Error parsing protocol data: %s", err.Error()))
		return
	}
	response := ArticleReferenceResponse{
		Id:              raw_article.ID,
		CreatedAt:       raw_article.CreatedAt.String(),
		UpdatedAt:       raw_article.UpdatedAt.String(),
		Title:           raw_article.Title,
		Authors:         raw_article.Authors,
		Journal:         raw_article.Journal,
		Year:            raw_article.Year,
		Pmid:            raw_article.Pmid,
		Doi:             raw_article.Doi,
		LinkedProtocols: linkedProtocols,
	}

	json_utils.RespondWithJSON(w, http.StatusOK, response)
}

func HandleDeleteArticleRefByID(c *config.Config, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	parsed_id, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteArticleReference(ctx, parsed_id.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting article reference: %s", parsed_id.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "article reference deleted"})
}

func HandleAddArticleToProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.AddArticleReferenceToProtocol(ctx, database.AddArticleReferenceToProtocolParams{
		ReferenceID: ids.ID,
		ProtocolID:  ids.ProtocolID,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error adding reference to protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "article reference added to protocol"})

}

func HandleRemoveArticleFromProtocol(c *config.Config, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ids, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.RemoveArticleReferenceFromProtocol(ctx, database.RemoveArticleReferenceFromProtocolParams{
		ReferenceID: ids.ID,
		ProtocolID:  ids.ProtocolID,
	})

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error removing reference from protocol: %s", ids.ID.String()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "article reference removed from protocol"})

}
