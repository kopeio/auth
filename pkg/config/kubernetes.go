package config

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	"kope.io/auth/api/v1alpha1"
)

type KubernetesConfigStore struct {
	namespace string
	cache     cache.Cache
}

var _ Provider = &KubernetesConfigStore{}

func (k *KubernetesConfigStore) AuthProvider(ctx context.Context, key string) (*v1alpha1.AuthProvider, error) {
	var authProvider v1alpha1.AuthProvider
	if err := k.cache.Get(ctx, types.NamespacedName{Namespace: k.namespace, Name: key}, &authProvider); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return &authProvider, nil
	// return &AuthProviderConfig{
	// 	ResourceVersion: authProvider.ResourceVersion,
	// 	ProviderID:      authProvider.Name,
	// 	PermitEmails:    authProvider.Spec.PermitEmails,
	// }, nil
}

// func (c *KubernetesConfigStore) AuthConfiguration() (*v1alpha1.AuthConfiguration, error) {
// 	attempt := 0
// 	for {
// 		keys := c.authConfigurations.ListKeys()
// 		if len(keys) == 0 {
// 			return &v1alpha1.AuthConfiguration{}, nil
// 		}

// 		if len(keys) > 1 {
// 			klog.Warningf("Multiple AuthConfiguration found %s", keys)
// 		}

// 		obj, found, err := c.authConfigurations.GetByKey(keys[0])
// 		if err != nil {
// 			return nil, fmt.Errorf("error retrieving authconfiguration: %v", err)
// 		}
// 		if found {
// 			return obj.(*v1alpha1.AuthConfiguration), nil
// 		}
// 		attempt++
// 		if attempt > 10 {
// 			return nil, fmt.Errorf("caught in mutation loop reading AuthConfiguration")
// 		}
// 	}
// }

func NewKubernetesConfigStore(ctx context.Context, restConfig *rest.Config, namespace string) (*KubernetesConfigStore, error) {
	scheme := runtime.NewScheme()
	if err := v1alpha1.AddToScheme(scheme); err != nil {
		return nil, err
	}

	options := cache.Options{Namespace: namespace, Scheme: scheme}
	cacheBuilder := cache.BuilderWithOptions(options)

	cacher, err := cacheBuilder(restConfig, options)
	if err != nil {
		return nil, fmt.Errorf("error building cache: %w", err)
	}
	if _, err := cacher.GetInformer(ctx, &v1alpha1.AuthProvider{}); err != nil {
		return nil, fmt.Errorf("error getting informer for AuthProvider: %w", err)
	}

	go func() {
		if err := cacher.Start(ctx); err != nil {
			klog.Fatalf("cacher failed unexpectedly: %w", err)
		}
	}()

	if !cacher.WaitForCacheSync(ctx) {
		return nil, fmt.Errorf("unable to sync caches")
	}

	return &KubernetesConfigStore{
		namespace: namespace,
		cache:     cacher,
	}, nil
}
