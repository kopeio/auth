package session

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"kope.io/auth/pkg/oauth/pb"
)

type Session struct {
	pb.SessionData
}

type UserInfo struct {
	Email          string
	ProviderID     string
	ProviderUserID string
}

func (s *Session) IsExpired() bool {
	if s.ExpiresOn != 0 && s.ExpiresOn < time.Now().Unix() {
		return true
	}
	return false
}

func (s *Session) Age() time.Duration {
	timestamp := time.Unix(int64(s.Timestamp), 0)
	return time.Now().Truncate(time.Second).Sub(timestamp)
}

func (s *Session) Marshal() ([]byte, error) {
	b, err := proto.Marshal(&s.SessionData)
	if err != nil {
		return nil, fmt.Errorf("error serializing data: %v", err)
	}
	return b, nil
}

func UnmarshalSession(b []byte) (*Session, error) {
	s := &Session{}
	if err := proto.Unmarshal(b, &s.SessionData); err != nil {
		return nil, fmt.Errorf("error parsing cookie data: %v", err)
	}

	return s, nil
}
