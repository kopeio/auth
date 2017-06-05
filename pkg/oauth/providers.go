package oauth

import (
	"fmt"
	"github.com/golang/glog"
	"kope.io/auth/pkg/oauth/providers"
	"kope.io/auth/pkg/oauth/providers/google"
)

func (s *Server) getProvider(id string) (providers.Provider, error) {
	config, err := s.Config.AuthProvider(id)
	if err != nil {
		return nil, fmt.Errorf("error getting configuration for %q: %v", id, err)
	}

	if config == nil {
		glog.Warningf("provider %q not configured", id)
		return nil, nil
	}

	s.providersMutex.Lock()
	defer s.providersMutex.Unlock()

	existing := s.providers[id]
	if existing != nil && existing.Config().ResourceVersion == config.ResourceVersion {
		return existing, nil
	}

	provider := google.NewGoogleProvider(config)
	if s.providers == nil {
		s.providers = make(map[string]providers.Provider)
	}
	s.providers[id] = provider
	return provider, nil
}
