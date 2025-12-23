package router

import (
	"net/http"

	"github.com/beavercli/beaver_api/internal/service"
)

// @Summary      List contributors
// @Description  Returns a paginated list of contributors
// @Tags         contributors
// @Produce      json
// @Param        page       query  int  false  "Page number"     default(1)
// @Param        page_size  query  int  false  "Items per page"  default(20)
// @Security     BearerAuth
// @Success      200  {object}  ContributorsPageResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/contributors [get]
func (s *server) handleListContributors(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, err := toPageQuery(query)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	contribList, err := s.service.GetContributorsPage(r.Context(), service.PageParam{
		Page:     page.Page,
		PageSize: page.PageSize,
	})
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, toPage(toContributors(contribList.Items), contribList.Total, page.Page, page.PageSize))
}
