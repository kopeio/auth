package config

import (
	"context"

	"kope.io/auth/api/v1alpha1"
)

type Provider interface {
	AuthProvider(ctx context.Context, key string) (*v1alpha1.AuthProvider, error)
}

// type AuthProviderConfig struct {
// 	ResourceVersion string

// 	ProviderID string

// 	PermitEmails []string

// 	ClientID     string
// 	ClientSecret string
// }
