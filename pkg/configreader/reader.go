package configreader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/golang/glog"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

type ManagedConfiguration struct {
	Config  runtime.Object
	Decoder runtime.Decoder
}

func (c *ManagedConfiguration) Read(p string) (*runtime.Object, error) {
	glog.Infof("Reading config file %q", p)

	yamlBytes, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %q: %v", p, err)
	}

	jsonBytes, err := yaml.ToJSON(yamlBytes)
	if err != nil {
		return nil, fmt.Errorf("error converting YAML config file %q: %v", p, err)
	}

	if err := json.Unmarshal(jsonBytes, c.Config); err != nil {
		return nil, fmt.Errorf("error parsing YAML config file %q: %v", p, err)
	}

	return nil, nil
}

func (c *ManagedConfiguration) ReadFromKubernetes(k8sClient kubernetes.Interface, namespace string, name string) (runtime.Object, error) {
	glog.Infof("Querying kubernetes for config %s/%s", namespace, name)

	configMap, err := k8sClient.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			// TODO(authprovider-q): Create?  Probably not as it would interfere with a manager (race condition)
			glog.Infof("Config map not found %s/%s", namespace, name)
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving configmap from kubernetes: %v", err)
	}

	return c.decodeConfigMap(configMap, namespace+"/"+name)
}

func (c *ManagedConfiguration) decodeConfigMap(configMap *v1.ConfigMap, name string) (runtime.Object, error) {
	configData := configMap.Data["config"]
	if configData == "" {
		glog.Warningf("No config section found in config map %s", name)
		return nil, nil
	}

	obj, _, err := c.Decoder.Decode([]byte(configData), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error decoding config map %s: %v", name, err)
	}

	// TODO(authprovider-q): So how do changes work?  Do we dynamically load, or rely on k8s reloading?
	// Presumably we rely on k8s reloading for control, but then e.g. configmaps aren't versioned...

	return obj, nil
}
