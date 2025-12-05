package router

import "net/http"

// @Summary      List contributors
// @Description  Returns a paginated list of contributors
// @Tags         contributors
// @Produce      json
// @Param        page       query  int  false  "Page number"     default(1)
// @Param        page_size  query  int  false  "Items per page"  default(20)
// @Success      200  {object}  ContributorsPageResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /contributors [get]
func (s *server) handleListContributors(w http.ResponseWriter, r *http.Request) {

}
