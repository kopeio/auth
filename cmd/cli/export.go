package main

import (
	"io"

	"github.com/spf13/cobra"
	"kope.io/auth/pkg/cmd"
)

var (
	export_long = longDescription(`
	export
	`)

	export_short = shortDescription(`export`)
)

func NewCmdExport(f cmd.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: export_short,
		Long:  export_long,
	}

	// export subcommands
	cmd.AddCommand(NewCmdExportKubecfg(f, out))

	return cmd
}
