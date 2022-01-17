package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"k8s.io/klog/v2"
	"kope.io/auth/pkg/config"
	"kope.io/auth/pkg/keystore"
	"kope.io/auth/pkg/oauth"
	//"kope.io/auth/pkg/tokenstore"
)

type HTTPServer struct {
	config config.Provider

	staticDir string

	oauthServer *oauth.Server
	//tokenStore  tokenstore.Interface
}

func NewHTTPServer(ctx context.Context, config config.Provider, keyStore keystore.KeyStore) (*HTTPServer, error) {
	keyset, err := keyStore.KeySet(ctx, "oauth")
	if err != nil {
		return nil, fmt.Errorf("error initializing keyset: %w", err)
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

		// staticDir: staticDir,

		oauthServer: oauthServer,
		//tokenStore:  tokenStore,
	}

	s.oauthServer.UserMapper = s.mapUser

	return s, nil
}

func (s *HTTPServer) ListenAndServe(listen string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/oauth2/start", s.oauthStart)
	mux.HandleFunc("/oauth2/callback", s.oauthCallback)
	mux.HandleFunc("/oauth2/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// mux.HandleFunc("/api/whoami", s.apiWhoAmI)
	// mux.HandleFunc("/api/tokens", s.apiTokens)
	// mux.HandleFunc("/api/kubeconfig", s.apiKubeconfig)
	// mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.NotFound(w, r)
	// })

	// mux.HandleFunc("/portal/actions/login", s.portalActionLogin)
	// mux.HandleFunc("/portal/actions/logout", s.portalActionLogout)
	// mux.HandleFunc("/portal/actions/kubeconfig", s.portalActionKubeconfig)
	// mux.HandleFunc("/portal/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.NotFound(w, r)
	// })

	// staticServer := http.FileServer(http.Dir(s.staticDir))
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	path := r.URL.Path
	// 	glog.Infof("%s %s", r.Method, path)
	// 	if path == "/" {
	// 		s.portalIndex(w, r)
	// 	} else {
	// 		staticServer.ServeHTTP(w, r)
	// 	}
	// })

	server := &http.Server{
		Addr:    listen,
		Handler: mux,
	}
	klog.Infof("listening on %s", listen)
	return server.ListenAndServe()
}

func (s *HTTPServer) oauthCallback(rw http.ResponseWriter, req *http.Request) {
	code, err := s.oauthServer.OAuthCallback(rw, req)
	if err != nil {
		s.internalError(rw, req, err)
		return
	}

	if code != 0 {
		http.Error(rw, "", code)
	}
}

func (s *HTTPServer) oauthStart(rw http.ResponseWriter, req *http.Request) {
	code, err := s.oauthServer.OAuthStart(rw, req)
	if err != nil {
		s.internalError(rw, req, err)
		return
	}

	if code != 0 {
		http.Error(rw, "", code)
	}
}

func (s *HTTPServer) internalError(rw http.ResponseWriter, req *http.Request, err error) {
	klog.Warningf("internal error processing %s %s: %v", req.Method, req.URL, err)

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
		klog.Warningf("error sending response: %v", err)
	}
}
