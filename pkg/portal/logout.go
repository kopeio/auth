package portal

import (
	"net/http"
)

func (s *HTTPServer) portalActionLogout(rw http.ResponseWriter, req *http.Request) {
	s.oauthServer.Logout(rw, req)

	http.Redirect(rw, req, "/", http.StatusFound)
}
