package portal

import (
	authenticationv1beta1 "k8s.io/client-go/pkg/apis/authentication/v1beta1"
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/golang/glog"
	"strings"
)

type WhoAmIResponse struct {
	User  string `json:"user,omitempty"`
	Email string `json:"email,omitempty"`
}

func (s*HTTPServer) apiWhoAmI(rw http.ResponseWriter, req *http.Request) {
	session, status := s.OAuthProxy.Authenticate(rw, req)
	if status == http.StatusAccepted {
		response := &WhoAmIResponse{
			User: session.User,
			Email: session.Email,
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

func (s*HTTPServer) internalError(rw http.ResponseWriter, req *http.Request, err error) {
	glog.Warningf("internal error processing %s %s: %v", req.Method, req.URL, err)

	http.Error(rw, "internal error", http.StatusInternalServerError)
}

func (s*HTTPServer) sendJson(rw http.ResponseWriter, req *http.Request, response interface{}) {
	responseJson, err := json.Marshal(response)
	if err != nil {
		s.internalError(rw, req, err)
		return
	}

	_, err = rw.Write(responseJson)
	if err != nil {
		glog.Warningf("error sending response: %v", err)
	}
}

type TokenResponse struct {
	Value string `json:"value"`
}

func (s*HTTPServer) apiTokens(rw http.ResponseWriter, req *http.Request) {
	session, status := s.OAuthProxy.Authenticate(rw, req)
	if status != http.StatusAccepted {
		http.Error(rw, "unauthorized request", http.StatusUnauthorized)
		return
	}

	// TODO: This probably isn't unique / safe
	uid := strings.Replace(session.Email, "@", "-", -1)


	if req.Method == "POST" {
		userInfo := &authenticationv1beta1.UserInfo{
			UID: uid,
			Username: session.Email,
		}
		token, err := s.Tokenstore.CreateToken(userInfo)
		if err != nil {
			s.internalError(rw, req, err)
			return
		}
		response := &TokenResponse{
			Value: token.Encode(),
		}
		s.sendJson(rw, req, response)
		return
	}

	if req.Method == "GET" {
		tokens, err := s.Tokenstore.ListTokens(uid)
		if err != nil {
			s.internalError(rw, req, err)
			return
		}
		response := []*TokenResponse{}
		for _, t := range tokens {
			response = append(response, &TokenResponse{
				Value: t.Encode(),
			})
		}
		s.sendJson(rw, req, response)
		return
	}

	http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
}