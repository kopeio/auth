package configreader

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	"kope.io/auth/pkg/api/apiserver"
	"kope.io/auth/pkg/apis/componentconfig"
	"kope.io/auth/pkg/apis/componentconfig/v1alpha1"
)

func TestReadConfigMap(t *testing.T) {
	//apiContext, err := api.NewAPIContext("")
	//if err != nil {
	//	t.Fatalf("error creating API context: %v", err)
	//}
	//
	//componentconfiginstall.Install(apiContext.GroupFactoryRegistry, apiContext.Registry, apiContext.Scheme)

	mc := &ManagedConfiguration{
		Decoder: apiserver.Codecs.UniversalDecoder(),
	}
	configMap := &v1.ConfigMap{}
	configMap.Data = make(map[string]string)
	configString := `
apiVersion: auth.kope.io/v1alpha1
kind: AuthConfiguration
metadata:
  creationTimestamp: null
spec:
  authProviders:
  - id: "123"
    name: Some provider
    oAuthConfig:
      clientID: ABCDEFG
      clientSecret: HIJKLMNOP
    permitEmails:
    - '*@google.com'
`
	configMap.Data["config"] = configString
	obj, err := mc.decodeConfigMap(configMap, "test/test")
	if err != nil {
		t.Fatalf("failed to decode config map: %v", err)
	}

	expected := &componentconfig.AuthConfiguration{
		Spec: componentconfig.AuthConfigurationSpec{
			AuthProviders: []componentconfig.AuthProviderSpec{
				{
					ID:           "123",
					Name:         "Some provider",
					PermitEmails: []string{"*@google.com"},
					OAuthConfig: &componentconfig.OAuthConfig{
						ClientID:     "ABCDEFG",
						ClientSecret: "HIJKLMNOP",
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(expected, obj) {
		t.Logf("Expected: %s", apiserver.MustToYAML(v1alpha1.SchemeGroupVersion, expected))
		t.Logf("Actual: %v", apiserver.MustToYAML(v1alpha1.SchemeGroupVersion, obj))
		t.Fatalf("Unexpected value")
	}
}
