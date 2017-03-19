package api

import (
	"bytes"
	"fmt"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/apimachinery/announced"
	"k8s.io/apimachinery/pkg/apimachinery/registered"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

type APIContext struct {
	GroupFactoryRegistry announced.APIGroupFactoryRegistry

	Registry *registered.APIRegistrationManager

	Scheme *runtime.Scheme

	Codecs serializer.CodecFactory
}

func NewAPIContext(apiVersions string) (*APIContext, error) {
	c := &APIContext{}
	c.GroupFactoryRegistry = make(announced.APIGroupFactoryRegistry)

	registry, err := registered.NewAPIRegistrationManager(apiVersions)
	if err != nil {
		return nil, fmt.Errorf("Could not construct version manager: %v (apiVersions=%q)", err, apiVersions)
	}
	c.Registry = registry

	c.Scheme = runtime.NewScheme()

	c.Codecs = serializer.NewCodecFactory(c.Scheme)

	return c, nil
}

func (a *APIContext) ToYAML(gvk runtime.GroupVersioner, obj runtime.Object) (string, error) {
	yaml, ok := runtime.SerializerInfoForMediaType(a.Codecs.SupportedMediaTypes(), "application/yaml")
	if !ok {
		return "", fmt.Errorf("YAML Serializer not registered")
	}
	//k := obj.GetObjectKind()
	//gvk := k.GroupVersionKind()

	encoder := a.Codecs.EncoderForVersion(yaml.Serializer, gvk)

	var w bytes.Buffer
	if err := encoder.Encode(obj, &w); err != nil {
		return "", fmt.Errorf("error encoding object: %v", err)
	}

	return w.String(), nil
}

func (a *APIContext) MustToYAML(gvk runtime.GroupVersioner, obj runtime.Object) string {
	s, err := a.ToYAML(gvk, obj)
	if err != nil {
		glog.Fatalf("unexpected error marshalling to yaml: %v", err)
	}
	return s
}
