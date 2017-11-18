package main

import (
	"io"

	"github.com/spf13/cobra"
	"kope.io/auth/pkg/cmd"
)

var (
	create_long = longDescription(`
	create
	`)

	create_short = shortDescription(`create`)
)

func NewCmdCreate(f cmd.Factory, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: create_short,
		Long:  create_long,
	}

	// create subcommands
	cmd.AddCommand(NewCmdCreateToken(f, out))
	cmd.AddCommand(NewCmdCreateUser(f, out))

	return cmd
}
