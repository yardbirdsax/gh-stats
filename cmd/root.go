/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	startDateString *string
	startDate time.Time
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-stats",
	Short: "An extension for the GitHub CLI for getting statistics about various things",
	Long: ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	defaultStartDate := time.Now().Add(-7 * 24 * time.Hour).Truncate(24 * time.Hour).Format("2006-01-02")
	startDateString = rootCmd.PersistentFlags().String("start-date", defaultStartDate, "the date at which to start when qualifying searches")
	formats := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02"}
	var err error
	for _, format := range formats {
		startDate, err = time.Parse(format, *startDateString)
		if err == nil {
			break
		}
	}
	if err != nil {
		panic("value entered for start time is not a parseable date/time string")
	}
}
