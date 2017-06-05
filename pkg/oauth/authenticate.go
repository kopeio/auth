package oauth

import (
	"fmt"
	"github.com/golang/glog"
	"kope.io/auth/pkg/oauth/session"
	"net/http"
	"time"
)

func (p *Server) Authenticate(rw http.ResponseWriter, req *http.Request) (*session.Session, error) {
	var saveSession, clearSession, revalidated bool
	remoteAddr := getRemoteAddr(req)

	session, err := p.loadCookiedSession(req)
	if err != nil {
		glog.Infof("error loading session %s %s", remoteAddr, err)
		session = nil
	}

	if session == nil {
		return nil, nil
	}

	provider, err := p.getProvider(session.SessionData.ProviderId)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		glog.Infof("provider %q not found", session.SessionData.ProviderId)
		return nil, nil
	}

	{
		sessionAge := session.Age()

		if sessionAge > p.CookieRefresh && p.CookieRefresh != time.Duration(0) {
			glog.Infof("%s refreshing %s old session cookie for %s (refresh after %s)", remoteAddr, sessionAge, session, p.CookieRefresh)
			saveSession = true
		}
	}

	if ok, err := provider.RefreshSessionIfNeeded(session); err != nil {
		glog.Infof("%s removing session. error refreshing access token %s %s", remoteAddr, err, session)
		clearSession = true
		session = nil
	} else if ok {
		saveSession = true
		revalidated = true
	}

	if session != nil && session.IsExpired() {
		glog.Infof("%s removing session. token expired %s", remoteAddr, session)
		session = nil
		saveSession = false
		clearSession = true
	}

	if saveSession && !revalidated && session != nil {
		ok, err := provider.RevalidateSession(session)
		if err != nil {
			return nil, fmt.Errorf("error revalidating session: %v", err)
		}

		if !ok {
			glog.Infof("%s removing session. error revalidating %s", remoteAddr, session)
			saveSession = false
			session = nil
			clearSession = true
		}
	}

	//if session != nil && session.Email != "" && !p.Validator(session.Email) {
	//	glog.Infof("%s Permission Denied: removing session %s", remoteAddr, session)
	//	session = nil
	//	saveSession = false
	//	clearSession = true
	//}

	if saveSession && session != nil {
		err := p.saveSession(rw, req, session)
		if err != nil {
			return session, err
		}
	}

	if clearSession {
		p.clearCookie(rw, req)
	}

	//if session == nil {
	//	session, err = p.CheckBasicAuth(req)
	//	if err != nil {
	//		log.Printf("%s %s", remoteAddr, err)
	//	}
	//}

	//if session == nil {
	//	return session, http.StatusForbidden
	//}

	//// At this point, the user is authenticated. proxy normally
	//if p.PassBasicAuth {
	//	req.SetBasicAuth(session.User, p.BasicAuthPassword)
	//	req.Header["X-Forwarded-User"] = []string{session.User}
	//	if session.Email != "" {
	//		req.Header["X-Forwarded-Email"] = []string{session.Email}
	//	}
	//}
	//if p.PassAccessToken && session.AccessToken != "" {
	//	req.Header["X-Forwarded-Access-Token"] = []string{session.AccessToken}
	//}
	//if session.Email == "" {
	//	rw.Header().Set("GAP-Auth", session.User)
	//} else {
	//	rw.Header().Set("GAP-Auth", session.Email)
	//}

	return session, nil
}
