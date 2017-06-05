package portal

import (
	"net/http"
)

func (s *HTTPServer) oauthCallback(rw http.ResponseWriter, req *http.Request) {
	code, err := s.oauthServer.OAuthCallback(rw, req)
	if err != nil {
		s.internalError(rw, req, err)
		return
	}

	if code != 0 {
		http.Error(rw, "", code)
	}
}

func (s *HTTPServer) oauthStart(rw http.ResponseWriter, req *http.Request) {
	code, err := s.oauthServer.OAuthStart(rw, req)
	if err != nil {
		s.internalError(rw, req, err)
		return
	}

	if code != 0 {
		http.Error(rw, "", code)
	}
}
