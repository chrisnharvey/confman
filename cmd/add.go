package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/chrisnharvey/confman/internal/config"
	"github.com/chrisnharvey/confman/internal/fs/link"
	"github.com/spf13/cobra"
)

type AddCmd struct {
	Config      *config.Config
	LinkFactory LinkFactory
}

//go:generate mockery --name LinkFactory
type LinkFactory interface {
	NewLink(source, destination string) link.Link
}

func NewAddCmd(config *config.Config, linkFactory LinkFactory) *AddCmd {
	return &AddCmd{
		Config:      config,
		LinkFactory: linkFactory,
	}
}

func (l *AddCmd) Register(rootCmd *cobra.Command) {
	rootCmd.AddCommand(l.GetCmd())
}

func (l *AddCmd) GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add a configuration file to the Confman repository [source] [destination]",
		// This should move the config file to the destination and add the mapping in the yaml
		// file. A symlink should be created in the source location pointing to the destination.
		RunE: l.RunAddCmd,
		Args: cobra.ExactArgs(2),
	}
}

func (l *AddCmd) RunAddCmd(cmd *cobra.Command, args []string) error {
	// Move the config file to the destination
	// Add the mapping in the yaml file
	// Create a symlink in the source location pointing to the destination

	filePath, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	if _, exists := l.Config.Paths[filePath]; exists {
		return fmt.Errorf("mapping already exists for %s", filePath)
	}

	link := l.LinkFactory.NewLink(filePath, args[1])

	err = link.Create()
	if err != nil {
		return err
	}

	if err := l.Config.AddPath(filePath, args[1]); err != nil {
		_ = link.Restore()
		return err
	}

	return nil
}
