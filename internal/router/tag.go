package router

import "net/http"

// @Summary      List tags
// @Description  Returns a paginated list of tags
// @Tags         tags
// @Produce      json
// @Param        offset  query  int  false  "Offset"  default(0)
// @Param        limit   query  int  false  "Limit"   default(20)
// @Success      200  {array}   Tag
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/tags [get]
func (s *server) handleListTags(w http.ResponseWriter, r *http.Request) {

}
