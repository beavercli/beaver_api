package router

import (
	"net/http"

	"github.com/beavercli/beaver_api/internal/service"
)

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
	query := r.URL.Query()
	page, err := toPageQuery(query)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	langList, err := s.service.GetLanguagesPage(r.Context(), service.PageParam{
		Page:     page.Page,
		PageSize: page.PageSize,
	})
	jsonResponse(w, http.StatusOK, toPage(toLanguages(langList.Items), langList.Total, page.Page, page.PageSize))
}
