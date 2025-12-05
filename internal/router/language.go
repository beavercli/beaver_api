package router

import "net/http"

// @Summary      List languages
// @Description  Returns a paginated list of programming languages
// @Tags         languages
// @Produce      json
// @Param        page       query  int  false  "Page number"     default(1)
// @Param        page_size  query  int  false  "Items per page"  default(20)
// @Success      200  {object}  LanguagesPageResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /languages [get]
func (s *server) handleListLanguages(w http.ResponseWriter, r *http.Request) {

}
