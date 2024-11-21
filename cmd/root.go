package cmd

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var (
	// AppName for server
	AppName string
	// AppVersion for server
	AppVersion string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:              "data-transmission",
	TraverseChildren: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	spew.Config = *spew.NewDefaultConfig()
	spew.Config.ContinueOnMethod = true

	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

}
