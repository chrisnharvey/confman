package cmd

import (
	"github.com/chrisnharvey/confman/internal/config"
	"github.com/spf13/cobra"
)

type DetachCmd struct {
	Config *config.Config
}

func NewDetachCmd(config *config.Config) *DetachCmd {
	return &DetachCmd{
		Config: config,
	}
}

func (l *DetachCmd) Register(rootCmd *cobra.Command) {
	rootCmd.AddCommand(l.GetCmd())
}

func (l *DetachCmd) GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "detach",
		Short: "Detach a config file from Confman and move it back to its original location",
		RunE: l.RunDetachCmd,
	}
}

func (l *DetachCmd) RunDetachCmd(cmd *cobra.Command, args []string) error {
	// This should remove the mapping from the yaml file and delete the symlink. Then move
	// the file back to its original location.

	return nil
}