package portal

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"

	"encoding/base64"

	"k8s.io/client-go/rest"
	componentconfig "kope.io/auth/pkg/apis/componentconfig/v1alpha1"
	authclient "kope.io/auth/pkg/client/clientset_generated/clientset"
	"kope.io/auth/pkg/keystore"
	oauth2proxy "kope.io/auth/pkg/oauth2proxy"
	"kope.io/auth/pkg/tokenstore"
)

type HTTPServer struct {
	options *componentconfig.AuthConfiguration

	listen    string
	staticDir string

	OAuthProxy *oauth2proxy.OAuthProxy
	Tokenstore tokenstore.Interface
}

func NewHTTPServer(o *componentconfig.AuthConfiguration, authProviders []componentconfig.AuthProvider, listen string, staticDir string, cookieSecret keystore.SharedSecretSet) (*HTTPServer, error) {
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
	authClient, err := authclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error building user client: %v", err)
	}

	tokenStore := tokenstore.NewAPITokenStore(authClient)

	b := oauth2proxy.NewOptions()
	b.HttpAddress = listen
	b.CookieName = "_auth_portal"

	if len(authProviders) == 0 {
		return nil, fmt.Errorf("AuthProvider must be configured")
	}
	if len(authProviders) != 1 {
		return nil, fmt.Errorf("Only a single AuthProvider is currently supported")
	}

	var validator func(string) bool

	for _, authProvider := range authProviders {
		if authProvider.OAuthConfig.ClientID == "" {
			return nil, fmt.Errorf("OAuthConfig ClientID not set for %q", authProvider.Name)
		}
		if authProvider.OAuthConfig.ClientSecret == "" {
			return nil, fmt.Errorf("OAuthConfig ClientSecret not set for %q", authProvider.Name)
		}
		glog.Warningf("Using static cookie secret")
		// TODO: Implement rotation etc ...pass it down...
		sharedSecret, err := cookieSecret.EnsureSharedSecret()
		if err != nil {
			return nil, fmt.Errorf("error building shared secret: %v", err)
		}
		b.CookieSecret = base64.URLEncoding.EncodeToString(sharedSecret.SecretData())

		b.ClientID = authProvider.OAuthConfig.ClientID
		b.ClientSecret = authProvider.OAuthConfig.ClientSecret

		validator, err = buildValidator(authProvider.PermitEmails)
		if err != nil {
			return nil, fmt.Errorf("error building validator: %v", err)
		}
		b.EmailDomains = authProvider.PermitEmails
	}

	// Refresh cookies every hour
	b.CookieRefresh = time.Hour

	// Dummy values to pass validation
	b.Upstreams = []string{"http://127.0.0.1:8888"}

	if err := b.Validate(); err != nil {
		return nil, fmt.Errorf("Configuration error: %v", err)
	}

	proxy := oauth2proxy.NewOAuthProxy(b, validator)

	s := &HTTPServer{
		options: o,

		listen:    listen,
		staticDir: staticDir,

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
	mux.HandleFunc("/oauth2/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	mux.HandleFunc("/api/whoami", s.apiWhoAmI)
	mux.HandleFunc("/api/tokens", s.apiTokens)
	mux.HandleFunc("/api/kubeconfig", s.apiKubeconfig)
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	mux.HandleFunc("/portal/actions/login", s.portalActionLogin)
	mux.HandleFunc("/portal/actions/logout", s.portalActionLogout)
	mux.HandleFunc("/portal/actions/kubeconfig", s.portalActionKubeconfig)
	mux.HandleFunc("/portal/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	mux.HandleFunc("/", s.portalIndex)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(s.staticDir))))

	server := &http.Server{
		Addr:    s.listen,
		Handler: mux,
	}
	return server.ListenAndServe()
}

func buildValidator(permitEmails []string) (func(string) bool, error) {
	var allowAll bool
	var exact []string
	var suffixes []string
	for _, permitEmail := range permitEmails {
		wildcardCount := strings.Count(permitEmail, "*")
		if wildcardCount == 0 {
			if permitEmail == "" {
				// TODO: Move to validation?
				// TODO: Maybe ignore invalid rules?
				return nil, fmt.Errorf("empty permitEmail not allowed")
			}
			exact = append(exact, permitEmail)
		} else if wildcardCount == 1 && strings.HasPrefix(permitEmail, "*") {
			if permitEmail == "*" {
				allowAll = true
			} else {
				// TODO: Block dangerous things i.e. require *@ or *. ?
				suffixes = append(suffixes, permitEmail[1:])
			}
		} else {
			return nil, fmt.Errorf("Cannot parse permitEmail rule: %q", permitEmail)
		}
	}

	validator := func(email string) bool {
		if email == "" {
			return false
		}
		email = strings.TrimSpace(strings.ToLower(email))
		if allowAll {
			return true
		}
		for _, s := range exact {
			if s == email {
				return true
			}
		}
		for _, suffix := range suffixes {
			if strings.HasSuffix(email, suffix) {
				return true
			}
		}

		return false
	}
	return validator, nil
}
