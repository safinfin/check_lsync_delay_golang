package cmd

import (
	"errors"
	"os"

	"github.com/mackerelio/checkers"
	"github.com/safinfin/check_lsync_delay_golang/actions"
	"github.com/spf13/cobra"
)

var (
	file     string
	warning  int64
	critical int64

	rootCmd = &cobra.Command{
		Use:   "check_lsync_delay",
		Short: "A nagios plugins to check sync delay of lsyncd",
		Long:  `A nagios plugins to check sync delay of lsyncd`,
		Args: func(cmd *cobra.Command, args []string) error {
			if warning == 0 && critical == 0 {
				return errors.New("at least one of either warning or critical is required")
			}
			if warning < 0 || critical < 0 {
				return errors.New("warning/critical must be greater than or equal to 0")
			}
			if critical != 0 && warning > critical {
				return errors.New("critical must be greater than or equal to warning")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ckr := actions.Run(file, warning, critical)
			ckr.Name = "lsync delay"
			ckr.Exit()
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(int(checkers.UNKNOWN))
	}
}

func init() {
	rootCmd.Flags().SortFlags = false
	rootCmd.Flags().StringVarP(&file, "file", "f", "", "the path to file")
	rootCmd.MarkFlagRequired("file")
	rootCmd.Flags().Int64VarP(&warning, "warning", "w", 0, "the number of seconds result in warning status")
	rootCmd.Flags().Int64VarP(&critical, "critical", "c", 0, "the number of seconds result in critical status")
}
