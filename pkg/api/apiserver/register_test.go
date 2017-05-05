package apiserver

import (
	"testing"

	"bytes"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"kope.io/auth/pkg/apis/componentconfig/v1alpha1"

	"kope.io/auth/pkg/api/apiserver"
)

func TestDefaultSchema(t *testing.T) {
	//apiContext, err := api.NewAPIContext("")
	//if err != nil {
	//	t.Fatalf("error creating API context: %v", err)
	//}
	//
	//Install(apiContext.GroupFactoryRegistry, apiContext.Registry, apiContext.Scheme)

	yaml, ok := runtime.SerializerInfoForMediaType(apiserver.Codecs.SupportedMediaTypes(), "application/yaml")
	if !ok {
		t.Fatalf("no YAML serializer registered")
	}
	gv := v1alpha1.SchemeGroupVersion
	encoder := Codecs.EncoderForVersion(yaml.Serializer, gv)

	obj := &v1alpha1.AuthConfiguration{
		Spec: v1alpha1.AuthConfigurationSpec{
			GenerateKubeconfig: &v1alpha1.GenerateKubeconfig{
				Server: "https://api.example.com",
			},
			AuthProviders: []v1alpha1.AuthProviderSpec{
				{
					ID:           "123",
					Name:         "Some provider",
					PermitEmails: []string{"*@google.com"},
					OAuthConfig: &v1alpha1.OAuthConfig{
						ClientID:     "ABCDEFG",
						ClientSecret: "HIJKLMNOP",
					},
				},
			},
		},
	}
	var w bytes.Buffer
	if err := encoder.Encode(obj, &w); err != nil {
		t.Fatalf("error encoding object")
	}

	actual := w.String()
	expected := `
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
  generateKubeconfig:
    server: https://api.example.com
`
	if strings.TrimSpace(actual) != strings.TrimSpace(expected) {
		t.Logf("Expected: %s", expected)
		t.Logf("Actual: %s", actual)
		t.Fatalf("Unexpected value")
	}
}
