package portal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"github.com/golang/glog"
	auth "kope.io/auth/pkg/apis/auth/v1alpha1"
	"kope.io/auth/pkg/kubeconfig"
)

func (s *HTTPServer) portalActionKubeconfig(rw http.ResponseWriter, req *http.Request) {
	// TODO: Is this safe against XSS attacks?
	glog.Errorf("TODO: Validate portalActionKubeconfig against XSS attacks")

	if req.Method == "GET" {
		authn, err := s.authenticate(rw, req)
		if err != nil {
			s.internalError(rw, req, err)
			return
		}
		if authn == nil {
			http.Error(rw, "unauthorized request", http.StatusUnauthorized)
			return
		}

		token, err := s.createToken(authn)
		if err != nil {
			s.internalError(rw, req, err)
			return
		}
		if token == nil {
			http.NotFound(rw, req)
			return
		}

		authConfiguration, err := s.config.AuthConfiguration()
		if err != nil {
			s.internalError(rw, req, err)
			return
		}
		name := authConfiguration.GenerateKubeconfig.Name
		apiEndpoint := authConfiguration.GenerateKubeconfig.Server

		if apiEndpoint == "" && name != "" {
			// Try to infer the apiEndpoint from the name (follow the kops convention)
			apiEndpoint = "https://api." + name
		}
		if name == "" && apiEndpoint != "" {
			// Try to infer the name from the apiEndpoint
			u, err := url.Parse(apiEndpoint)
			if err != nil {
				glog.Warningf("error parsing api endpoint %q", apiEndpoint)
			} else {
				name = u.Host
				name = strings.TrimPrefix(name, "api.")
			}
		}

		if name == "" {
			s.internalError(rw, req, fmt.Errorf("cannot determine cluster name"))
			return
		}

		// TODO: What about if we use a letsencrypt cert?
		caCertificate, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
		if err != nil {
			if os.IsNotExist(err) {
				glog.Warningf("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt does not exist")
			} else {
				glog.Warningf("error reading /var/run/secrets/kubernetes.io/serviceaccount/ca.crt: %v", err)
			}
			caCertificate = nil
		}


		kubeconfigData, err := kubeconfig.BuildKubeconfig(apiEndpoint, caCertificate, authn, token)
		if err != nil {
			s.internalError(rw, req, err)
			return
		}

		rw.Header().Set("Content-Disposition", "attachment; filename=\""+"kubeconfig."+name+"\"")
		rw.Header().Set("Content-Type", "application/octet-stream")

		s.writeResponse(rw, req, []byte(kubeconfigData))
		return
	}

	s.methodNotAllowed(rw, req)
}

func (s *HTTPServer) apiKubeconfig(rw http.ResponseWriter, req *http.Request) {
	// TODO: Is this safe against XSS attacks?
	glog.Errorf("TODO: Validate apiKubeconfig against XSS attacks")

	if req.Method == "GET" {
		user, err := s.authenticate(rw, req)
		if err != nil {
			s.internalError(rw, req, err)
			return
		}

		if user == nil {
			http.Error(rw, "unauthorized request", http.StatusUnauthorized)
			return
		}

		token, err := s.createToken(user)
		if err != nil {
			s.internalError(rw, req, err)
			return
		}

		if token == nil {
			http.NotFound(rw, req)
			return
		}

		response := kubeconfig.EncodeToken(user, token)

		s.writeResponse(rw, req, []byte(response))
		return
	}

	s.methodNotAllowed(rw, req)
}

func (s *HTTPServer) createToken(user *auth.User) (*auth.TokenSpec, error) {
	bestToken := kubeconfig.FindBestToken(user)

	if bestToken == nil {
		hashed := false // really hard to use
		tokenSpec, err := s.tokenStore.CreateToken(user, hashed)
		if err != nil {
			return nil, fmt.Errorf("error creating token: %v", err)
		}
		bestToken = tokenSpec
	}

	if bestToken == nil {
		glog.Infof("could not find token for user %q", user.Spec.Username)
		return nil, nil
	}

	return bestToken, nil
}
