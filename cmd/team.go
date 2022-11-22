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

var (
	teamName *string
	orgName  *string
)

// teamCmd represents the mine command
var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		teamReviwes, err := pr.TeamReviews(*orgName, *teamName, startDate)
		if err != nil {
			fmt.Printf("error: %v", err)
			os.Exit(1)
		}
		r, _ := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
		)
		out, err := r.Render(teamReviwes.AsMarkdownTable())
		if err != nil {
			fmt.Printf("error: %v", err)
		}
		fmt.Print(out)
	},
}

func init() {
	prCmd.AddCommand(teamCmd)
	orgName = teamCmd.PersistentFlags().String("org-name", "", "The name of the organization that the team belongs to.")
	teamName = teamCmd.PersistentFlags().String("name", "", "The name of the team to view PR statistics for.")
}
