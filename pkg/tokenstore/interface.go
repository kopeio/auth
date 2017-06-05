package tokenstore

import (
	"encoding/base64"

	authenticationv1beta1 "k8s.io/client-go/pkg/apis/authentication/v1beta1"
	auth "kope.io/auth/pkg/apis/auth/v1alpha1"
	"kope.io/auth/pkg/oauth/session"
)

type Interface interface {
	LookupToken(token string) (*authenticationv1beta1.UserInfo, error)
	CreateToken(u *auth.User, hashSecret bool) (*auth.TokenSpec, error)
	FindUserByUID(uid string) (*auth.User, error)
	//FindExistingUser(i *auth.IdentitySpec) (*auth.User, error)
	MapToUser(info *session.UserInfo, create bool) (*auth.User, error)

	Run(stopCh <-chan struct{})
}

type TokenInfo struct {
	UserID  string
	TokenID string
	Secret  []byte
}

func (t *TokenInfo) Encode() string {
	return t.UserID + "/" + t.TokenID + "/" + base64.URLEncoding.EncodeToString(t.Secret)
}
