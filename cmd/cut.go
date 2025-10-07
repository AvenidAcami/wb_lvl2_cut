package cmd

import "github.com/spf13/cobra"

var cutCmd = &cobra.Command{}

func init() {
	rootCmd.AddCommand(cutCmd)
}
