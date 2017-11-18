package main

import (
	"io"
	"github.com/spf13/cobra"
	"kope.io/auth/pkg/cmd"
	"fmt"
)

var (
	export_kubecfg_long = longDescription(`
	kubecfg <username>
	`)

	export_kubecfg_short = shortDescription(`kubecfg <username>`)
)

func NewCmdExportKubecfg(f cmd.Factory, out io.Writer) *cobra.Command {
	options := &cmd.ExportKubecfgOptions{}

	cmd := &cobra.Command{
		Use:   "kubecfg <username>",
		Aliases: []string{"kubeconfig"},
		Short: export_kubecfg_short,
		Long:  export_kubecfg_long,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("must specify name of user to create a kubeconfig for")
			}
			options.Username = args[0]
			return nil
		},
		Run: func(c *cobra.Command, args []string) {
			err := cmd.RunExportKubecfg(f, out, options)
			if err != nil {
				exitWithError(err)
			}
		},
	}

	return cmd
}
