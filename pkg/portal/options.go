package portal

type Options struct {
	Listen string `json:"listen"`

	EmailDomains []string `json:"emailDomains"`

	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`

	CookieSecret string `json:"cookieSecret"`

	StaticDir string `json:"staticDir"`
}
