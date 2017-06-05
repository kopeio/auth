package oauth

import (
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"kope.io/auth/pkg/oauth/pb"
)

type State struct {
	pb.StateData
}

func unmarshalState(value string) (*State, error) {
	b, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("error decoding state date: %v", err)
	}
	s := &State{}
	if err := proto.Unmarshal(b, &s.StateData); err != nil {
		return nil, fmt.Errorf("error parsing state data: %v", err)
	}
	return s, nil
}

func (s *State) Marshal() (string, error) {
	b, err := proto.Marshal(&s.StateData)
	if err != nil {
		return "", fmt.Errorf("error serializing data: %v", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
