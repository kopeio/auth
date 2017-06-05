package portal

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	auth "kope.io/auth/pkg/apis/auth/v1alpha1"
	"kope.io/auth/pkg/tokenstore"
)

type KubeConfig struct {
	Kind           string                    `json:"kind"`
	ApiVersion     string                    `json:"apiVersion"`
	CurrentContext string                    `json:"current-context"`
	Clusters       []*KubectlClusterWithName `json:"clusters"`
	Contexts       []*KubectlContextWithName `json:"contexts"`
	Users          []*KubectlUserWithName    `json:"users"`
}

type KubectlClusterWithName struct {
	Name    string         `json:"name"`
	Cluster KubectlCluster `json:"cluster"`
}

type KubectlCluster struct {
	Server                   string `json:"server,omitempty"`
	CertificateAuthorityData []byte `json:"certificate-authority-data,omitempty"`
}

type KubectlContextWithName struct {
	Name    string         `json:"name"`
	Context KubectlContext `json:"context"`
}

type KubectlContext struct {
	Cluster string `json:"cluster"`
	User    string `json:"user"`
}

type KubectlUserWithName struct {
	Name string      `json:"name"`
	User KubectlUser `json:"user"`
}

type KubectlUser struct {
	ClientCertificateData []byte `json:"client-certificate-data,omitempty"`
	ClientKeyData         []byte `json:"client-key-data,omitempty"`
	Password              string `json:"password,omitempty"`
	Username              string `json:"username,omitempty"`
	Token                 string `json:"token,omitempty"`
}

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

		cluster := KubectlCluster{
			Server: apiEndpoint,
			CertificateAuthorityData: caCertificate,
		}
		context := KubectlContext{
			Cluster: name,
			User:    name,
		}
		user := KubectlUser{
			Token: token.Encode(),
		}
		config := &KubeConfig{
			ApiVersion:     "v1",
			Kind:           "Config",
			CurrentContext: name,
			Clusters: []*KubectlClusterWithName{
				{
					Name:    name,
					Cluster: cluster,
				},
			},
			Contexts: []*KubectlContextWithName{
				{
					Name:    name,
					Context: context,
				},
			},
			Users: []*KubectlUserWithName{
				{
					Name: name,
					User: user,
				},
			},
		}

		response, err := yaml.Marshal(config)
		if err != nil {
			s.internalError(rw, req, fmt.Errorf("error serializing kubeconfig to yaml: %v", err))
			return
		}

		rw.Header().Set("Content-Disposition", "attachment; filename=\""+"kubeconfig."+name+"\"")
		rw.Header().Set("Content-Type", "application/octet-stream")

		s.writeResponse(rw, req, []byte(response))
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

		response := token.Encode()

		s.writeResponse(rw, req, []byte(response))
		return
	}

	s.methodNotAllowed(rw, req)
}

func (s *HTTPServer) createToken(user *auth.User) (*tokenstore.TokenInfo, error) {
	var bestToken *auth.TokenSpec

	for _, t := range user.Spec.Tokens {
		// TODO: Pick best token
		bestToken = t
	}

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

	tokenInfo := &tokenstore.TokenInfo{
		UserID:  string(user.UID),
		TokenID: bestToken.ID,
		Secret:  bestToken.RawSecret,
	}

	return tokenInfo, nil
}
