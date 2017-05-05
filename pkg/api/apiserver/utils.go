package apiserver

import (
	"bytes"
	"fmt"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/runtime"
)

func ToYAML(gvk runtime.GroupVersioner, obj runtime.Object) (string, error) {
	yaml, ok := runtime.SerializerInfoForMediaType(Codecs.SupportedMediaTypes(), "application/yaml")
	if !ok {
		return "", fmt.Errorf("YAML Serializer not registered")
	}
	//k := obj.GetObjectKind()
	//gvk := k.GroupVersionKind()

	encoder := Codecs.EncoderForVersion(yaml.Serializer, gvk)

	var w bytes.Buffer
	if err := encoder.Encode(obj, &w); err != nil {
		return "", fmt.Errorf("error encoding object: %v", err)
	}

	return w.String(), nil
}

func MustToYAML(gvk runtime.GroupVersioner, obj runtime.Object) string {
	s, err := ToYAML(gvk, obj)
	if err != nil {
		glog.Fatalf("unexpected error marshalling to yaml: %v", err)
	}
	return s
}
