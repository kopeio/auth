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
