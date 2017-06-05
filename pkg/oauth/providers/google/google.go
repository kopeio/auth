package google

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang/glog"
	"kope.io/auth/pkg/apis/componentconfig/v1alpha1"
	"kope.io/auth/pkg/oauth/pb"
	"kope.io/auth/pkg/oauth/providers"
	"kope.io/auth/pkg/oauth/session"
)

type GoogleProvider struct {
	providers.DefaultProvider
	//*ProviderData
	//RedeemRefreshURL *url.URL
	//// GroupValidator is a function that determines if the passed email is in
	//// the configured Google group.
	//GroupValidator func(string) bool
}

var _ providers.Provider = &GoogleProvider{}

func NewGoogleProvider(config *v1alpha1.AuthProvider) *GoogleProvider {
	c := providers.ProviderConfiguration{}

	c.ProviderName = config.Description
	if c.ProviderName == "" {
		c.ProviderName = "Google"
	}
	c.ClientID = config.OAuthConfig.ClientID
	c.ClientSecret = config.OAuthConfig.ClientSecret

	if c.LoginURL == nil {
		c.LoginURL = &url.URL{Scheme: "https",
			Host: "accounts.google.com",
			Path: "/o/oauth2/v2/auth",
			// to get a refresh token. see https://developers.google.com/identity/protocols/OAuth2WebServer#offline
			RawQuery: "access_type=offline",
		}
	}
	if c.RedeemURL == nil {
		c.RedeemURL = &url.URL{Scheme: "https",
			Host: "www.googleapis.com",
			Path: "/oauth2/v4/token"}
	}
	if c.ValidateURL == nil {
		c.ValidateURL = &url.URL{Scheme: "https",
			Host: "www.googleapis.com",
			Path: "/oauth2/v1/tokeninfo"}
	}
	if c.Scope == "" {
		c.Scope = "profile email"
	}

	g := &GoogleProvider{
		providers.DefaultProvider{
			KubernetesConfig:      config,
			ProviderConfiguration: c,
		},
		//// Set a default GroupValidator to just always return valid (true), it will
		//// be overwritten if we configured a Google group restriction.
		//GroupValidator: func(email string) bool {
		//	return true
		//},
	}

	return g
}

func emailFromIdToken(idToken string) (string, error) {
	// id_token is a base64 encode ID token payload
	// https://developers.google.com/accounts/docs/OAuth2Login#obtainuserinfo
	jwtString := strings.Split(idToken, ".")
	b, err := jwtDecodeSegment(jwtString[1])
	if err != nil {
		return "", err
	}

	// See  Obtain user information from the ID token in https://developers.google.com/identity/protocols/OpenIDConnect

	//glog.Infof("JWT: %q", string(b))

	var jwt struct {
		// Always provided

		Subject  string `json:"sub"` // 	An identifier for the user, unique among all Google accounts and never reused. A Google account can have multiple emails at different points in time, but the sub value is never changed. Use sub within your application as the unique-identifier key for the user.
		Issuer   string `json:"iss"` // The Issuer Identifier for the Issuer of the response. Always https://accounts.google.com or accounts.google.com for Google ID tokens.
		IssuedAt int64  `json:"iat"` // The time the ID token was issued, represented in Unix time (integer seconds).
		Expiry   int64  `json:"exp"` // The time the ID token expires, represented in Unix time (integer seconds).

		// Usually provided
		Email         string `json:"email"`          // The user's email address. This may not be unique and is not suitable for use as a primary key. Provided only if your scope included the string "email".
		EmailVerified bool   `json:"email_verified"` // True if the user's e-mail address has been verified; otherwise false.

		// Provided for Google Apps
		HostedDomain string `json:"hd"` // The hosted G Suite domain of the user. Provided only if the user belongs to a hosted domain.

		// Sometimes provided
		Name    string `json:"name"`    // 	The user's full name, in a displayable form.
		Picture string `json:"picture"` // The URL of the user's profile picture.
	}
	err = json.Unmarshal(b, &jwt)
	if err != nil {
		return "", fmt.Errorf("error parsing JWT: %v", err)
	}
	if jwt.Email == "" {
		return "", errors.New("missing email")
	}
	if !jwt.EmailVerified {
		return "", fmt.Errorf("email %s not listed as verified", jwt.Email)
	}
	return jwt.Email, nil
}

func jwtDecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}

func (p *GoogleProvider) Redeem(redirectURL, code string) (*session.Session, *session.UserInfo, error) {
	if code == "" {
		return nil, nil, errors.New("missing code")
	}

	params := url.Values{}
	params.Add("redirect_uri", redirectURL)
	params.Add("client_id", p.ClientID)
	params.Add("client_secret", p.ClientSecret)
	params.Add("code", code)
	params.Add("grant_type", "authorization_code")
	var req *http.Request
	req, err := http.NewRequest("POST", p.RedeemURL.String(), bytes.NewBufferString(params.Encode()))
	if err != nil {
		return nil, nil, fmt.Errorf("error building redemption HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error making redemption HTTP request: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading HTTP response: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("unexpected response %q from redemption HTTP request: %v", resp.Status, err)
	}

	glog.V(2).Infof("Redeem response %s", string(body))

	var jsonResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		IdToken      string `json:"id_token"`
		TokenType    string `json:"token_type"`
	}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing redemption HTTP response: %v", err)
	}
	email, err := emailFromIdToken(jsonResponse.IdToken)
	if err != nil {
		return nil, nil, err
	}
	s := &session.Session{
		SessionData: pb.SessionData{
			AccessToken:  jsonResponse.AccessToken,
			ExpiresOn:    time.Now().Unix() + jsonResponse.ExpiresIn,
			RefreshToken: jsonResponse.RefreshToken,
			ProviderId:   p.Config().Name,
			Timestamp:    time.Now().Unix(),
		},
	}
	extra := &session.UserInfo{
		Email:          email,
		ProviderID:     p.KubernetesConfig.Name,
		ProviderUserID: email,
	}
	return s, extra, nil
}

//// SetGroupRestriction configures the GoogleProvider to restrict access to the
//// specified group(s). AdminEmail has to be an administrative email on the domain that is
//// checked. CredentialsFile is the path to a json file containing a Google service
//// account credentials.
//func (p *GoogleProvider) SetGroupRestriction(groups []string, adminEmail string, credentialsReader io.Reader) {
//	adminService := getAdminService(adminEmail, credentialsReader)
//	p.GroupValidator = func(email string) bool {
//		return userInGroup(adminService, groups, email)
//	}
//}
//
//func getAdminService(adminEmail string, credentialsReader io.Reader) *admin.Service {
//	data, err := ioutil.ReadAll(credentialsReader)
//	if err != nil {
//		log.Fatal("can't read Google credentials file:", err)
//	}
//	conf, err := google.JWTConfigFromJSON(data, admin.AdminDirectoryUserReadonlyScope, admin.AdminDirectoryGroupReadonlyScope)
//	if err != nil {
//		log.Fatal("can't load Google credentials file:", err)
//	}
//	conf.Subject = adminEmail
//
//	client := conf.Client(oauth2.NoContext)
//	adminService, err := admin.New(client)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return adminService
//}
//
//func userInGroup(service *admin.Service, groups []string, email string) bool {
//	user, err := fetchUser(service, email)
//	if err != nil {
//		log.Printf("error fetching user: %v", err)
//		return false
//	}
//	id := user.Id
//	custID := user.CustomerId
//
//	for _, group := range groups {
//		members, err := fetchGroupMembers(service, group)
//		if err != nil {
//			log.Printf("error fetching group members: %v", err)
//			return false
//		}
//
//		for _, member := range members {
//			switch member.Type {
//			case "CUSTOMER":
//				if member.Id == custID {
//					return true
//				}
//			case "USER":
//				if member.Id == id {
//					return true
//				}
//			}
//		}
//	}
//	return false
//}
//
//func fetchUser(service *admin.Service, email string) (*admin.User, error) {
//	user, err := service.Users.Get(email).Do()
//	return user, err
//}
//
//func fetchGroupMembers(service *admin.Service, group string) ([]*admin.Member, error) {
//	members := []*admin.Member{}
//	pageToken := ""
//	for {
//		req := service.Members.List(group)
//		if pageToken != "" {
//			req.PageToken(pageToken)
//		}
//		r, err := req.Do()
//		if err != nil {
//			return nil, err
//		}
//		for _, member := range r.Members {
//			members = append(members, member)
//		}
//		if r.NextPageToken == "" {
//			break
//		}
//		pageToken = r.NextPageToken
//	}
//	return members, nil
//}
//
//// ValidateGroup validates that the provided email exists in the configured Google
//// group(s).
//func (p *GoogleProvider) ValidateGroup(email string) bool {
//	return p.GroupValidator(email)
//}

func (p *GoogleProvider) RefreshSessionIfNeeded(s *session.Session) (bool, error) {
	if s == nil || s.ExpiresOn > time.Now().Unix() || s.RefreshToken == "" {
		return false, nil
	}

	newToken, duration, err := p.redeemRefreshToken(s.RefreshToken)
	if err != nil {
		return false, err
	}

	//// re-check that the user is in the proper google group(s)
	//if !p.ValidateGroup(s.Email) {
	//	return false, fmt.Errorf("%s is no longer in the group(s)", s.Email)
	//}

	origExpiration := s.ExpiresOn
	s.AccessToken = newToken
	s.ExpiresOn = time.Now().Unix() + int64(duration.Seconds())
	log.Printf("refreshed access token %s (expired on %s)", s, origExpiration)
	return true, nil
}

func (p *GoogleProvider) redeemRefreshToken(refreshToken string) (string, time.Duration, error) {
	// https://developers.google.com/identity/protocols/OAuth2WebServer#refresh
	params := url.Values{}
	params.Add("client_id", p.ClientID)
	params.Add("client_secret", p.ClientSecret)
	params.Add("refresh_token", refreshToken)
	params.Add("grant_type", "refresh_token")
	var req *http.Request
	req, err := http.NewRequest("POST", p.RedeemURL.String(), bytes.NewBufferString(params.Encode()))
	if err != nil {
		return "", 0, fmt.Errorf("error building refreshing request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("error doing refreshing request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("error reading refreshing response: %v", err)
	}

	if resp.StatusCode != 200 {
		return "", 0, fmt.Errorf("got %d from refresh request", resp.Status)
	}

	var data struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}

	// TODO: Anything else of value?
	glog.Infof("TEMPORARY refresh response: %s", string(body))

	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", 0, fmt.Errorf("error parsing refresh response: %v", err)
	}

	expires := time.Duration(data.ExpiresIn) * time.Second
	return data.AccessToken, expires, nil
}
