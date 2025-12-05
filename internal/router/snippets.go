package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/beavercli/beaver_api/internal/service"
)

// @Summary      Get snippet by ID
// @Description  Returns a code snippet by its ID
// @Tags         snippets
// @Produce      json
// @Param        id  path  int  true  "Snippet ID"
// @Success      200  {object}  Snippet
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /snippets/{id} [get]
func (s *server) handleGetSnippet(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("SnippetID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid Snippet ID"}); err != nil {
			fmt.Println("Cannot send error response")
		}
		return
	}

	snippet, err := s.service.GetSnippet(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: CONVERT TO RESPONSE DTO
	jsonResponse(w, http.StatusOK, toSnippet(snippet))
}

// @Summary      List snippets
// @Description  Returns a paginated list of snippets with tags and languages
// @Tags         snippets
// @Produce      json
// @Param        page       query  int     false  "Page number"      default(1)
// @Param        page_size  query  int     false  "Items per page"   default(20)
// @Param        language   query  string  false  "Filter by language"
// @Param        tag        query  string  false  "Filter by tag"
// @Success      200  {object}  SnippetsPageResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /snippets [get]
func (s *server) handleListSnippets(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(query.Get("page_size"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var languageID *int64
	if id, err := strconv.ParseInt(query.Get("language"), 10, 64); err == nil {
		languageID = &id
	}

	var tagIDs []int64
	for _, s := range query["tag"] {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			tagIDs = append(tagIDs, id)
		}
	}
	snippetList, err := s.service.GetSnippetsPage(r.Context(), service.ListSnippetsParams{
		Page:       page,
		PageSize:   pageSize,
		LanguageID: languageID,
		TagIDs:     tagIDs,
	})
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	snippetsSummaries := toSnippetSummaries(snippetList.Items)
	snippetPage := toPage(snippetsSummaries, snippetList.Total, page, pageSize)
	jsonResponse(w, http.StatusOK, snippetPage)
}

// @Summary      Create snippet
// @Description  Creates a new code snippet
// @Tags         snippets
// @Accept       json
// @Produce      json
// @Param        body  body  CreateSnippetRequest  true  "Snippet data"
// @Success      201  {object}  Snippet
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /snippets [post]
func (s *server) handleCreateSnippet(w http.ResponseWriter, r *http.Request) {

}
