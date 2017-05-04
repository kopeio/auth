package portal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
	auth "kope.io/auth/pkg/apis/auth/v1alpha1"
	"kope.io/auth/pkg/tokenstore"
)

type UserInfo struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

func (s *HTTPServer) authenticate(rw http.ResponseWriter, req *http.Request) (*auth.User, error) {
	// TODO: Cache in context ?
	session, status := s.OAuthProxy.Authenticate(rw, req)
	if status == http.StatusAccepted {
		id := s.OAuthProxy.MapToIdentity(session)

		glog.Infof("looking up identity %q", id)

		u, err := s.Tokenstore.MapToUser(id, false)
		if err != nil {
			return nil, fmt.Errorf("error finding user: %v", err)
		}
		if u == nil {
			glog.Infof("user not found %q", id.Username)
			return nil, nil
		}

		return u, nil
	} else {
		return nil, nil
	}
}

func (s *HTTPServer) apiWhoAmI(rw http.ResponseWriter, req *http.Request) {
	auth, err := s.authenticate(rw, req)
	if err != nil {
		s.internalError(rw, req, err)
		return
	}

	if auth == nil {
		http.Error(rw, "unauthorized request", http.StatusUnauthorized)
		return
	}

	response := &UserInfo{
		ID:       auth.Name,
		Username: auth.Spec.Username,
	}
	responseJson, err := json.Marshal(response)
	if err != nil {
		s.internalError(rw, req, err)
		return
	}

	rw.Write(responseJson)
}

func (s *HTTPServer) internalError(rw http.ResponseWriter, req *http.Request, err error) {
	glog.Warningf("internal error processing %s %s: %v", req.Method, req.URL, err)

	http.Error(rw, "internal error", http.StatusInternalServerError)
}

func (s *HTTPServer) methodNotAllowed(rw http.ResponseWriter, req *http.Request) {
	http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
}

func (s *HTTPServer) sendJson(rw http.ResponseWriter, req *http.Request, response interface{}) {
	responseJson, err := json.Marshal(response)
	if err != nil {
		s.internalError(rw, req, err)
		return
	}

	s.writeResponse(rw, req, responseJson)
}

func (s *HTTPServer) writeResponse(rw http.ResponseWriter, req *http.Request, data []byte) {
	_, err := rw.Write(data)
	if err != nil {
		glog.Warningf("error sending response: %v", err)
	}
}

type TokenResponse struct {
	Value string `json:"value"`
}

func (s *HTTPServer) apiTokens(rw http.ResponseWriter, req *http.Request) {
	session, status := s.OAuthProxy.Authenticate(rw, req)
	if status != http.StatusAccepted {
		http.Error(rw, "unauthorized request", http.StatusUnauthorized)
		return
	}

	id := s.OAuthProxy.MapToIdentity(session)

	glog.Infof("looking up identity %q", id)

	u, err := s.Tokenstore.MapToUser(id, false)
	if err != nil {
		glog.Infof("error finding user: %v", err)
		s.internalError(rw, req, err)
		return
	}
	if u == nil {
		glog.Infof("user not found %q", id.Username)
		http.NotFound(rw, req)
		return
	}

	// strings.Replace(session.Email, "@", "-", -1)

	if req.Method == "POST" {
		hashed := false // really hard to use
		tokenSpec, err := s.Tokenstore.CreateToken(u, hashed)
		if err != nil {
			s.internalError(rw, req, err)
			return
		}
		tokenInfo := &tokenstore.TokenInfo{
			UserID:  string(u.UID),
			TokenID: tokenSpec.ID,
			Secret:  tokenSpec.RawSecret,
		}
		response := &TokenResponse{
			Value: tokenInfo.Encode(),
		}
		s.sendJson(rw, req, response)
		return
	}

	if req.Method == "GET" {
		response := []*TokenResponse{}
		for _, t := range u.Spec.Tokens {
			tokenInfo := &tokenstore.TokenInfo{
				UserID:  string(u.UID),
				TokenID: t.ID,
				Secret:  t.RawSecret,
			}
			response = append(response, &TokenResponse{
				Value: tokenInfo.Encode(),
			})
		}
		s.sendJson(rw, req, response)
		return
	}

	s.methodNotAllowed(rw, req)
}
