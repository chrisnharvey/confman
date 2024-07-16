package cmd

import (
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the config ",
	RunE: RunInitCmd,
	// This should list all the files in the yaml file.
}

func RunInitCmd(cmd *cobra.Command, args []string) error {
	// Who is the current user?
	// Get the config file from the user's home directory?

	// Read the yaml file


	return nil
}