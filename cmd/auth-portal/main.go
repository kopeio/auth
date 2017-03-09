package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/yaml"
	"kope.io/auth/pkg/portal"
	"os"
)

func main() {
	var o portal.Options
	o.Listen = ":8080"

	o.StaticDir = "/webapp"

	o.ClientID = os.Getenv("OAUTH2_CLIENT_ID")
	o.ClientSecret = os.Getenv("OAUTH2_CLIENT_SECRET")
	o.CookieSecret = os.Getenv("OAUTH2_COOKIE_SECRET")

	flag.Set("logtostderr", "true")

	configFile := os.Getenv("CONFIG")
	if configFile != "" {
		fmt.Fprintf(os.Stderr, "Reading config file %q", configFile)

		yamlBytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading config file %q: %v\n", configFile, err)
			os.Exit(1)
		}

		jsonBytes, err := yaml.ToJSON(yamlBytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error converting YAML config file %q: %v\n", configFile, err)
			os.Exit(1)
		}

		if err := json.Unmarshal(jsonBytes, &o); err != nil {
			fmt.Fprintf(os.Stderr, "error parsing YAML config file %q: %v\n", configFile, err)
			os.Exit(1)
		}
	}
	flag.StringVar(&o.Listen, "listen", o.Listen, "host/port on which to listen")

	flag.Parse()

	err := run(&o)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unexpected error: %v\n", err)
		os.Exit(1)
	}
}

func run(o *portal.Options) error {
	p, err := portal.NewHTTPServer(o)
	if err != nil {
		return err
	}

	return p.ListenAndServe()
}