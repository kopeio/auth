package portal

import (
	"encoding/json"
	"net/http"
	"strings"
)

var indexTemplate = `
<html>
<head>
    <meta charset="utf-8">
    <title>Kubernetes Auth Portal</title>
</head>
<body>
<style>
    @import url('https://fonts.googleapis.com/css?family=Roboto:300,400,500:latin');

    html {
        font-family: 'Roboto', sans-serif;
    }
</style>
<div id="app" />
<script type="text/javascript">
    window.initialProps = {{initialProps}};
</script>
<script src="/static/bundle.js" type="text/javascript"></script>
</body>
</html>
`

func (s *HTTPServer) portalIndex(rw http.ResponseWriter, req *http.Request) {
	status, err := s.status(rw, req)
	if err != nil {
		s.internalError(rw, req, err)
		return
	}

	statusJson, err := json.Marshal(status)
	if err != nil {
		s.internalError(rw, req, err)
		return
	}

	html := strings.Replace(indexTemplate, "{{initialProps}}", string(statusJson), -1)
	rw.Write([]byte(html))
}

type PortalStatus struct {
	User *UserInfo `json:"user"`
}

func (s *HTTPServer) status(rw http.ResponseWriter, req *http.Request) (*PortalStatus, error) {
	auth, err := s.authenticate(rw, req)
	if err != nil {
		return nil, err
	}

	status := &PortalStatus{}

	if auth != nil {
		status.User = &UserInfo{
			ID:       auth.Metadata.Name,
			Username: auth.Spec.Username,
		}
	}

	return status, nil
}
