package oauth

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"time"

	"k8s.io/klog/v2"

	"kope.io/auth/pkg/session"
)

func (p *Server) loadCookiedSession(req *http.Request) (*session.Session, error) {
	c, err := req.Cookie(p.CookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, nil
		} else {
			return nil, fmt.Errorf("error reading cookie %q: %w", p.CookieName, err)
		}
	}

	sessionData, err := p.decrypt(c)
	if err != nil {
		return nil, err
	}

	return sessionData, nil
}

func (p *Server) decrypt(httpCookie *http.Cookie) (*session.Session, error) {
	valueBytes, err := base64.URLEncoding.DecodeString(httpCookie.Value)
	if err != nil {
		return nil, fmt.Errorf("cookie not valid base64 encoded: %w", err)
	}

	plaintext, err := p.Keyset.Decrypt(valueBytes)
	if err != nil {
		return nil, fmt.Errorf("error decrypting cookie: %w", err)
	}

	return session.UnmarshalSession(plaintext)
}

func (p *Server) makeCookie(req *http.Request, value []byte, expiration time.Duration, now time.Time) *http.Cookie {
	domain := req.Host
	if h, _, err := net.SplitHostPort(domain); err == nil {
		domain = h
	}
	//if p.CookieDomain != "" {
	//	if !strings.HasSuffix(domain, p.CookieDomain) {
	//		log.Printf("Warning: request host is %q but using configured cookie domain of %q", domain, p.CookieDomain)
	//	}
	//	domain = p.CookieDomain
	//}

	valueString := ""
	if value != nil {
		valueString = base64.URLEncoding.EncodeToString(value)
	}
	return &http.Cookie{
		Name:     p.CookieName,
		Value:    valueString,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true, //p.CookieHttpOnly,
		Secure:   true, // p.CookieSecure,
		Expires:  now.Add(expiration),
	}
}

func (p *Server) clearCookie(rw http.ResponseWriter, req *http.Request) {
	klog.Infof("clearing cookie")
	http.SetCookie(rw, p.makeCookie(req, nil, time.Hour*-1, time.Now()))
}

func (p *Server) setCookie(rw http.ResponseWriter, req *http.Request, value []byte) {
	http.SetCookie(rw, p.makeCookie(req, value, p.CookieExpiry, time.Now()))
}

func (p *Server) saveSession(rw http.ResponseWriter, req *http.Request, s *session.Session) error {
	plaintext, err := s.Marshal()
	if err != nil {
		return err
	}

	encrypted, err := p.Keyset.Encrypt(plaintext)
	if err != nil {
		return fmt.Errorf("error encrypting cookie: %w", err)
	}

	p.setCookie(rw, req, encrypted)
	return nil
}
