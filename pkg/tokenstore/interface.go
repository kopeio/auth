package tokenstore

import (
	authenticationv1beta1 "k8s.io/api/authentication/v1beta1"
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
