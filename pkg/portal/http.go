package portal

import (
	"net/http"
	"github.com/golang/glog"
	oauth2proxy "github.com/kopeio/kauth/pkg/oauth2proxy"
	"github.com/kopeio/kauth/pkg/tokenstore"
	"fmt"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	"strings"
)

type HTTPServer struct {
	options    *Options
	OAuthProxy *oauth2proxy.OAuthProxy
	Tokenstore tokenstore.Interface
}

func NewHTTPServer(o *Options) (*HTTPServer, error) {
	if o.Namespace == "" {
		return nil, fmt.Errorf("Namespace must be specified (either through the NAMESPACE env var or -namespace flag")
	}

	glog.V(2).Infof("Using namespace %q", o.Namespace)

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("error building kubernetes configuration: %v", err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error building kubernetes client: %v", err)
	}

	tokenStore := tokenstore.NewSecrets(clientset, o.Namespace)

	b := oauth2proxy.NewOptions()
	b.EmailDomains = o.EmailDomains
	b.HttpAddress = o.Listen
	b.CookieName = "_kauth_portal"

	b.CookieSecret = o.CookieSecret
	b.ClientID = o.ClientID
	b.ClientSecret = o.ClientSecret


	// Dummy values to pass validation
	b.Upstreams = []string{"http://127.0.0.1:8888" }

	if err := b.Validate(); err != nil {
		return nil, fmt.Errorf("Configuration error: %v", err)
	}

	validator := buildValidator(o.EmailDomains)

	proxy := oauth2proxy.NewOAuthProxy(b, validator)

	s := &HTTPServer{
		options: o,
		OAuthProxy: proxy,
		Tokenstore: tokenStore,
	}

	return s, nil
}

func (s*HTTPServer) ListenAndServe() error {
	stopCh := make(chan struct{})
	go s.Tokenstore.Run(stopCh)

	mux := http.NewServeMux()

	mux.HandleFunc("/oauth2/start", s.OAuthProxy.OAuthStart)
	mux.HandleFunc("/oauth2/callback", s.OAuthProxy.OAuthCallback)
	mux.HandleFunc("/api/whoami", s.apiWhoAmI)
	mux.HandleFunc("/api/tokens", s.apiTokens)

	server := &http.Server{
		Addr:   s.options.Listen,
		Handler: mux,
	}
	return server.ListenAndServe()
}

func buildValidator(domains []string) func(string) (bool) {
	var allowAll bool
	for i, domain := range domains {
		if domain == "*" {
			allowAll = true
			continue
		}
		domains[i] = fmt.Sprintf("@%s", strings.ToLower(domain))
	}

	validator := func(email string) (valid bool) {
		if email == "" {
			return
		}
		email = strings.ToLower(email)
		for _, domain := range domains {
			valid = valid || strings.HasSuffix(email, domain)
		}
		//if !valid {
		//	valid = validUsers.IsValid(email)
		//}
		if allowAll {
			valid = true
		}
		return valid
	}
	return validator
}