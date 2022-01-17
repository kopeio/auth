package oauth

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"k8s.io/klog/v2"

	"kope.io/auth/pkg/config"
	"kope.io/auth/pkg/keystore"
	"kope.io/auth/pkg/oauth/providers"
	"kope.io/auth/pkg/session"
)

const stateCookieName = "_oauth2_state"

type UserMapperFn func(ctx context.Context, session *session.Session, info *session.UserInfo) error

type Server struct {
	CookieName    string
	CookieExpiry  time.Duration
	CookieRefresh time.Duration

	Keyset keystore.KeySet

	UserMapper UserMapperFn
	Config     config.Provider

	providersMutex sync.Mutex
	providers      map[string]providers.Provider
}

func (p *Server) OAuthStart(rw http.ResponseWriter, req *http.Request) (int, error) {
	ctx := req.Context()

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
	state.Nonce = strconv.FormatInt(rand.Int63(), 16)

	stateString, err := state.Marshal()
	if err != nil {
		return 0, fmt.Errorf("error marshalling state: %w", err)
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     stateCookieName,
		Value:    stateString,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   600,
	})

	// TODO: SameSite?

	// TODO: Set secure if scheme == https?

	redirectURI := p.GetRedirectURI(req)

	provider, err := p.getProvider(ctx, providerID)
	if err != nil {
		return 0, err
	}
	if provider == nil {
		return http.StatusBadRequest, fmt.Errorf("provider %q not configured", providerID)
	}

	http.Redirect(rw, req, provider.GetLoginURL(ctx, redirectURI, stateString), 302)
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

func (p *Server) GetRedirectURI(req *http.Request) string {
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
	u.Scheme = req.URL.Scheme
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	u.Host = req.Host

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
	ctx := req.Context()

	remoteAddr := getRemoteAddr(req)

	// finish the oauth cycle
	err := req.ParseForm()
	if err != nil {
		return 0, err
	}

	stateCookie, err := req.Cookie(stateCookieName)
	if err != nil {
		return 0, fmt.Errorf("error getting state cookie: %w", err)
	}
	stateParameter := req.URL.Query().Get("state")

	if stateParameter != stateCookie.Value {
		klog.Warningf("state in cookie does not match state in request")
		return 0, fmt.Errorf("state mismatch %q vs %q", stateCookie.Value, stateParameter)
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     stateCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	errorString := req.Form.Get("error")
	if errorString != "" {
		return 403, fmt.Errorf("permission denied: %v", errorString)
	}
	state, err := unmarshalState(stateParameter)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("invalid state parameter")
	}
	providerID := state.ProviderId
	redirect := state.Redirect
	if !strings.HasPrefix(redirect, "/") {
		redirect = "/"
	}

	provider, err := p.getProvider(ctx, providerID)
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

	redirectURI := p.GetRedirectURI(req)
	session, info, err := provider.Redeem(ctx, redirectURI, code)
	if err != nil {
		return 0, fmt.Errorf("error redeeming code: %w", err)
	}

	// set cookie, or deny
	if err := p.UserMapper(ctx, session, info); err != nil {
		klog.Infof("%s error mapping to user: %v", remoteAddr, err)
		return 0, err
	}

	// if uid == "" {
	// 	klog.Infof("%s Permission Denied: %q is unauthorized", remoteAddr, info)
	// 	return 403, nil
	// }

	// session.SessionData.KubernetesUid = string(uid)

	klog.Infof("%s authentication complete %s", remoteAddr, session)

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
