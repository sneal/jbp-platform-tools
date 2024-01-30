package cmd

import (
	"github.com/spf13/cobra"
)

const cfClientKey = "CFClient"

var rootCmd = &cobra.Command{
	Use:   "platform-tools",
	Short: "Tools to help rationalize about your CF system.",
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
