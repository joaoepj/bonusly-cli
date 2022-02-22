/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "bonusly",
	Short: "CLI wrapper for Bonusly API written in Go",
	Long: `CLI wrapper for Bonusly API written in go
    This terminal applicatin provides the ability to interact with parts of the Bonusly API.
    See bonusly --help for all possible commands.`,
	Version: "0.2.0",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
