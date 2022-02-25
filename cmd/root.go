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
	Long: `bonusly is a CLI wrapper for bonus.ly API written in go
This application provides the ability to interact with parts of the Bonusly API.`,
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
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}
