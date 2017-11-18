package main

import (
	goflag "flag"
	"fmt"
	"os"
)

func main() {
	Execute()
}

// exitWithError will terminate execution with an error result
// It prints the error to stderr and exits with a non-zero exit code
func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "\n%v\n", err)
	os.Exit(1)
}

func Execute() {
	goflag.Set("logtostderr", "true")
	goflag.CommandLine.Parse([]string{})

	c := NewCmdRoot(os.Stdout)

	if err := c.Execute(); err != nil {
		exitWithError(err)
	}
}
