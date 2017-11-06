package kubeconfig

import (
	"encoding/base64"
	"kope.io/auth/pkg/apis/auth/v1alpha1"
)

func FindBestToken(user *v1alpha1.User) (*v1alpha1.TokenSpec) {
	var bestToken *v1alpha1.TokenSpec

	for _, t := range user.Spec.Tokens {
		// TODO: Pick best token
		bestToken = t
	}

	return bestToken
}

func EncodeToken(user *v1alpha1.User, token *v1alpha1.TokenSpec) string {
	return string(user.UID) + "/" + token.ID + "/" + base64.URLEncoding.EncodeToString(token.RawSecret)
}
