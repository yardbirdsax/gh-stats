/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// prCmd represents the pr command
var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Show statisctics around pull requests.",
	Long: `Shows statistics around pull requests, such as reviews, opened, etc.`,
}

func init() {
	rootCmd.AddCommand(prCmd)
}
