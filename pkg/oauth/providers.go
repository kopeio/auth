package oauth

import (
	"context"
	"fmt"

	"k8s.io/klog/v2"
	"kope.io/auth/pkg/oauth/providers"
	"kope.io/auth/pkg/oauth/providers/factory"
)

func (s *Server) getProvider(ctx context.Context, id string) (providers.Provider, error) {
	config, err := s.Config.AuthProvider(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting configuration for %q: %v", id, err)
	}

	if config == nil {
		klog.Warningf("provider %q not configured", id)
		return nil, nil
	}

	s.providersMutex.Lock()
	defer s.providersMutex.Unlock()

	existing := s.providers[id]
	if existing != nil && existing.Config().ResourceVersion == config.ResourceVersion {
		return existing, nil
	}

	provider, err := factory.New(config)
	if err != nil {
		return nil, fmt.Errorf("error creating provider: %w", err)
	}
	if s.providers == nil {
		s.providers = make(map[string]providers.Provider)
	}
	s.providers[id] = provider
	return provider, nil
}
