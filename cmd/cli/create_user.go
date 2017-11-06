package main

import (
	"io"

	"github.com/spf13/cobra"
	"kope.io/auth/pkg/cmd"
	"fmt"
)

var (
	create_user_long = longDescription(`
	create <username>
	`)

	create_user_short = shortDescription(`create <username>`)
)

func NewCmdCreateUser(f cmd.Factory, out io.Writer) *cobra.Command {
	options := &cmd.CreateUserOptions{}

	cmd := &cobra.Command{
		Use:   "user",
		Short: create_user_short,
		Long:  create_user_long,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("must specify name of user to create")
			}
			options.Username = args[0]
			return nil
		},
		Run: func(c *cobra.Command, args []string) {
			err := cmd.RunCreateUser(f, out, options)
			if err != nil {
				exitWithError(err)
			}
		},
	}

	return cmd
}
