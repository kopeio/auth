package oauth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/types"
	"kope.io/auth/pkg/configreader"
	"kope.io/auth/pkg/keystore"
	"kope.io/auth/pkg/oauth/providers"
	"kope.io/auth/pkg/oauth/session"
)

type UserMapperFn func(session *session.Session, info *session.UserInfo) (types.UID, error)

type Server struct {
	CookieName    string
	CookieExpiry  time.Duration
	CookieRefresh time.Duration

	Keyset keystore.KeySet

	UserMapper UserMapperFn
	Config     *configreader.ManagedConfiguration

	providersMutex sync.Mutex
	providers      map[string]providers.Provider
}

func (p *Server) OAuthStart(rw http.ResponseWriter, req *http.Request) (int, error) {
	err := req.ParseForm()
	if err != nil {
		return http.StatusBadRequest, err
	}

	providerID := req.FormValue("provider")
	if providerID == "" {
		return http.StatusBadRequest, fmt.Errorf("provider must be specified")
	}

	redirect := req.FormValue("rd")
	if redirect == "" {
		redirect = "/"
	}

	state := &State{}
	state.ProviderId = providerID
	state.Redirect = redirect
	stateString, err := state.Marshal()
	if err != nil {
		return 0, fmt.Errorf("error marshalling state: %v", err)
	}

	redirectURI := p.GetRedirectURI(req.Host)

	provider, err := p.getProvider(providerID)
	if err != nil {
		return 0, err
	}
	if provider == nil {
		return http.StatusBadRequest, fmt.Errorf("provider %q not configured", providerID)
	}

	http.Redirect(rw, req, provider.GetLoginURL(redirectURI, stateString), 302)
	return 0, nil
}

func (p *Server) Logout(rw http.ResponseWriter, req *http.Request) {
	p.clearCookie(rw, req)
}

//func (p *Server) GetRedirect(req *http.Request) (string, error) {
//	err := req.ParseForm()
//	if err != nil {
//		return "", err
//	}
//
//	redirect := req.FormValue("rd")
//	if redirect == "" {
//		redirect = "/"
//	}
//
//	return redirect, err
//}

func (p *Server) GetRedirectURI(host string) string {
	//// default to the request Host if not set
	//if p.redirectURL.Host != "" {
	//	return p.redirectURL.String()
	//}
	//var u url.URL
	//u = *p.redirectURL
	//if u.Scheme == "" {
	//	if p.CookieSecure {
	//		u.Scheme = "https"
	//	} else {
	//		u.Scheme = "http"
	//	}
	//}
	//u.Host = host
	//return u.String()

	var u url.URL
	u.Scheme = "https"
	u.Host = host

	u.Path = "/oauth2/callback"

	return u.String()
}

func getRemoteAddr(req *http.Request) (s string) {
	s = req.RemoteAddr
	if req.Header.Get("X-Real-IP") != "" {
		s += fmt.Sprintf(" (%q)", req.Header.Get("X-Real-IP"))
	}
	return
}

func (p *Server) OAuthCallback(rw http.ResponseWriter, req *http.Request) (int, error) {
	remoteAddr := getRemoteAddr(req)

	// finish the oauth cycle
	err := req.ParseForm()
	if err != nil {
		return 0, err
	}
	errorString := req.Form.Get("error")
	if errorString != "" {
		return 403, fmt.Errorf("permission denied: %v", errorString)
	}
	state, err := unmarshalState(req.Form.Get("state"))
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid state parameter")
	}
	providerID := state.ProviderId
	redirect := state.Redirect
	if !strings.HasPrefix(redirect, "/") {
		redirect = "/"
	}

	provider, err := p.getProvider(providerID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if provider == nil {
		return http.StatusBadRequest, fmt.Errorf("provider %q not configured", providerID)
	}

	code := req.Form.Get("code")
	if code == "" {
		return http.StatusBadRequest, errors.New("missing code")
	}

	redirectURI := p.GetRedirectURI(req.Host)
	session, info, err := provider.Redeem(redirectURI, code)
	if err != nil {
		return 0, fmt.Errorf("error redeeming code: %v", err)
	}

	// set cookie, or deny
	uid, err := p.UserMapper(session, info)
	if err != nil {
		glog.Infof("%s error mapping to user: %s", remoteAddr, err)
		return 0, err
	}

	if uid == "" {
		glog.Infof("%s Permission Denied: %q is unauthorized", remoteAddr, info)
		return 403, nil
	}

	session.SessionData.KubernetesUid = string(uid)

	glog.Infof("%s authentication complete %s", remoteAddr, session)

	//id := p.MapToIdentity(session)
	//
	//_, err := tokenStore.MapToUser(id, true)
	//if err != nil {
	//	glog.Infof("%s error mapping to user: %s", remoteAddr, err)
	//	return 0, err
	//}

	err = p.saveSession(rw, req, session)
	if err != nil {
		return 0, err
	}

	http.Redirect(rw, req, redirect, 302)
	return 0, nil
}

//func (p *Server) saveSession(rw http.ResponseWriter, req *http.Request, s *pb.SessionData) error {
//	value, err := p.provider.CookieForSession(s, p.CookieCipher)
//	if err != nil {
//		return err
//	}
//	p.SetCookie(rw, req, value)
//	return nil
//}
