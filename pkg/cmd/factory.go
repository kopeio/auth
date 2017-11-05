package cmd

import "kope.io/auth/pkg/client/clientset_generated/clientset"

// Factory provides what is effectively injection for the commands
type Factory interface {
	// Clientset returns the interface to the API clients
	Clientset() (clientset.Interface, error)
}
