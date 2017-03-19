package providers

import (
	"testing"
	"time"

	"kope.io/auth/pkg/assert"
	"kope.io/auth/pkg/cookie/pb"
)

func TestRefresh(t *testing.T) {
	p := &ProviderData{}
	refreshed, err := p.RefreshSessionIfNeeded(&SessionState{
		pb.SessionData{
			ExpiresOn: time.Now().Unix() - (11 * 60),
		},
	})
	assert.Equal(t, false, refreshed)
	assert.Equal(t, nil, err)
}
