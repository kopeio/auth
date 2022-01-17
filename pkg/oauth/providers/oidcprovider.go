package providers

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"

	"kope.io/auth/api/v1alpha1"
	"kope.io/auth/pkg/session"
	"kope.io/auth/pkg/session/pb"
)

type GenericOIDCProvider struct {
	config       *v1alpha1.AuthProvider
	OAuth2Config oauth2.Config
}

func NewGeneric(conf oauth2.Config, config *v1alpha1.AuthProvider) (*GenericOIDCProvider, error) {
	return &GenericOIDCProvider{
		OAuth2Config: conf,
		config:       config,
	}, nil
}

func (p *GenericOIDCProvider) Config() *v1alpha1.AuthProvider {
	return p.config
}

func (p *GenericOIDCProvider) GetLoginURL(ctx context.Context, redirectURI, state string) string {

	conf := p.OAuth2Config
	conf.RedirectURL = redirectURI

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL(state) //  oauth2.AccessTypeOffline ?
	return url
}

func (p *GenericOIDCProvider) Redeem(ctx context.Context, redirectURI, code string) (*session.Session, *oauth2.Token, error) {
	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	token, err := p.OAuth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to redeem token: %w", err)
	}

	s := &session.Session{
		SessionData: pb.SessionData{
			AccessToken:  token.AccessToken,
			ExpiresOn:    token.Expiry.Unix(),
			RefreshToken: token.RefreshToken,
			ProviderId:   p.config.Name,
			Timestamp:    time.Now().Unix(),
		},
	}

	return s, token, nil
}

// func jwtDecodeSegment(seg string) ([]byte, error) {
// 	if l := len(seg) % 4; l > 0 {
// 		seg += strings.Repeat("=", 4-l)
// 	}

// 	return base64.URLEncoding.DecodeString(seg)
// }

// func emailFromIdToken(idToken string) (string, error) {
// 	// id_token is a base64 encode ID token payload
// 	// https://developers.google.com/accounts/docs/OAuth2Login#obtainuserinfo
// 	jwtString := strings.Split(idToken, ".")
// 	b, err := jwtDecodeSegment(jwtString[1])
// 	if err != nil {
// 		return "", err
// 	}

// 	// See  Obtain user information from the ID token in https://developers.google.com/identity/protocols/OpenIDConnect

// 	//glog.Infof("JWT: %q", string(b))

// 	var jwt struct {
// 		// Always provided

// 		Subject  string `json:"sub"` // 	An identifier for the user, unique among all Google accounts and never reused. A Google account can have multiple emails at different points in time, but the sub value is never changed. Use sub within your application as the unique-identifier key for the user.
// 		Issuer   string `json:"iss"` // The Issuer Identifier for the Issuer of the response. Always https://accounts.google.com or accounts.google.com for Google ID tokens.
// 		IssuedAt int64  `json:"iat"` // The time the ID token was issued, represented in Unix time (integer seconds).
// 		Expiry   int64  `json:"exp"` // The time the ID token expires, represented in Unix time (integer seconds).

// 		// Usually provided
// 		Email         string `json:"email"`          // The user's email address. This may not be unique and is not suitable for use as a primary key. Provided only if your scope included the string "email".
// 		EmailVerified bool   `json:"email_verified"` // True if the user's e-mail address has been verified; otherwise false.

// 		// Provided for Google Apps
// 		HostedDomain string `json:"hd"` // The hosted G Suite domain of the user. Provided only if the user belongs to a hosted domain.

// 		// Sometimes provided
// 		Name    string `json:"name"`    // 	The user's full name, in a displayable form.
// 		Picture string `json:"picture"` // The URL of the user's profile picture.
// 	}
// 	err = json.Unmarshal(b, &jwt)
// 	if err != nil {
// 		return "", fmt.Errorf("error parsing JWT: %w", err)
// 	}
// 	if jwt.Email == "" {
// 		return "", errors.New("missing email")
// 	}
// 	if !jwt.EmailVerified {
// 		return "", fmt.Errorf("email %s not listed as verified", jwt.Email)
// 	}
// 	return jwt.Email, nil
// }
