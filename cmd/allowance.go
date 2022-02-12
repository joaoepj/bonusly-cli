package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var force bool

type UserData struct {
	Timestamp time.Time `json:"timestamp"`
	Data      []byte    `json:"data"`
}
type Config struct {
	ApiToken string `yaml:"apiToken"`
}

type User struct {
	Result struct {
		GivingBalance  int `json:"giving_balance"`
		EarningBalance int `json:"earning_balance"`
	} `json:"result"`
}

func getUser(id, apiToken string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://bonus.ly/api/v1/users/me", nil)
	req.Header.Add("Authorization", "Bearer "+apiToken)
	req.Header.Add("HTTP_APPLICATION_NAME", "bonuslyCLI")

	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func readApiToken() string {
	dat, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	config := Config{}

	if err := yaml.Unmarshal([]byte(dat), &config); err != nil {
		panic(err)
	}
	return config.ApiToken
}

func saveUserDataToDisk(userData UserData) {
	data, err := json.Marshal(userData)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("userData.txt", data, 0666); err != nil {
		panic(err)
	}
}
func readUserDataFromDisk() UserData {
	dat, err := os.ReadFile("userData.txt")
	if err != nil {
		panic(err)
	}
	userData := UserData{}
	if err := json.Unmarshal(dat, &userData); err != nil {
		panic(err)
	}
	return userData
}

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
		userData := readUserDataFromDisk()
		if userData.Timestamp.Add(24*time.Hour).Before(time.Now()) || force == true {
			// fetch new data from server
			fmt.Println("getting new data from server")
			fmt.Printf("force flag set?: %t\n", force)
			fmt.Printf("timestamp exceedeed?: %t\n", userData.Timestamp.Add(24*time.Hour).Before(time.Now()))
			apiToken := readApiToken()
			rawUserData, userDataErr := getUser("me", apiToken)

			if userDataErr != nil {
				panic(userDataErr)
			}
			userData.Data = rawUserData
			userData.Timestamp = time.Now()
		}
		user := User{}
		if err := json.Unmarshal(userData.Data, &user); err != nil {
			panic(err)
		}
		fmt.Printf("You still have %d Bonusly left to give away this month.\n", user.Result.GivingBalance)
		fmt.Printf("You still have %d Bonusly left to spend on rewards this month.\n", user.Result.EarningBalance)
		saveUserDataToDisk(userData)
	},
}

func init() {
	rootCmd.AddCommand(allowanceCmd)

	allowanceCmd.Flags().BoolVarP(&force, "force", "f", false, "Force bonuslyCLI to fetch new data from the server")
}
