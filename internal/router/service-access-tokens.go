package router

import "net/http"

// @Summary      Issue service access token
// @Description  Creates a long-lived token for third-party integrations. The token secret is returned only at creation time.
// @Tags         service-access-tokens
// @Accept       json
// @Produce      json
// @Param        request  body      CreateServiceAccessTokenRequest  true  "Token label and optional expiry"
// @Success      201      {object}  ServiceAccessToken
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /api/v1/service-access-tokens [post]
func (s *server) handleCreateServiceAccessToken(w http.ResponseWriter, r *http.Request) {
	return
}

// @Summary      List service access tokens
// @Description  Returns a paginated list of issued service access tokens without the token secret.
// @Tags         service-access-tokens
// @Produce      json
// @Param        page       query  int  false  "Page number"     default(1)
// @Param        page_size  query  int  false  "Items per page"  default(20)
// @Success      200  {object}  ServiceAccessTokensPageResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/service-access-tokens [get]
func (s *server) handleGetServiceAccessToken(w http.ResponseWriter, r *http.Request) {
	return
}

// @Summary      Revoke service access token
// @Description  Revokes an existing service access token by ID.
// @Tags         service-access-tokens
// @Param        token_id  query  int  true  "Service access token ID"
// @Success      204
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/service-access-tokens [delete]
func (s *server) handleDeleteServiceAccessToken(w http.ResponseWriter, r *http.Request) {
	return
}
