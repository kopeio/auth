package main

import (
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"kope.io/auth/pkg/api/apiserver"
	authclient "kope.io/auth/pkg/client/clientset_generated/clientset"
	"kope.io/auth/pkg/k8sauth"
	"kope.io/auth/pkg/tokenstore"
	"net/http"
	"os"
)

func main() {
	var o Options
	o.Listen = ":8080"

	flag.Set("logtostderr", "true")

	flag.StringVar(&o.Listen, "listen", o.Listen, "host/port on which to listen")

	flag.Parse()

	err := run(&o)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unexpected error: %v\n", err)
		os.Exit(1)
	}
}

type Options struct {
	Listen string
}

func run(o *Options) error {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("error building kubernetes configuration: %v", err)
	}

	authClient, err := authclient.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error building auth client: %v", err)
	}

	// creates the clientset
	//k8sClient, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	return fmt.Errorf("error building kubernetes client: %v", err)
	//}
	//if err := authclient.RegisterResource(k8sClient); err != nil {
	//	return fmt.Errorf("error registering third party resource: %v", err)
	//}

	{
		o := apiserver.NewAuthServerOptions(os.Stdout, os.Stderr)
		if err := o.Complete(); err != nil {
			return err
		}
		if err := o.Validate(nil); err != nil {
			return err
		}
		if err := o.RunAuthServer(wait.NeverStop); err != nil {
			return err
		}
	}

	tokenStore := tokenstore.NewAPITokenStore(authClient)

	stopCh := make(chan struct{})
	go tokenStore.Run(stopCh)

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
		Addr:    o.Listen,
		Handler: mux,
	}
	return server.ListenAndServe()
}
