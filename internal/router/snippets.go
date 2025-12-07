package router

import (
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
	id, err := strconv.ParseInt(r.PathValue("SnippetID"), 10, 64)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	snippet, err := s.service.GetSnippet(r.Context(), id)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

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
	v := r.URL.Query()
	p, err := toPageQuery(v)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	f, err := toSnippetListFilterArg(v)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	snippetList, err := s.service.GetSnippetsPage(r.Context(), service.ListSnippetsParams{
		Page:       p.Page,
		PageSize:   p.PageSize,
		LanguageID: f.LanguageID,
		TagIDs:     f.TagIDs,
	})
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	snippetsSummaries := toSnippetSummaries(snippetList.Items)
	snippetPage := toPage(snippetsSummaries, snippetList.Total, p.Page, p.PageSize)
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
