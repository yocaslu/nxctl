/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"nxctl/internal/utils/nxlog"
	"os"

	"github.com/spf13/cobra"
)

var debug bool = false
var MODULE_NAME string = "cmd"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nxctl",
	Short: "Nexcloud CLI Helper",
	Long: `nxctl is a command line application that simplify
	installation and managment of Nextcloud.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		nxlog.InitLogger(debug)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.test.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolVar(&debug, "debug", false, "Enable debug information")
}
