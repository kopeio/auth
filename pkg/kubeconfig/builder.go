package kubeconfig

import (
	"kope.io/auth/pkg/apis/auth/v1alpha1"
	"fmt"
	"github.com/ghodss/yaml"
)

func BuildKubeconfig(apiEndpoint string, caCertificate []byte, user *v1alpha1.User, token *v1alpha1.TokenSpec) ([]byte, error) {
	name := user.Spec.Username
	if name == "" {
		name = user.Name
	}

	cluster := KubectlCluster{
		Server:                   apiEndpoint,
		CertificateAuthorityData: caCertificate,
	}
	context := KubectlContext{
		Cluster: name,
		User:    name,
	}
	kubectlUser := KubectlUser{
		Token: EncodeToken(user, token),
	}
	config := &KubeConfig{
		ApiVersion:     "v1",
		Kind:           "Config",
		CurrentContext: name,
		Clusters: []*KubectlClusterWithName{
			{
				Name:    name,
				Cluster: cluster,
			},
		},
		Contexts: []*KubectlContextWithName{
			{
				Name:    name,
				Context: context,
			},
		},
		Users: []*KubectlUserWithName{
			{
				Name: name,
				User: kubectlUser,
			},
		},
	}

	response, err := yaml.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("error serializing kubeconfig to yaml: %v", err)
	}

	return response, nil
}
