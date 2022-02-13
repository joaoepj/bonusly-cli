package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"bonusly/utils"

	"github.com/spf13/cobra"
)

var force bool

// allowanceCmd represents the allowance command
var allowanceCmd = &cobra.Command{
	Use:   "allowance",
	Short: "Get your current Bonuslys for spending and giving away",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		userData := utils.ReadUserDataFromDisk()
		if userData.Timestamp.Add(24*time.Hour).Before(time.Now()) || force == true {
			// fetch new data from server
			fmt.Println("getting new data from server")
			fmt.Printf("force flag set?: %t\n", force)
			fmt.Printf("timestamp exceedeed?: %t\n", userData.Timestamp.Add(24*time.Hour).Before(time.Now()))
			apiToken := utils.ReadApiToken()
			rawUserData, userDataErr := utils.GetUser("me", apiToken)

			if userDataErr != nil {
				panic(userDataErr)
			}
			userData.Data = rawUserData
			userData.Timestamp = time.Now()
		}
		user := utils.User{}
		if err := json.Unmarshal(userData.Data, &user); err != nil {
			panic(err)
		}
		fmt.Printf("You still have %d Bonusly left to give away this month.\n", user.Result.GivingBalance)
		fmt.Printf("You still have %d Bonusly left to spend on rewards this month.\n", user.Result.EarningBalance)
		utils.SaveUserDataToDisk(userData)
	},
}

func init() {
	rootCmd.AddCommand(allowanceCmd)

	allowanceCmd.Flags().BoolVarP(&force, "force", "f", false, "Force bonuslyCLI to fetch new data from the server")
}
