package providers

import (
	"net/url"

	"kope.io/auth/pkg/apis/componentconfig/v1alpha1"
	"kope.io/auth/pkg/oauth/session"
)

type Provider interface {
	Config() *v1alpha1.AuthProvider
	//Data() *ProviderData
	//GetEmailAddress(*SessionState) (string, error)
	Redeem(redirectURI string, code string) (*session.Session, *session.UserInfo, error)
	//ValidateGroup(string) bool
	RevalidateSession(*session.Session) (bool, error)
	GetLoginURL(redirectURI, state string) string
	RefreshSessionIfNeeded(*session.Session) (bool, error)
	//SessionFromCookie(string, *cookie.Cipher) (*SessionState, error)
	//CookieForSession(*SessionState, *cookie.Cipher) (string, error)
}

type ProviderConfiguration struct {
	LoginURL       *url.URL
	RedeemURL      *url.URL
	ValidateURL    *url.URL
	Scope          string
	ProviderName   string
	ApprovalPrompt string

	ClientID     string
	ClientSecret string
}
