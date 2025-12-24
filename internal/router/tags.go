package router

import (
	"net/http"

	"github.com/beavercli/beaver_api/internal/service"
)

// @Summary		List tags
// @Description	Returns a paginated list of tags
// @Tags			tags
// @Produce		json
// @Param			page		query	int	false	"Page number"		default(1)
// @Param			page_size	query	int	false	"Items per page"	default(20)
// @Security		BearerAuth
// @Success		200	{object}	TagsPageResponse
// @Failure		500	{object}	ErrorResponse
// @Router			/api/v1/tags [get]
func (s *server) handleListTags(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, err := toPageQuery(query)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	tags, err := s.service.GetTagsPage(r.Context(), service.PageParam{
		Page:     page.Page,
		PageSize: page.PageSize,
	})
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, toPage(toTags(tags.Items), tags.Total, page.Page, page.PageSize))
}
