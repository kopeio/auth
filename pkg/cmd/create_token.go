package cmd

import (
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"kope.io/auth/pkg/apis/auth/v1alpha1"
	"time"
	"strconv"
	crypto_rand "crypto/rand"
)

type CreateTokenOptions struct {
	Username string

	ID    string
	Token []byte
}

func RunCreateToken(f Factory, out io.Writer, o *CreateTokenOptions) error {
	if o.Username == "" {
		return fmt.Errorf("must specify username for token creation")
	}
	clientset, err := f.Clientset()
	if err != nil {
		return err
	}

	user, err := clientset.AuthV1alpha1().Users().Get(o.Username, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error reading user: %v", err)
	}

	if o.ID == "" {
		o.ID = strconv.FormatInt(time.Now().Unix(), 10)
	}

	rawSecret := o.Token
	if len(rawSecret) == 0 {

		secret := make([]byte, 32, 32)
		_, err := crypto_rand.Read(secret)
		if err != nil {
			return fmt.Errorf("error generating random token: %v", err)
		}
		rawSecret = secret
	}

	token := &v1alpha1.TokenSpec{
		ID:        o.ID,
		RawSecret: rawSecret,
	}
	user.Spec.Tokens = append(user.Spec.Tokens, token)

	_, err = clientset.AuthV1alpha1().Users().Update(user)
	if err != nil {
		return fmt.Errorf("error writing updated user: %v", err)
	}

	fmt.Fprintf(out, "created token for user: %s\n", user.Name)

	return nil
}
