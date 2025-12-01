package router

import "net/http"

// @Summary      List languages
// @Description  Returns a paginated list of programming languages
// @Tags         languages
// @Produce      json
// @Param        offset  query  int  false  "Offset"  default(0)
// @Param        limit   query  int  false  "Limit"   default(20)
// @Success      200  {array}   Language
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/languages [get]
func (s *server) handleListLanguages(w http.ResponseWriter, r *http.Request) {

}
