package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/pflag"
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
	//o.Server = "http://127.0.0.1:8080"
	//o.ServerInsecure = false

	pflag.Set("logtostderr", "true")
	flag.CommandLine.Parse([]string{"--logtostderr=true"})

	o.AuthServer = apiserver.NewAuthServerOptions(os.Stdout, os.Stderr)

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	//pflag.StringVar(&o.Server, "server", o.Server, "url on which to connect to server")
	pflag.StringVar(&o.Listen, "listen", o.Listen, "host/port on which to listen")
	//pflag.BoolVar(&o.ServerInsecure, "insecure-skip-tls-verify", o.ServerInsecure, "skip verification of server certificate (this is insecure)")

	o.AuthServer.AddFlags(pflag.CommandLine)

	pflag.Parse()

	// HACK: Create /tmp, so we don't need to create it in the base image
	if err := os.MkdirAll("/tmp", 0777|os.ModeTemporary); err != nil {
		glog.Warning("failed to mkdir /tmp: %v", err)
	}

	err := run(&o)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unexpected error: %v\n", err)
		os.Exit(1)
	}
}

type Options struct {
	Listen         string
	//Server         string
	//ServerInsecure bool
	AuthServer     *apiserver.AuthServerOptions
}

func run(o *Options) error {
	// creates the clientset
	//k8sClient, err := kubernetes.NewForConfig(config)
	//if err != nil {
	//	return fmt.Errorf("error building kubernetes client: %v", err)
	//}
	//if err := authclient.RegisterResource(k8sClient); err != nil {
	//	return fmt.Errorf("error registering third party resource: %v", err)
	//}

	{
		if err := o.AuthServer.Complete(); err != nil {
			return err
		}
		if err := o.AuthServer.Validate(nil); err != nil {
			return err
		}
		go func() {
			if err := o.AuthServer.RunAuthServer(wait.NeverStop); err != nil {
				glog.Fatalf("error running API server: %v", err)
			}
		}()
	}

	 //creates the in-cluster config
	authRestConfig, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("error building kubernetes configuration: %v", err)
	}

	//u, err := url.Parse(o.Server)
	//if err != nil {
	//	return fmt.Errorf("Invalid server flag: %q", o.Server)
	//}
	//
	//authRestConfig := &rest.Config{
	//	Host: u.Host,
	//}
	//
	//if o.ServerInsecure {
	//	authRestConfig.Insecure = o.ServerInsecure
	//
	//	// Avoid "specifying a root certificates file with the insecure flag is not allowed"
	//	authRestConfig.TLSClientConfig.CAData = nil
	//	authRestConfig.TLSClientConfig.CAFile = ""
	//}

	authClient, err := authclient.NewForConfig(authRestConfig)
	if err != nil {
		return fmt.Errorf("error building user client: %v", err)
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
	glog.Infof("starting hook server on %s", o.Listen)
	return server.ListenAndServe()
}
