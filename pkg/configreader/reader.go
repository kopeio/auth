package configreader

import (
	"fmt"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"kope.io/auth/pkg/apis/componentconfig/v1alpha1"
	authclient "kope.io/auth/pkg/client/clientset_generated/clientset"
)

type ManagedConfiguration struct {
	//Config  runtime.Object
	//Decoder runtime.Decoder

	client authclient.Interface

	authConfigurations cache.Indexer
	authProviders      cache.Indexer
}

func New(client authclient.Interface) *ManagedConfiguration {
	return &ManagedConfiguration{
		client: client,
	}
}

func (c *ManagedConfiguration) AuthConfiguration() (*v1alpha1.AuthConfiguration, error) {
	attempt := 0
	for {
		keys := c.authConfigurations.ListKeys()
		if len(keys) == 0 {
			return &v1alpha1.AuthConfiguration{}, nil
		}

		if len(keys) > 1 {
			glog.Warningf("Multiple AuthConfiguration found %s", keys)
		}

		obj, found, err := c.authConfigurations.GetByKey(keys[0])
		if err != nil {
			return nil, fmt.Errorf("error retrieving authconfiguration: %v", err)
		}
		if found {
			return obj.(*v1alpha1.AuthConfiguration), nil
		}
		attempt++
		if attempt > 10 {
			return nil, fmt.Errorf("caught in mutation loop reading AuthConfiguration")
		}
	}
}

func (c *ManagedConfiguration) AuthProvider(name string) (*v1alpha1.AuthProvider, error) {
	obj, found, err := c.authProviders.GetByKey(name)
	if err != nil {
		return nil, fmt.Errorf("error retrieving AuthProvider configuration: %v", err)
	}
	if found {
		return obj.(*v1alpha1.AuthProvider), nil
	}
	//glog.Infof("provider %q not found; actual %s", name, c.authProviders.ListKeys())
	return nil, nil
}

func (c *ManagedConfiguration) StartWatches(stopCh <-chan struct{}) error {
	{
		watcher := cache.NewListWatchFromClient(c.client.ConfigV1alpha1().RESTClient(), "authconfigurations", "", fields.Everything())
		indexer, informer := cache.NewIndexerInformer(watcher, &v1alpha1.AuthConfiguration{}, 0,
			cache.ResourceEventHandlerFuncs{},
			cache.Indexers{})

		go informer.Run(stopCh)

		// Wait for all involved caches to be synced, before processing items from the queue is started
		if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
			return fmt.Errorf("Timed out waiting for caches to sync")
		}

		c.authConfigurations = indexer
	}

	{
		watcher := cache.NewListWatchFromClient(c.client.ConfigV1alpha1().RESTClient(), "authproviders", "", fields.Everything())
		indexer, informer := cache.NewIndexerInformer(watcher, &v1alpha1.AuthProvider{}, 0,
			cache.ResourceEventHandlerFuncs{},
			cache.Indexers{})

		go informer.Run(stopCh)

		// Wait for all involved caches to be synced, before processing items from the queue is started
		if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
			return fmt.Errorf("Timed out waiting for caches to sync")
		}

		c.authProviders = indexer
	}

	return nil
}
