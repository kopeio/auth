package cmd

import (
	"kope.io/auth/pkg/client/clientset_generated/clientset"
	"k8s.io/client-go/rest"
)

// Factory provides what is effectively injection for the commands
type Factory interface {
	// Clientset returns the interface to the API clients
	Clientset() (clientset.Interface, error)

	Config() (*rest.Config, error)
}
