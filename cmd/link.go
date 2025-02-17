package cmd

import (
	"fmt"

	"github.com/chrisnharvey/confman/internal/config"
	"github.com/chrisnharvey/confman/internal/fs/link"
	"github.com/spf13/cobra"
)

type LinkCmd struct {
	Config      *config.Config
	LinkFactory LinkFactory
}

func NewLinkCmd(config *config.Config, linkFactory LinkFactory) *LinkCmd {
	return &LinkCmd{
		Config:      config,
		LinkFactory: linkFactory,
	}
}

func (l *LinkCmd) Register(rootCmd *cobra.Command) {
	rootCmd.AddCommand(l.GetCmd())
}

func (l *LinkCmd) GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "link",
		Short: "Creates symlinks for all configuration files managed by Confman",
		RunE:  l.RunLinkCmd,
		// This should list all the files in the yaml file.
	}
}

func (l *LinkCmd) RunLinkCmd(cmd *cobra.Command, args []string) error {
	// Who is the current user?
	// Get the config file from the user's home directory?

	// Read the yaml file

	links := link.NewLinks(l.Config.Paths)

	for _, link := range links {
		fmt.Printf("%s -> %s", link.Destination, link.Source)

		if !link.DestinationExists() {
			fmt.Println(" (file missing from repository)")
			continue
		}

		if link.IsLinked() {
			fmt.Println(" (already linked)")
			continue
		}

		if link.SourceExists() && link.IsSourceSymlink() {
			fmt.Println(" (invalid source symlink)")
			continue
		}

		if !link.CanBeLinked() {
			fmt.Println(" (can not be linked)")
			continue
		}

		if err := link.Link(); err != nil {
			fmt.Println(" (error creating link)")
			continue
		}

		fmt.Println(" (linked)")
	}

	return nil
}
