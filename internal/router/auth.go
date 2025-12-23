package router

import (
	"net/http"
	"strconv"
)

// @Summary      Start GitHub device login
// @Description  Starts the GitHub device OAuth flow and returns verification URL, user code, and polling token
// @Tags         auth
// @Produce      json
// @Success      200  {object}  DeviceOAuth
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /auth/github/login [post]
func (s *server) handleGithubLogin(w http.ResponseWriter, r *http.Request) {
	dr, err := s.service.GetDeviceRequest(r.Context())
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, toDeviceOAuth(dr))
}

// @Summary      Complete GitHub device login
// @Description  Exchanges the device flow token for a GitHub user and upserts the user in the database
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      GithubPullRequest  true  "Device flow token returned from /auth/github/login"
// @Success      200      {object}  DeviceAuthResult
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /auth/github/device/poll [post]
func (s *server) handleGitHubDeviceStatus(w http.ResponseWriter, r *http.Request) {
	p, err := toGithubPullRequest(r)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	ar, err := s.service.GithubDevicePoll(r.Context(), p.Token)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, toDeviceAuthResult(ar))
}

// @Summary      Refresh session tokens
// @Description  Rotates a refresh token and returns a new access/refresh pair
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RefreshToken  true  "Refresh token payload"
// @Security     BearerAuth
// @Success      200      {object}  TokenPair
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /auth/refresh [post]
func (s *server) handleTokenRotate(w http.ResponseWriter, r *http.Request) {
	t, err := toRefreshToken(r)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	uID, err := strconv.ParseInt(t.UserID, 10, 64)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return

	}
	tp, err := s.service.RotateTokens(r.Context(), uID, t.RefreshToken)
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, toTokenPair(tp))
}

// @Summary      Logout
// @Description  Clears user session and logs out
// @Tags         auth
// @Security     BearerAuth
// @Success      200  {object}  MessageResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /auth/logout [post]
func (s *server) handleLogout(w http.ResponseWriter, r *http.Request) {

}

// @Summary      Get current user
// @Description  Returns the currently authenticated user
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  User
// @Failure      401  {object}  ErrorResponse
// @Router       /auth/me [get]
func (s *server) handleMe(w http.ResponseWriter, r *http.Request) {

}
