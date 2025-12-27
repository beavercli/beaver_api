package router

import (
	"net/http"
	"strconv"

	"github.com/beavercli/beaver_api/internal/service"
)

// @Summary		Issue service access token
// @Description	Creates a long-lived token for third-party integrations. The token secret is returned only at creation time.
// @Tags			service-access-tokens
// @Accept			json
// @Produce		json
// @Param			request	body	CreateServiceAccessTokenRequest	true	"Token label and optional expiry"
// @Security		BearerAuth
// @Success		201	{object}	ServiceAccessToken
// @Failure		400	{object}	ErrorResponse
// @Failure		401	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router			/api/v1/service-access-tokens [post]
func (s *server) handleCreateServiceAccessToken(w http.ResponseWriter, r *http.Request) {
	p, err := getCreateServiceAccessTokenRequest(r)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := getUserIDFromCtx(r.Context())
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return

	}
	serviceTokenArgs, err := toServiceCreateServiceAccessToken(p, userID)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	t, err := s.service.CreateServceAccessToken(r.Context(), serviceTokenArgs)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonResponse(w, http.StatusCreated, toServiceAccessToken(t))
}

// @Summary		List service access tokens
// @Description	Returns a paginated list of issued service access tokens without the token secret.
// @Tags			service-access-tokens
// @Produce		json
// @Param			page		query	int	false	"Page number"		default(1)
// @Param			page_size	query	int	false	"Items per page"	default(20)
// @Security		BearerAuth
// @Success		200	{object}	ServiceAccessTokensPageResponse
// @Failure		401	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router			/api/v1/service-access-tokens [get]
func (s *server) handleGetServiceAccessTokens(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	pq, err := toPageQuery(v)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := getUserIDFromCtx(r.Context())
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	sl, err := s.service.ListServiceAccessTokens(r.Context(), userID, service.PageParam{
		Page:     pq.Page,
		PageSize: pq.PageSize,
	})
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, toPage(sl.Items, sl.Total, pq.Page, pq.PageSize))
}

// @Summary		Revoke service access token
// @Description	Revokes an existing service access token by ID.
// @Tags			service-access-tokens
// @Param			ID path	int	true	"Service access token ID"
// @Security		BearerAuth
// @Success		204
// @Failure		400	{object}	ErrorResponse
// @Failure		401	{object}	ErrorResponse
// @Failure		404	{object}	ErrorResponse
// @Failure		500	{object}	ErrorResponse
// @Router			/api/v1/service-access-tokens/{ID} [delete]
func (s *server) handleDeleteServiceAccessToken(w http.ResponseWriter, r *http.Request) {
	tokenID, err := strconv.ParseInt(r.PathValue("ID"), 10, 64)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.service.DeleteServiceAccessToken(r.Context(), tokenID); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonResponse(w, http.StatusCreated, nil)
}
