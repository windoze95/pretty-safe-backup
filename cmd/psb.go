package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	// VERSION is set during build
	VERSION string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "psb",
	Short: "A neat program to help you set up the best backup solution for your Linux, Android, and macOS devices",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	VERSION = version

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
