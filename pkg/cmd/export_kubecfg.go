package cmd

import (
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"kope.io/auth/pkg/kubeconfig"
)

type ExportKubecfgOptions struct {
	Username string
	Groups   []string
}

func RunExportKubecfg(f Factory, out io.Writer, o *ExportKubecfgOptions) error {
	if o.Username == "" {
		return fmt.Errorf("must specify username for whom we should generate a kubecfg")
	}

	clientset, err := f.Clientset()
	if err != nil {
		return err
	}

	user, err := clientset.AuthV1alpha1().Users().Get(o.Username, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error reading user: %v", err)
	}

	token := kubeconfig.FindBestToken(user)
	if token == nil {
		return fmt.Errorf("user has no tokens")
	}

	config, err := f.Config()
	if err != nil {
		return err
	}

	b, err := kubeconfig.BuildKubeconfig(config.Host, config.CAData, user, token)
	if err != nil {
		return err
	}

	if _, err := out.Write(b); err != nil {
		return err
	}

	return nil
}
