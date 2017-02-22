package portal

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"kope.io/auth/pkg/tokenstore"
	"net/http"
	"kope.io/auth/pkg/apis/auth"
)

type WhoAmIResponse struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

func (s *HTTPServer) apiWhoAmI(rw http.ResponseWriter, req *http.Request) {
	session, status := s.OAuthProxy.Authenticate(rw, req)
	if status == http.StatusAccepted {
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

		response := &WhoAmIResponse{
			ID: u.Metadata.Name,
			Username: u.Spec.Username,
		}
		responseJson, err := json.Marshal(response)
		if err != nil {
			http.Error(rw, fmt.Sprintf("internal error: %v", err), http.StatusInternalServerError)
		} else {
			rw.Write(responseJson)
		}
	} else {
		http.Error(rw, "unauthorized request", http.StatusUnauthorized)
	}
}

func (s *HTTPServer) internalError(rw http.ResponseWriter, req *http.Request, err error) {
	glog.Warningf("internal error processing %s %s: %v", req.Method, req.URL, err)

	http.Error(rw, "internal error", http.StatusInternalServerError)
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
			UserID:  string(u.Metadata.UID),
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
				UserID:  string(u.Metadata.UID),
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

	http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
}

func (s *HTTPServer) apiKubeconfig(rw http.ResponseWriter, req *http.Request) {
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

	if req.Method == "GET" {
		var bestToken *auth.TokenSpec

		for _, t := range u.Spec.Tokens {
			// TODO: Pick best token
			bestToken = t
		}

		if bestToken == nil {
			hashed := false // really hard to use
			tokenSpec, err := s.Tokenstore.CreateToken(u, hashed)
			if err != nil {
				s.internalError(rw, req, err)
				return
			}
			bestToken = tokenSpec
		}

		if bestToken == nil {
			glog.Infof("user not found %q", id.Username)
			http.NotFound(rw, req)
			return
		}

		tokenInfo := &tokenstore.TokenInfo{
			UserID:  string(u.Metadata.UID),
			TokenID: bestToken.ID,
			Secret:  bestToken.RawSecret,
		}

		response := tokenInfo.Encode()
		s.writeResponse(rw, req, []byte(response))
		return
	}

	http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
}
