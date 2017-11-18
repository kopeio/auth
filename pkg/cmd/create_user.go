package cmd

import (
	"fmt"
	"io"
	"kope.io/auth/pkg/apis/auth/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateUserOptions struct {
	Username string
	Groups   []string
}

func RunCreateUser(f Factory, out io.Writer, o *CreateUserOptions) error {
	if o.Username == "" {
		return fmt.Errorf("must specify username to create")
	}
	clientset, err := f.Clientset()
	if err != nil {
		return err
	}

	user := &v1alpha1.User{
		ObjectMeta: v1.ObjectMeta{
			Name: o.Username,
		},
		Spec: v1alpha1.UserSpec{
			Username: o.Username,
			Groups:   o.Groups,
		},
	}

	if _, err := clientset.AuthV1alpha1().Users().Create(user); err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	fmt.Fprintf(out, "created user: %s\n", user.Name)

	return nil
}
