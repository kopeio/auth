package portal

import (
	"fmt"
	"k8s.io/client-go/rest"
	"kope.io/auth/pkg/apis/auth"
	oauth2proxy "kope.io/auth/pkg/oauth2proxy"
	"kope.io/auth/pkg/tokenstore"
	"net/http"
	"strings"
	"time"
)

type HTTPServer struct {
	options    *Options
	OAuthProxy *oauth2proxy.OAuthProxy
	Tokenstore tokenstore.Interface
}

func NewHTTPServer(o *Options) (*HTTPServer, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("error building kubernetes configuration: %v", err)
	}
	//// creates the clientset
	//clientset, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	return nil, fmt.Errorf("error building kubernetes client: %v", err)
	//}
	authClient, err := auth.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error building auth client: %v", err)
	}

	tokenStore := tokenstore.NewThirdPartyStorage(authClient)

	b := oauth2proxy.NewOptions()
	b.EmailDomains = o.EmailDomains
	b.HttpAddress = o.Listen
	b.CookieName = "_auth_portal"

	b.CookieSecret = o.CookieSecret
	b.ClientID = o.ClientID
	b.ClientSecret = o.ClientSecret

	// Refresh cookies every hour
	b.CookieRefresh = time.Hour

	// Dummy values to pass validation
	b.Upstreams = []string{"http://127.0.0.1:8888"}

	if err := b.Validate(); err != nil {
		return nil, fmt.Errorf("Configuration error: %v", err)
	}

	validator := buildValidator(o.EmailDomains)

	proxy := oauth2proxy.NewOAuthProxy(b, validator)

	s := &HTTPServer{
		options:    o,
		OAuthProxy: proxy,
		Tokenstore: tokenStore,
	}

	return s, nil
}

func (s *HTTPServer) ListenAndServe() error {
	stopCh := make(chan struct{})
	go s.Tokenstore.Run(stopCh)

	mux := http.NewServeMux()

	mux.HandleFunc("/oauth2/start", s.OAuthProxy.OAuthStart)
	mux.HandleFunc("/oauth2/callback", s.oauthCallback)
	mux.HandleFunc("/oauth2/", func(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) })

	mux.HandleFunc("/api/whoami", s.apiWhoAmI)
	mux.HandleFunc("/api/tokens", s.apiTokens)
	mux.HandleFunc("/api/kubeconfig", s.apiKubeconfig)
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) })

	mux.HandleFunc("/portal/actions/login", s.portalActionLogin)
	mux.HandleFunc("/portal/actions/logout", s.portalActionLogout)
	mux.HandleFunc("/portal/actions/kubeconfig", s.portalActionKubeconfig)
	mux.HandleFunc("/portal/", func(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) })

	mux.HandleFunc("/", s.portalIndex)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(s.options.StaticDir))))

	server := &http.Server{
		Addr:    s.options.Listen,
		Handler: mux,
	}
	return server.ListenAndServe()
}

func buildValidator(domains []string) func(string) bool {
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
