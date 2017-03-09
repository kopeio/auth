package providers

import (
	"kope.io/auth/pkg/assert"
	"kope.io/auth/pkg/cookie"
	"kope.io/auth/pkg/cookie/proto"
	"strings"
	"testing"
	"time"
)

const secret = "0123456789abcdefghijklmnopqrstuv"
const altSecret = "0000000000abcdefghijklmnopqrstuv"

func TestSessionStateSerialization(t *testing.T) {
	c, err := cookie.NewCipher([]byte(secret))
	assert.Equal(t, nil, err)
	c2, err := cookie.NewCipher([]byte(altSecret))
	assert.Equal(t, nil, err)
	s := &SessionState{
		proto.SessionData{
			Email:        "user@domain.com",
			AccessToken:  "token1234",
			ExpiresOn:    time.Now().Unix() + 3600,
			RefreshToken: "refresh4321",
		},
	}
	encoded, err := s.EncodeSessionState(c)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, strings.Count(encoded, "|"))

	ss, err := DecodeSessionState(encoded, c)
	t.Logf("%#v", ss)
	assert.Equal(t, nil, err)
	assert.Equal(t, s.Email, ss.Email)
	assert.Equal(t, s.AccessToken, ss.AccessToken)
	assert.Equal(t, s.ExpiresOn, ss.ExpiresOn)
	assert.Equal(t, s.RefreshToken, ss.RefreshToken)

	// ensure a different cipher can't decode properly (ie: it gets gibberish)
	ss, err = DecodeSessionState(encoded, c2)
	t.Logf("%#v", ss)
	assert.Equal(t, nil, err)
	assert.Equal(t, s.Email, ss.Email)
	assert.Equal(t, s.ExpiresOn, ss.ExpiresOn)
	assert.NotEqual(t, s.AccessToken, ss.AccessToken)
	assert.NotEqual(t, s.RefreshToken, ss.RefreshToken)
}

func TestSessionStateSerializationNoCipher(t *testing.T) {

	s := &SessionState{
		proto.SessionData{
			Email:        "user@domain.com",
			AccessToken:  "token1234",
			ExpiresOn:    time.Now().Unix() + 3600,
			RefreshToken: "refresh4321",
		},
	}
	encoded, err := s.EncodeSessionState(nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, s.Email, encoded)

	// only email should have been serialized
	ss, err := DecodeSessionState(encoded, nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, s.Email, ss.Email)
	assert.Equal(t, "", ss.AccessToken)
	assert.Equal(t, "", ss.RefreshToken)
}

func TestSessionStateUserOrEmail(t *testing.T) {

	s := &SessionState{
		proto.SessionData{
			Email: "user@domain.com",
			User:  "just-user",
		},
	}
	assert.Equal(t, "user@domain.com", s.userOrEmail())
	s.Email = ""
	assert.Equal(t, "just-user", s.userOrEmail())
}

func TestExpired(t *testing.T) {
	s := &SessionState{proto.SessionData{ExpiresOn: time.Now().Unix() - 60}}
	assert.Equal(t, true, s.IsExpired())

	s = &SessionState{proto.SessionData{ExpiresOn: time.Now().Unix() + 60}}
	assert.Equal(t, false, s.IsExpired())

	s = &SessionState{}
	assert.Equal(t, false, s.IsExpired())
}
