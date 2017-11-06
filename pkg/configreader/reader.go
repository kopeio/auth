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

	client    authclient.Interface

	authConfigurations cache.Indexer
	authProviders      cache.Indexer
}

func New(client authclient.Interface) *ManagedConfiguration {
	return &ManagedConfiguration{
		client:    client,
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

//func (c *ManagedConfiguration) Read(p string) (*runtime.Object, error) {
//	glog.Infof("Reading config file %q", p)
//
//	yamlBytes, err := ioutil.ReadFile(p)
//	if err != nil {
//		return nil, fmt.Errorf("error reading config file %q: %v", p, err)
//	}
//
//	jsonBytes, err := yaml.ToJSON(yamlBytes)
//	if err != nil {
//		return nil, fmt.Errorf("error converting YAML config file %q: %v", p, err)
//	}
//
//	if err := json.Unmarshal(jsonBytes, c.Config); err != nil {
//		return nil, fmt.Errorf("error parsing YAML config file %q: %v", p, err)
//	}
//
//	return nil, nil
//}
//
//func (c *ManagedConfiguration) ReadFromKubernetes(k8sClient kubernetes.Interface, namespace string, name string) (runtime.Object, error) {
//	glog.Infof("Querying kubernetes for config %s/%s", namespace, name)
//
//	configMap, err := k8sClient.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
//	if err != nil {
//		if apierrors.IsNotFound(err) {
//			// TODO(authprovider-q): Create?  Probably not as it would interfere with a manager (race condition)
//			glog.Infof("Config map not found %s/%s", namespace, name)
//			return nil, nil
//		}
//		return nil, fmt.Errorf("error retrieving configmap from kubernetes: %v", err)
//	}
//
//	return c.decodeConfigMap(configMap, namespace+"/"+name)
//}
//
//func (c *ManagedConfiguration) decodeConfigMap(configMap *v1.ConfigMap, name string) (runtime.Object, error) {
//	configData := configMap.Data["config"]
//	if configData == "" {
//		glog.Warningf("No config section found in config map %s", name)
//		return nil, nil
//	}
//
//	obj, _, err := c.Decoder.Decode([]byte(configData), nil, nil)
//	if err != nil {
//		return nil, fmt.Errorf("error decoding config map %s: %v", name, err)
//	}
//
//	// TODO(authprovider-q): So how do changes work?  Do we dynamically load, or rely on k8s reloading?
//	// Presumably we rely on k8s reloading for control, but then e.g. configmaps aren't versioned...
//
//	return obj, nil
//}
