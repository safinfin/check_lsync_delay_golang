package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version string = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("check_lsync_delay version", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
