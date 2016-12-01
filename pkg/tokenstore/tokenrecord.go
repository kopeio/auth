package tokenstore

import (
	authenticationv1beta1 "k8s.io/client-go/pkg/apis/authentication/v1beta1"
)

type TokenRecord struct {
	ID         string `json:"id,omitempty"`
	Secret     []byte `json:"hashed,omitempty"`
}

type UserRecord struct {
	User *authenticationv1beta1.UserInfo `json:"user,omitempty"`
}
