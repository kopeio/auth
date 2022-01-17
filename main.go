package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"os"

	cryptorand "crypto/rand"
	mathrand "math/rand"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	"kope.io/auth/pkg/config"
	"kope.io/auth/pkg/httpserver"
	"kope.io/auth/pkg/keystore"
)

func main() {
	err := run(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cryptoSeed()

	listen := ":8080"
	flag.StringVar(&listen, "listen", listen, "endpoint on which to listen")

	klog.InitFlags(nil)
	flag.Parse()

	restConfig, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("error building kubernetes client configuration: %w", err)
	}

	k8sClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return fmt.Errorf("error building kubernetes client: %w", err)
	}

	namespaceBytes, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return fmt.Errorf("error reading namespace from %q: %w", "/var/run/secrets/kubernetes.io/serviceaccount/namespace", err)
	}
	namespace := string(namespaceBytes)

	keystore, err := keystore.NewKubernetesKeyStore(k8sClient, namespace, "auth")
	if err != nil {
		return err
	}
	go keystore.WatchForever(ctx)

	config, err := config.NewKubernetesConfigStore(ctx, restConfig, namespace)
	if err != nil {
		return err
	}

	httpServer, err := httpserver.NewHTTPServer(ctx, config, keystore)
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}
	return httpServer.ListenAndServe(listen)
}

func cryptoSeed() {
	data := make([]byte, 8)
	_, err := cryptorand.Read(data)
	if err != nil {
		klog.Fatalf("error seeding random numbers: %v", err)
	}
	seed := binary.BigEndian.Uint64(data)
	mathrand.Seed(int64(seed))
}
