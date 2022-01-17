package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
	"kope.io/auth/api/v1alpha1"
	"kope.io/auth/pkg/oauth/providers"
	"kope.io/auth/pkg/session"
)

type GithubAuthProvider struct {
	*providers.GenericOIDCProvider
}

func New(config *v1alpha1.AuthProvider) (*GithubAuthProvider, error) {

	conf := oauth2.Config{
		ClientID:     config.Spec.ClientID,
		ClientSecret: config.Spec.ClientSecret,

		Scopes: []string{"user:email"}, //", "SCOPE2"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}

	generic, err := providers.NewGeneric(conf, config)
	if err != nil {
		return nil, err
	}
	return &GithubAuthProvider{
		GenericOIDCProvider: generic,
	}, nil
}

func (p *GithubAuthProvider) Redeem(ctx context.Context, redirectURI, code string) (*session.Session, *session.UserInfo, error) {
	sessionInfo, token, err := p.GenericOIDCProvider.Redeem(ctx, redirectURI, code)
	if err != nil {
		return nil, nil, err
	}

	httpClient := p.OAuth2Config.Client(ctx, token)

	userInfo := &session.UserInfo{
		ProviderID: sessionInfo.ProviderId,
	}

	if err := p.populateUserInfo(ctx, httpClient, userInfo); err != nil {
		return nil, nil, err
	}
	return sessionInfo, userInfo, nil
}

func doJSONGet(ctx context.Context, httpClient *http.Client, url string, dest interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	req = req.WithContext(ctx)

	response, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %s", response.Status)
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	if err := json.Unmarshal(b, dest); err != nil {
		return fmt.Errorf("error parsing response as JSON: %w", err)
	}

	return nil
}

func (p *GithubAuthProvider) populateUserInfo(ctx context.Context, httpClient *http.Client, userInfo *session.UserInfo) error {
	// https://docs.github.com/en/rest/reference/users#get-the-authenticated-user

	var response getUserResponse

	if err := doJSONGet(ctx, httpClient, "https://api.github.com/user", &response); err != nil {
		return fmt.Errorf("failed to get user from github API: %w", err)
	}

	if response.Email == "" {
		// TODO: email may be empty: https://stackoverflow.com/questions/35373995/github-user-email-is-null-despite-useremail-scope
		return fmt.Errorf("user email was not provided by github API")
	}

	userInfo.Email = response.Email
	return nil
}

type getUserResponse struct {
	Name     string `json:"name"`
	Company  string `json:"company"`
	Location string `json:"location"`
	Email    string `json:"email"`
}
