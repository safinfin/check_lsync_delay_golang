package cmd

import (
	"errors"
	"os"
	"time"

	"github.com/mackerelio/checkers"
	"github.com/spf13/cobra"
	"gitlab.heartbeats.jp/kaji/check_lsync_delay_golang/actions"
)

var opts actions.DelayOptions

var rootCmd = &cobra.Command{
	Use:     "check_lsync_delay",
	Short:   "check_lsync_delay is a nagios plugins to check file sync delay of lsyncd",
	Long:    `check_lsync_delay is a nagios plugins to check file sync delay of lsyncd.`,
	Version: Version,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(os.Args) == 1 {
			cmd.Help()
			os.Exit(int(checkers.UNKNOWN))
		}
		if opts.File == "" {
			return errors.New("required flag `-f (--file)` not set")
		}
		if opts.Warning == 0 && opts.Critical == 0 {
			return errors.New("at least one of either warning or critical is required")
		}
		if opts.Warning < 0 || opts.Critical < 0 {
			return errors.New("warning and critical must be greater than or equal to 0")
		}
		if opts.Critical != 0 && opts.Warning > opts.Critical {
			return errors.New("critical must be greater than or equal to warning")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		opts.Time = time.Now()
		opts.Check().Exit()
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(int(checkers.UNKNOWN))
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().StringVarP(&opts.File, "file", "f", "", "the file synchronized by lsyncd (required)")
	rootCmd.PersistentFlags().IntVarP(&opts.Warning, "warning", "w", 0, "the number of seconds result in warning status")
	rootCmd.PersistentFlags().IntVarP(&opts.Critical, "critical", "c", 0, "the number of seconds result in critical status")
}
