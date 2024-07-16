package main

import (
	"fmt"
	"os"

	"github.com/chrisnharvey/confman/cmd"
	"github.com/chrisnharvey/confman/internal/config"
	"github.com/spf13/cobra"
)

type Settings struct {
	ConfigPath string `mapstructure:"CONFMAN_CONFIG_PATH"`
}

var rootCmd = &cobra.Command{
  // Use:   "confman",
  Short: "Simple configuration manager",
}

type registerableCmd interface {
	Register(*cobra.Command)
}

func main() {
	// configFile, err := config.GetConfigFile()
	// if err != nil {
	// 	panic(err)
	// }

	// config, err := os.Open(configFile)
	// if err != nil {
	// 	fmt.Println(
	// 		fmt.Errorf("error opening config file: %s", err),
	// 	)
	// 	os.Exit(1)
	// }

	// viper.SetConfigType("envfile")
	// // stuff, _ := os.ReadFile(configFile)
	// // fmt.Println(string(stuff))
	// err = viper.ReadConfig(config)
	// if err != nil {
	// 	panic(err)
	// }

	// settings := &Settings{
	// 	ConfigPath: "/confman",
	// }

	// err = viper.Unmarshal(settings)
	// if err != nil {
	// 	panic(err)
	// }

	// config.GetConfigFrom(configFile)

	cfg, err := config.GetConfigFrom("/confman/.confman.yaml")
	if err != nil {
		panic(err)
	}

	// rootCmd.AddCommand(cmd.AddCmd)
	registerCmd(rootCmd, cmd.NewDetachCmd(cfg))
	registerCmd(rootCmd, cmd.NewListCmd(cfg))
	registerCmd(rootCmd, cmd.NewAddCmd(cfg))
	registerCmd(rootCmd, cmd.NewLinkCmd(cfg))
	// rootCmd.AddCommand(gitCmd)
	// rootCmd.AddCommand(cdCmd)

	if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func registerCmd(rootCmd *cobra.Command, cmd registerableCmd) {
	cmd.Register(rootCmd)
}