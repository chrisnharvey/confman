package cmd

import (
	"github.com/chrisnharvey/confman/internal/config"
	"github.com/spf13/cobra"
)

type InitCmd struct {
	Config *config.Config
}

func NewInitCmd(config *config.Config) *InitCmd {
	return &InitCmd{
		Config: config,
	}
}

func (l *InitCmd) Register(rootCmd *cobra.Command) {
	rootCmd.AddCommand(l.GetCmd())
}

func (l *InitCmd) GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize the Confman repository in the given location",
		RunE:  l.RunInitCmd,
	}
}

func (l *InitCmd) RunInitCmd(cmd *cobra.Command, args []string) error {
	return nil
}
