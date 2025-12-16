package router

import "net/http"

// @Summary      Start GitHub device login
// @Description  Starts the GitHub device OAuth flow and returns verification URL, user code, and polling token
// @Tags         auth
// @Produce      json
// @Success      200  {object}  DeviceOAuth
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /auth/github/login [get]
func (s *server) handleGithubLogin(w http.ResponseWriter, r *http.Request) {
	dr, err := s.service.GetDeviceRequest(r.Context())
	if err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	jsonResponse(w, http.StatusOK, toDeviceOAuth(dr))
}

// @Summary      GitHub callback
// @Description  Handles GitHub OAuth callback and creates user session
// @Tags         auth
// @Param        code   query  string  true  "Authorization code from GitHub"
// @Param        state  query  string  true  "State parameter for CSRF protection"
// @Success      302  "Redirect to frontend"
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /auth/github/callback [get]
func (s *server) handleGithubCallback(w http.ResponseWriter, r *http.Request) {

}

// @Summary      Logout
// @Description  Clears user session and logs out
// @Tags         auth
// @Success      200  {object}  MessageResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /auth/logout [post]
func (s *server) handleLogout(w http.ResponseWriter, r *http.Request) {

}

// @Summary      Get current user
// @Description  Returns the currently authenticated user
// @Tags         auth
// @Produce      json
// @Success      200  {object}  User
// @Failure      401  {object}  ErrorResponse
// @Router       /auth/me [get]
func (s *server) handleMe(w http.ResponseWriter, r *http.Request) {

}
