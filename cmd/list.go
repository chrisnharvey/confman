package cmd

import (
	"fmt"

	"github.com/chrisnharvey/confman/internal/config"
	"github.com/chrisnharvey/confman/internal/fs/link"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Config      *config.Config
	LinkFactory LinkFactory
}

func NewListCmd(config *config.Config, linkFactory LinkFactory) *ListCmd {
	return &ListCmd{
		Config:      config,
		LinkFactory: linkFactory,
	}
}

func (l *ListCmd) Register(rootCmd *cobra.Command) {
	rootCmd.AddCommand(l.GetCmd())
}

func (l *ListCmd) GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all configuration files and their status",
		RunE:  l.RunListCmd,
	}
}

func (l *ListCmd) RunListCmd(cmd *cobra.Command, args []string) error {
	links := link.NewLinks(l.Config.Paths)

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
