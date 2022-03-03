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
	Long: `With this command you can configure your settings. Start off by running
"bonusly config --token <your-api-token>" to be able to use the Bonusly CLI.

If you don't have a Bonusly API token yet, go visit https://bonus.ly/api to create one.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Flags())
		if cmd.Flag("token").Changed {
			if apiToken == "" {
				fmt.Println("Api token can't be empty!")
				return
			}
		} else {
			fmt.Println("Please specify at least one flag! See \"bonusly config --help\" for more information.")
			return
		}
		userData, err := utils.ReadUserDataFromDisk(verbose)
		if err != nil {
			if verbose {
				fmt.Println(err)
			}
		}
		userData.ApiToken = apiToken
		utils.SaveUserDataToDisk(userData)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringVar(&apiToken, "token", "", "Specify the token")
}