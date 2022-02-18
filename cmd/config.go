/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bonusly/utils"
	"fmt"

	"github.com/spf13/cobra"
)

var apiToken string

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Use this command to edit the user config",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(apiToken) == 0 {
			fmt.Println("Api token can't be empty!")
			return
		}
		userData, err := utils.ReadUserDataFromDisk(verbose)
		if err != nil {
			fmt.Println(err)
			return
		}
		userData.ApiToken = apiToken
		utils.SaveUserDataToDisk(userData)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringVar(&apiToken, "token", "", "Specify the token")
}
