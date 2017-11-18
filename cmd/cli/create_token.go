package main

import (
	"io"

	"github.com/spf13/cobra"
	"kope.io/auth/pkg/cmd"
	"fmt"
)

var (
	create_token_long = longDescription(`
	token <username>
	`)

	create_token_short = shortDescription(`token <username>`)
)

func NewCmdCreateToken(f cmd.Factory, out io.Writer) *cobra.Command {
	options := &cmd.CreateTokenOptions{}

	cmd := &cobra.Command{
		Use:   "token",
		Short: create_token_short,
		Long:  create_token_long,
		Args: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("must specify name of user for whom we are creating the token")
			}
			options.Username = args[0]
			return nil
		},
		Run: func(c *cobra.Command, args []string) {
			err := cmd.RunCreateToken(f, out, options)
			if err != nil {
				exitWithError(err)
			}
		},
	}

	return cmd
}
