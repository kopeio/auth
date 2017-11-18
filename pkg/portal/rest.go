package portal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/golang/glog"
	auth "kope.io/auth/pkg/apis/auth/v1alpha1"
	"kope.io/auth/pkg/kubeconfig"
)

type UserInfo struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

func (s *HTTPServer) authenticate(rw http.ResponseWriter, req *http.Request) (*auth.User, error) {
	// TODO: Cache in context ?
	session, err := s.oauthServer.Authenticate(rw, req)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}

	//id := s.mapToIdentity(session)
	//glog.Infof("looking up identity %q", id)

	u, err := s.tokenStore.FindUserByUID(session.KubernetesUid)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %v", err)
	}
	if u == nil {
		glog.Infof("user not found %q", session.KubernetesUid)
		return nil, nil
	}

	return u, nil
}

//func (s *HTTPServer) mapToIdentity(session *session.Session) (*auth.IdentitySpec) {
//	providerID := session.SessionData.ProviderId
//
//	// TODO: Store all information from provider?
//	providerUserID := session.SessionData.UserId
//
//	id := &auth.IdentitySpec{
//		ProviderID: providerID,
//		ID:         providerUserID,
//		Username:   providerUserID,
//	}
//
//	return id
//}

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
	session, err := s.oauthServer.Authenticate(rw, req)
	if err != nil {
		s.internalError(rw, req, err)
	}
	if session == nil {
		http.Error(rw, "unauthorized request", http.StatusUnauthorized)
		return
	}

	//id := s.mapToIdentity(session)
	//glog.Infof("looking up identity %q", id)

	u, err := s.tokenStore.FindUserByUID(session.KubernetesUid)
	if err != nil {
		glog.Infof("error finding user: %v", err)
		s.internalError(rw, req, err)
		return
	}
	if u == nil {
		glog.Infof("user not found %q", session.KubernetesUid)
		http.NotFound(rw, req)
		return
	}

	// strings.Replace(session.Email, "@", "-", -1)

	if req.Method == "POST" {
		hashed := false // really hard to use
		tokenSpec, err := s.tokenStore.CreateToken(u, hashed)
		if err != nil {
			s.internalError(rw, req, err)
			return
		}
		response := &TokenResponse{
			Value: kubeconfig.EncodeToken(u, tokenSpec),
		}
		s.sendJson(rw, req, response)
		return
	}

	if req.Method == "GET" {
		response := []*TokenResponse{}
		for _, t := range u.Spec.Tokens {
			response = append(response, &TokenResponse{
				Value: kubeconfig.EncodeToken(u, t),
			})
		}
		s.sendJson(rw, req, response)
		return
	}

	s.methodNotAllowed(rw, req)
}
