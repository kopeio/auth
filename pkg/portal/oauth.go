package portal

import "net/http"

func (s *HTTPServer) oauthCallback(rw http.ResponseWriter, req *http.Request) {
	code, err := s.OAuthProxy.OAuthCallback(rw, req, s.Tokenstore)
	if err != nil {
		s.internalError(rw, req, err)
		return
	}

	if code != 0 {
		http.Error(rw, "", code)
	}
}
