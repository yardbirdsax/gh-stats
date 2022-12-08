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
	endDateString *string
	endDate time.Time
	groupByField *string
	asCSV *bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-stats",
	Short: "An extension for the GitHub CLI for getting statistics about various things",
	Long: ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		formats := []string{time.RFC3339Nano, time.RFC3339, "2006-01-02"}
		for _, format := range formats {
			startDate, err = time.Parse(format, *startDateString)
			if err == nil {
				break
			}
		}
		if err != nil {
			return err
		}
		for _, format := range formats {
			endDate, err = time.Parse(format, *endDateString)
			if err == nil {
				break
			}
		}
		return err
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
	defaultStartDate := time.Now().Add(-7 * 24 * time.Hour).Truncate(24 * time.Hour).Format("2006-01-02")
	startDateString = rootCmd.PersistentFlags().String("start-date", defaultStartDate, "the date at which to start when qualifying searches")
	defaultEndDate := time.Now().Truncate(24 * time.Hour).Format("2006-01-02")
	endDateString = rootCmd.PersistentFlags().String("end-date", defaultEndDate, "the date at which to end when qualifying searches")
	groupByField = rootCmd.PersistentFlags().String("group-by", "CreatedAt", "the field on which to group results")
	asCSV = rootCmd.PersistentFlags().Bool("csv", false, "output results as CSV")
}
