package portal

import (
	"net/http"
)

func (s *HTTPServer) portalActionLogin(rw http.ResponseWriter, req *http.Request) {
	s.OAuthProxy.OAuthStart(rw, req)
}
