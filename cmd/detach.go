package cmd

import "github.com/spf13/cobra"

var DetachCmd = &cobra.Command{
	Use:   "detach",
	Short: "Detach a configuration file from the manager [file]",
	// This should remove the mapping from the yaml file and delete the symlink. Then move
	// the file back to its original location.
}