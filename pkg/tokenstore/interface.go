package tokenstore

import (
	authenticationv1beta1 "k8s.io/client-go/pkg/apis/authentication/v1beta1"
	"encoding/base64"
)

type Interface interface {
	LookupToken(token string) (*authenticationv1beta1.UserInfo, error)
	CreateToken(userInfo *authenticationv1beta1.UserInfo) (*TokenInfo, error)
	ListTokens(uid string) ([]*TokenInfo, error)

	Run(stopCh <-chan struct{})
}

type TokenInfo struct {
	UserID  string
	TokenID string
	Secret  []byte
}

func (t*TokenInfo) Encode() string {
	return t.UserID + "/" + t.TokenID + "/" + base64.URLEncoding.EncodeToString(t.Secret)
}