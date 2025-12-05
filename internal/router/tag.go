package router

import "net/http"

// @Summary      List tags
// @Description  Returns a paginated list of tags
// @Tags         tags
// @Produce      json
// @Param        page       query  int  false  "Page number"     default(1)
// @Param        page_size  query  int  false  "Items per page"  default(20)
// @Success      200  {object}  TagsPageResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /tags [get]
func (s *server) handleListTags(w http.ResponseWriter, r *http.Request) {

}
