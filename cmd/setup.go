package cmd

import (
	"github.com/orange-lightsaber/pretty-safe-backup/setup"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Create a backup operation",
	Long: `Starts a step-by-step process that assists in creating a new backup operation.
You will need to run this to properly configure the backup services.`,
	Run: func(cmd *cobra.Command, args []string) {
		setup.Build()
	},
}

func init() {
	RootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
