package cmd

import (
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateUserOptions struct {
}

func RunCreateUser(f Factory, out io.Writer, o *CreateUserOptions) error {
	clientset, err := f.Clientset()
	if err != nil {
		return err
	}

	allNamespaces := ""
	users, err := clientset.AuthV1alpha1().Users(allNamespaces).List(v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error listing users: %v", err)
	}

	for _, u := range users.Items {
		if _, err := fmt.Fprintf(out, "%s", u.Name); err != nil {
			return fmt.Errorf("error writing to output: %v", err)
		}
	}

	return nil
}
