package portal

import (
	"github.com/golang/glog"
	"net/http"
)

func (s *HTTPServer) portalActionLogin(rw http.ResponseWriter, req *http.Request) {
	code, err := s.oauthServer.OAuthStart(rw, req)
	if err != nil {
		s.internalError(rw, req, err)
		return
	}

	if code != 0 {
		glog.Warningf("error from login action %d: %v", code, err)
		rw.WriteHeader(code)
		return
	}
}
