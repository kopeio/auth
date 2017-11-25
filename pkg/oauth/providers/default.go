package providers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/golang/glog"
	"kope.io/auth/pkg/apis/componentconfig/v1alpha1"
	"kope.io/auth/pkg/oauth/session"
)

type DefaultProvider struct {
	KubernetesConfig *v1alpha1.AuthProvider
	ProviderConfiguration
}

func (p *DefaultProvider) Config() *v1alpha1.AuthProvider {
	return p.KubernetesConfig
}

// GetLoginURL with typical oauth parameters
func (p *DefaultProvider) GetLoginURL(redirectURI, state string) string {
	var a url.URL
	a = *p.LoginURL
	params, _ := url.ParseQuery(a.RawQuery)
	params.Set("redirect_uri", redirectURI)
	if p.ApprovalPrompt != "" {
		params.Set("approval_prompt", p.ApprovalPrompt)
	}
	params.Add("scope", p.Scope)
	params.Set("client_id", p.ClientID)
	params.Set("response_type", "code")
	if state != "" {
		params.Add("state", state)
	}
	a.RawQuery = params.Encode()
	return a.String()
}

// RevalidateSession performs occasional session revalidation, without a renewal
func (p *DefaultProvider) RevalidateSession(s *session.Session) (bool, error) {
	if s.AccessToken == "" {
		return false, fmt.Errorf("no access token")
	}
	if p.ValidateURL == nil {
		return false, fmt.Errorf("no ValidateURL")
	}
	endpoint := p.ValidateURL.String()

	//if len(header) == 0 {
	params := url.Values{"access_token": {s.AccessToken}}
	endpoint = endpoint + "?" + params.Encode()
	//}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return false, fmt.Errorf("error building token validation request: %v", err)
	}
	//req.Header = header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("error from token validation request: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading token validation response: %v", err)
	}
	glog.Infof("token validation response %s: %s", resp.Status, string(body))

	if resp.StatusCode == 200 {
		return true, nil
	}
	glog.Infof("token was rejected on validation: status %s - %s", resp.Status, string(body))
	return false, nil
}
