package factory

import (
	"kope.io/auth/api/v1alpha1"
	"kope.io/auth/pkg/oauth/providers"
	"kope.io/auth/pkg/oauth/providers/github"
)

func New(config *v1alpha1.AuthProvider) (providers.Provider, error) {
	return github.New(config)
}
