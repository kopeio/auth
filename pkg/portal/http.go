package portal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
	"kope.io/auth/pkg/configreader"
	"kope.io/auth/pkg/keystore"
	"kope.io/auth/pkg/oauth"
	"kope.io/auth/pkg/tokenstore"
)

type HTTPServer struct {
	config *configreader.ManagedConfiguration

	listen    string
	staticDir string

	oauthServer *oauth.Server
	tokenStore  tokenstore.Interface
}

func NewHTTPServer(config *configreader.ManagedConfiguration, listen string, staticDir string, keyStore keystore.KeyStore, tokenStore tokenstore.Interface) (*HTTPServer, error) {
	// creates the in-cluster config
	//config, err := rest.InClusterConfig()
	//if err != nil {
	//	return nil, fmt.Errorf("error building kubernetes configuration: %v", err)
	//}
	//// creates the clientset
	//clientset, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	return nil, fmt.Errorf("error building kubernetes client: %v", err)
	//}
	//authClient, err := authclient.NewForConfig(config)
	//if err != nil {
	//	return nil, fmt.Errorf("error building user client: %v", err)
	//}

	//tokenStore := tokenstore.NewAPITokenStore(authClient)

	//b := &oauth.Server{}
	//b.HttpAddress = listen
	//b.CookieName = "_auth_portal"

	//if len(authProviders) == 0 {
	//	return nil, fmt.Errorf("AuthProvider must be configured")
	//}
	//if len(authProviders) != 1 {
	//	return nil, fmt.Errorf("Only a single AuthProvider is currently supported")
	//}
	//
	//var validator func(string) bool
	//
	//for _, authProvider := range authProviders {
	//	if authProvider.OAuthConfig.ClientID == "" {
	//		return nil, fmt.Errorf("OAuthConfig ClientID not set for %q", authProvider.Name)
	//	}
	//	if authProvider.OAuthConfig.ClientSecret == "" {
	//		return nil, fmt.Errorf("OAuthConfig ClientSecret not set for %q", authProvider.Name)
	//	}
	//	glog.Warningf("Using static cookie secret")
	//	// TODO: Implement rotation etc ...pass it down...
	//	sharedSecret, err := cookieSecret.EnsureSharedSecret()
	//	if err != nil {
	//		return nil, fmt.Errorf("error building shared secret: %v", err)
	//	}
	//	b.CookieSecret = base64.URLEncoding.EncodeToString(sharedSecret.SecretData())
	//
	//	b.ClientID = authProvider.OAuthConfig.ClientID
	//	b.ClientSecret = authProvider.OAuthConfig.ClientSecret
	//
	//	validator, err = buildValidator(authProvider.PermitEmails)
	//	if err != nil {
	//		return nil, fmt.Errorf("error building validator: %v", err)
	//	}
	//	b.EmailDomains = authProvider.PermitEmails
	//}
	//
	//// Refresh cookies every hour
	//b.CookieRefresh = time.Hour
	//
	//// Dummy values to pass validation
	//b.Upstreams = []string{"http://127.0.0.1:8888"}
	//
	//if err := b.Validate(); err != nil {
	//	return nil, fmt.Errorf("Configuration error: %v", err)
	//}
	//
	//proxy := oauth2proxy.NewOAuthProxy(b, validator)

	keyset, err := keyStore.KeySet("oauth")
	if err != nil {
		return nil, fmt.Errorf("error initializing keyset: %v", err)
	}

	oauthServer := &oauth.Server{
		CookieName:    "_auth_portal",
		CookieExpiry:  time.Duration(168) * time.Hour,
		CookieRefresh: time.Duration(0),
		Keyset:        keyset,
		Config:        config,
		// UserMapper set below
	}

	s := &HTTPServer{
		config: config,

		listen:    listen,
		staticDir: staticDir,

		oauthServer: oauthServer,
		tokenStore:  tokenStore,
	}

	s.oauthServer.UserMapper = s.mapUser

	return s, nil
}

func (s *HTTPServer) ListenAndServe() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/oauth2/start", s.oauthStart)
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

	staticServer := http.FileServer(http.Dir(s.staticDir))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		glog.Infof("%s %s", r.Method, path)
		if path == "/" {
			s.portalIndex(w, r)
		} else {
			staticServer.ServeHTTP(w, r)
		}
	})

	server := &http.Server{
		Addr:    s.listen,
		Handler: mux,
	}
	return server.ListenAndServe()
}
