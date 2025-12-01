package router

import "net/http"

// @Summary      Get random snippet
// @Description  Returns a random code snippet, optionally filtered by language or tags
// @Tags         snippets
// @Produce      json
// @Param        language  query  string  false  "Filter by language"
// @Param        tag       query  string  false  "Filter by tag"
// @Success      200  {object}  Snippet
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/snippets/random [get]
func (s *server) handleGetRandomSnippet(w http.ResponseWriter, r *http.Request) {

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
// @Router       /api/v1/snippets [post]
func (s *server) handleCreateSnippet(w http.ResponseWriter, r *http.Request) {

}
