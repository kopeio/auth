package main

import (
	"io"

	"github.com/spf13/cobra"
	"kope.io/auth/pkg/cmd"
)

var (
	create_user_long = longDescription(`
	create
	`)

	create_user_short = shortDescription(`create`)
)

func NewCmdCreateUser(f cmd.Factory, out io.Writer) *cobra.Command {
	options := &cmd.CreateUserOptions{}

	cmd := &cobra.Command{
		Use:   "user",
		Short: create_user_short,
		Long:  create_user_long,
		Run: func(c *cobra.Command, args []string) {
			err := cmd.RunCreateUser(f, out, options)
			if err != nil {
				exitWithError(err)
			}
		},
	}

	return cmd
}
