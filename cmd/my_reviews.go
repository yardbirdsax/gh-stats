/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
	"github.com/yardbirdsax/gh-stats/internal/pr"
)

// mineCmd represents the mine command
var mineCmd = &cobra.Command{
	Use:   "my-reviews",
	Short: "PRs you've reviewed",
	Run: func(cmd *cobra.Command, args []string) {
		myReviews, err := pr.MyReviews(startDate, endDate, *groupByField)
		if err != nil {
			fmt.Printf("error: %v", err)
			os.Exit(1)
		}
		var out string
		if !*asCSV {
			r, _ := glamour.NewTermRenderer()
			out, err = r.Render(myReviews.AsMarkdownTable())
			if err != nil {
				fmt.Printf("error: %v", err)
			}
		} else {
			out = myReviews.AsCSV()
		}
		fmt.Print(out)
	},
}

func init() {
	prCmd.AddCommand(mineCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mineCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mineCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
