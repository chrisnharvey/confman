package cmd

import (
	"fmt"

	"github.com/chrisnharvey/confman/internal/config"
	"github.com/chrisnharvey/confman/internal/fs"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Config *config.Config
}

func NewListCmd(config *config.Config) *ListCmd {
	return &ListCmd{
		Config: config,
	}
}

func (l *ListCmd) Register(rootCmd *cobra.Command) {
	rootCmd.AddCommand(l.GetCmd())
}

func (l *ListCmd) GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all configuration files managed by the manager",
		RunE: l.RunListCmd,
	}
}

func (l *ListCmd) RunListCmd(cmd *cobra.Command, args []string) error {
	links := fs.NewLinks(l.Config.Paths)

	for _, link := range links {
		fmt.Printf("%s -> %s", link.Destination, link.Source)

		if !link.DestinationExists() {
			fmt.Println(" (file missing from repository)")
			continue
		}

		if !link.SourceExists() {
			fmt.Println(" (symlink missing)")
			continue
		}

		if !link.IsSourceSymlink() {
			fmt.Println(" (not a symlink)")
			continue
		}

		if !link.IsLinked() {
			fmt.Println(" (symlink target mismatch)")
			continue
		}

		fmt.Println(" (ok)")
	}

	return nil
}