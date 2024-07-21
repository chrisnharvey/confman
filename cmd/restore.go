package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/chrisnharvey/confman/internal/config"
	"github.com/chrisnharvey/confman/internal/fs/link"
	"github.com/spf13/cobra"
)

type RestoreCmd struct {
	Config      *config.Config
	LinkFactory LinkFactory
}

func NewRestoreCmd(config *config.Config, linkFactory LinkFactory) *RestoreCmd {
	return &RestoreCmd{
		Config:      config,
		LinkFactory: linkFactory,
	}
}

func (l *RestoreCmd) Register(rootCmd *cobra.Command) {
	rootCmd.AddCommand(l.GetCmd())
}

func (l *RestoreCmd) GetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore",
		Short: "Restore a config file from Confman and move it back to its original location",
		RunE:  l.RunRestoreCmd,
	}
}

func (l *RestoreCmd) RunRestoreCmd(cmd *cobra.Command, args []string) error {
	// This should remove the mapping from the yaml file and delete the symlink. Then move
	// the file back to its original location.

	filePath, err := filepath.Abs(args[0])
	if err != nil {
		return err
	}

	name, exists := l.Config.Paths[filePath]
	if !exists {
		return fmt.Errorf("file is not managed by confman %s", filePath)
	}

	link := link.NewLink(filePath, name)

	err = link.Restore()
	if err != nil {
		return err
	}

	if err := l.Config.RemovePath(filePath); err != nil {
		return err
	}

	return nil
}
