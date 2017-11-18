package main

import (
	goflag "flag"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"kope.io/auth/pkg/cmd"
)

var (
	root_long = longDescription(`
	kopeio-auth CLI tool
	`)

	root_short = shortDescription(`kopeio-auth CLI tool`)
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func NewCmdRoot(out io.Writer) *cobra.Command {
	options := &cmd.FactoryOptions{}

	home := homeDir()
	if home != "" {
		options.Kubeconfig = filepath.Join(home, ".kube", "config")
	}

	f := cmd.NewDefaultFactory(options)

	cmd := &cobra.Command{
		Use:   "kopeio-auth",
		Short: root_short,
		Long:  root_long,
	}

	cmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)

	cmd.PersistentFlags().StringVar(&options.Kubeconfig, "kubeconfig", options.Kubeconfig, "Path to the kubeconfig file to use for CLI requests.")

	// create subcommands
	//cmd.AddCommand(NewCmdCompletion(f, out))
	cmd.AddCommand(NewCmdCreate(f, out))
	cmd.AddCommand(NewCmdExport(f, out))

	return cmd
}
