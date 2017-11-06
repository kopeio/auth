package cmd

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
	"kope.io/auth/pkg/client/clientset_generated/clientset"
	"k8s.io/client-go/rest"
)

// DefaultFactory providers the default implementation of Factory
type DefaultFactory struct {
	clientset clientset.Interface

	options *FactoryOptions
}

var _ Factory = &DefaultFactory{}

type FactoryOptions struct {
	Kubeconfig string
}

// Clientset implements Factory::Clientset
func (f *DefaultFactory) Clientset() (clientset.Interface, error) {
	if f.clientset == nil {
		config, err := f.Config()
		if err != nil {
			return nil, err
		}

		client, err := clientset.NewForConfig(config)
		if err != nil {
			return nil, fmt.Errorf("error building client: %v", err)
		}

		f.clientset = client
	}

	return f.clientset, nil
}

func (f*DefaultFactory) Config() (*rest.Config, error) {
	kubeconfig := f.options.Kubeconfig
	if kubeconfig == "" {
		return nil, fmt.Errorf("kubeconfig path must be provided")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("error reading kubeconfig %q: %v", kubeconfig, err)
	}

	return config, nil
}

func NewDefaultFactory(options *FactoryOptions) Factory {
	f := &DefaultFactory{
		options: options,
	}
	return f
}
