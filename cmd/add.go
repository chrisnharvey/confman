package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/chrisnharvey/confman/internal/config"
	"github.com/spf13/cobra"
)

type AddCmd struct {
	Config *config.Config
}

func NewAddCmd(config *config.Config) *AddCmd {
	return &AddCmd{
		Config: config,
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

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	destPath := "/confman/" + args[1]
	if _, err := os.Stat(destPath); err == nil {
		return fmt.Errorf("file already exists at %s", destPath)
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, file)
    if err != nil {
        return fmt.Errorf("could not copy file: %v", err)
    }

	err = os.Remove(filePath)
    if err != nil {
		os.Remove(destPath) // cleanup
        return fmt.Errorf("could not remove source file: %v", err)
    }

	if err := os.Symlink(destPath, filePath); err != nil {
		return err
	}

	if err := l.Config.AddPath(filePath, args[1]); err != nil {
		os.Remove(destPath) // cleanup
		return err
	}

	return nil
}