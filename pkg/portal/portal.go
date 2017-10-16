package portal

import (
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

func (s *HTTPServer) portalIndex(rw http.ResponseWriter, req *http.Request) {
	index, err := ioutil.ReadFile(filepath.Join(s.staticDir, "index.html"))
	if err != nil {
		s.internalError(rw, req, err)
		return
	}

	html := string(index)

	{
		settings := make(map[string]interface{})

		settings["kubernetesUrl"] = ""
		authConfig, err := s.config.AuthConfiguration()
		if err != nil {
			glog.Warningf("error reading auth configuration: %v", err)
			authConfig = nil
		}
		if authConfig != nil {
			settings["kubernetesUrl"] = authConfig.GenerateKubeconfig.Server
		}

		settingsJson, err := json.Marshal(settings)
		if err != nil {
			s.internalError(rw, req, err)
			return
		}

		placeholder := "window.AppSettings={}"
		if !strings.Contains(html, placeholder) {
			glog.Warningf("index.html does not contain %q; will not be able to insert settings", placeholder)
		}
		html = strings.Replace(html, placeholder, "window.AppSettings="+string(settingsJson), -1)
	}

	{
		auth, err := s.authenticate(rw, req)
		if err != nil {
			s.internalError(rw, req, err)
			return
		}

		userJson := []byte("null")
		if auth != nil {
			user := &UserInfo{
				ID:       auth.Name,
				Username: auth.Spec.Username,
			}
			userJson, err = json.Marshal(user)
			if err != nil {
				s.internalError(rw, req, err)
				return
			}
		}

		placeholder := "window.User=null"
		if !strings.Contains(html, placeholder) {
			glog.Warningf("index.html does not contain %q; will not be able to insert user info", placeholder)
		}
		html = strings.Replace(html, placeholder, "window.User="+string(userJson), -1)
	}

	//content-encoding:gzip
	//content-type:text/html; charset=utf-8
	//date:Mon, 05 Jun 2017 00:07:28 GMT
	//server:nginx/1.11.12
	//status:200
	//strict-transport-security:max-age=15724800; includeSubDomains;

	rw.Write([]byte(html))
}
