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
// @Security     BearerAuth
// @Success      200  {object}  Snippet
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/snippets/{SnippetID} [get]
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
// @Param        language_id  query  int    false  "Filter by language ID"
// @Param        tag_id       query  []int  false  "Filter by tag IDs (repeat: tag_id=1&tag_id=2)"  collectionFormat(multi)
// @Security     BearerAuth
// @Success      200  {object}  SnippetsPageResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/snippets [get]
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
		PageParam: service.PageParam{
			Page:     p.Page,
			PageSize: p.PageSize,
		},
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
// @Param        body  body  IngestSnippetRequest  true  "Snippet data"
// @Security     BearerAuth
// @Success      201
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/snippets [post]
func (s *server) handleIngestSnippet(w http.ResponseWriter, r *http.Request) {
	p, err := toCreateSnippetRequestBody(r)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())

		return
	}

	if err := s.service.InjestSnippet(r.Context(), toCreateSnippetParams(p)); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}
