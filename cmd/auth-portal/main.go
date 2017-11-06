package main

import (
	"flag"
	"fmt"
	"os"

	cryptorand "crypto/rand"
	"encoding/binary"
	"github.com/golang/glog"
	//apierrors "k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	authclient "kope.io/auth/pkg/client/clientset_generated/clientset"
	"kope.io/auth/pkg/configreader"
	"kope.io/auth/pkg/keystore"
	"kope.io/auth/pkg/portal"
	"kope.io/auth/pkg/tokenstore"
	mathrand "math/rand"
	"io/ioutil"
)

//const CookieSigningSecretLength = 24

func main() {
	cryptoSeed()

	flag.Set("logtostderr", "true")

	// TODO(authprovider-q): Some parameters we don't really want configurable, because
	// we expect to be running in a container.  But maybe they would be useful for people
	// that want to run the code differently, so they probably warrant a flag or an env var.
	// Thoughts?
	listen := ":8080"
	flag.StringVar(&listen, "listen", listen, "host/port on which to listen")
	staticDir := "/webapp"
	flag.StringVar(&staticDir, "static-dir", staticDir, "location of static directory")
	server := ""
	flag.StringVar(&server, "server", server, "override API server location")
	insecureServer := false
	flag.BoolVar(&insecureServer, "insecure-skip-tls-verify", insecureServer, "skip verification of server certificate (this is insecure)")

	flag.Parse()

	err := run(listen, staticDir, server, insecureServer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unexpected error: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run(listen string, staticDir string, server string, insecureServer bool) error {
	// creates the in-cluster config
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("error building kubernetes client configuration: %v", err)
	}

	// creates the clientset
	k8sClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("error building kubernetes client: %v", err)
	}

	authRestConfig, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("error building kubernetes client configuration: %v", err)
	}
	if server != "" {
		authRestConfig.Host = server
	}
	if insecureServer {
		authRestConfig.Insecure = insecureServer

		// Avoid "specifying a root certificates file with the insecure flag is not allowed"
		authRestConfig.TLSClientConfig.CAData = nil
		authRestConfig.TLSClientConfig.CAFile = ""
	}

	authClient, err := authclient.NewForConfig(authRestConfig)
	if err != nil {
		return fmt.Errorf("error building auth client: %v", err)
	}

	//configs, err := authClient.ComponentconfigV1alpha1().AuthConfigurations(namespace).List(metav1.ListOptions{})
	//if err != nil {
	//		return fmt.Errorf("error reading AuthConfigurations from API: %v", err)
	//}
	//
	//authProviderList, err := authClient.ComponentconfigV1alpha1().AuthProviders(namespace).List(metav1.ListOptions{})
	//if err != nil {
	//	return fmt.Errorf("error reading AuthProviders from API: %v", err)
	//}

	stopCh := make(chan struct{})

	//apiContext, err := api.NewAPIContext(os.Getenv("API_VERSIONS"))
	//if err != nil {
	//	return fmt.Errorf("error initializing API: %v", err)
	//}
	//
	//componentconfiginstall.Install(apiContext.GroupFactoryRegistry, apiContext.Registry, apiContext.Scheme)

	//configDecoder := apiserver.Codecs.UniversalDecoder()
	//
	//configReader := &configreader.ManagedConfiguration{
	//	Decoder: configDecoder,
	//}

	configReader := configreader.New(authClient)
	if err := configReader.StartWatches(stopCh); err != nil {
		return fmt.Errorf("error starting configuration watches: %v", err)
	}

	//configFile := os.Getenv("CONFIG")
	//if configFile != "" {
	//	err := configReader.Read(configFile)
	//	if err != nil {
	//		return fmt.Errorf("error reading config file %q: %v\n", configFile, err)
	//	}
	//}

	//configObj, err := configReader.ReadFromKubernetes(k8sClient, namespace, name)
	//if err != nil {
	//	return fmt.Errorf("error reading configuration: %v", err)
	//}

	//// TODO(authprovider-q): Should we deal with v1alpha1 or unversioned when we own the API?
	//// (I guess the same question with our User objects)
	//config := configObj.(*authprovider.AuthConfiguration)

	namespaceBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return fmt.Errorf("error reading namespace from %q: %v", "/var/run/secrets/kubernetes.io/serviceaccount/namespace", err)
	}
	namespace := string(namespaceBytes)

	keyStore, err := keystore.NewKubernetesKeyStore(k8sClient, namespace, "auth")
	if err != nil {
		return err
	}
	go keyStore.Run(stopCh)

	//sharedSecretSet, err := secretStore.EnsureSharedSecretSet("cookie-signing", generateCookieSigningSecrets)
	//if err != nil {
	//	return err
	//}

	//o.ClientID = os.Getenv("OAUTH2_CLIENT_ID")
	//o.ClientSecret = os.Getenv("OAUTH2_CLIENT_SECRET")
	//o.CookieSecret = os.Getenv("OAUTH2_COOKIE_SECRET")

	tokenStore := tokenstore.NewAPITokenStore(authClient)
	go tokenStore.Run(stopCh)

	p, err := portal.NewHTTPServer(configReader, listen, staticDir, keyStore, tokenStore)
	if err != nil {
		return err
	}

	return p.ListenAndServe()
}

//func generateCookieSigningSecrets() ([]byte, error) {
//	data := make([]byte, CookieSigningSecretLength)
//	_, err := cryptorand.Read(data)
//	if err != nil {
//		return nil, fmt.Errorf("error generating cookie signing secret: %v", err)
//	}
//	return data, nil
//}

func cryptoSeed() {
	data := make([]byte, 8)
	_, err := cryptorand.Read(data)
	if err != nil {
		glog.Fatalf("error seeding random numbers: %v", err)
	}
	seed := binary.BigEndian.Uint64(data)
	mathrand.Seed(int64(seed))
}
