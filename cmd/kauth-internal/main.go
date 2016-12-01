package main

import (
	"github.com/kopeio/kauth/pkg/k8sauth"
	"net/http"
	"fmt"
	"flag"
	"os"
	"github.com/kopeio/kauth/pkg/tokenstore"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	var o Options
	o.Listen = ":8080"
	o.Namespace = os.Getenv("NAMESPACE")

	flag.Set("logtostderr", "true")

	flag.StringVar(&o.Listen, "listen", o.Listen, "host/port on which to listen")
	flag.StringVar(&o.Namespace, "namespace", o.Namespace, "kubernetes namespace in which to store secrets")

	flag.Parse()

	err := run(&o)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unexpected error: %v\n", err)
		os.Exit(1)
	}
}

type Options struct {
	Listen    string
	Namespace string
}

func run(o *Options) error {
	if o.Namespace == "" {
		return fmt.Errorf("Namespace must be specified (either through the NAMESPACE env var or -namespace flag")
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("error building kubernetes configuration: %v", err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error building kubernetes client: %v", err)
	}

	tokenStore := tokenstore.NewSecrets(clientset, o.Namespace)

	w := &k8sauth.Webhook{
		Tokenstore: tokenStore,
	}

	mux := http.NewServeMux()

	// TODO: healthz
	//healthz.InstallHandler(mux, lbc.nginx)

	//http.HandleFunc("/build", func(w http.ResponseWriter, r *http.Request) {
	//	w.WriteHeader(http.StatusOK)
	//	fmt.Fprint(w, "build: %v - %v", gitRepo, version)
	//})
	//
	//http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
	//	c.Stop()
	//})
	//
	//if *profiling {
	//	mux.HandleFunc("/debug/pprof/", pprof.Index)
	//	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	//	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	//}


	mux.Handle("/hooks/authn", w)

	server := &http.Server{
		Addr:   o.Listen,
		Handler: mux,
	}
	return server.ListenAndServe()
}