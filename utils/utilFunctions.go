package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	// "reflect"
	"time"

	"gopkg.in/yaml.v3"
)

type userApiResponse struct {
	Success bool   `json:"success"`
	Result  User   `json:"result"`
	Message []byte `json:"message"`
}

type UserData struct {
	Timestamp time.Time `json:"timestamp"`
	Data      []byte    `json:"data"`
}
type config struct {
	ApiToken string `yaml:"apiToken"`
}

type User struct {
	GivingBalance  int    `json:"giving_balance"`
	EarningBalance int    `json:"earning_balance"`
	Email          string `json:"email"`
}

func fetchCurrentGivingBalance() int {
	return 5
}

func makeRequest(method, url string, payload []byte) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	req.Header.Add("Authorization", "Bearer "+readApiToken())
	req.Header.Add("HTTP_APPLICATION_NAME", "bonuslyCLI")

	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func GetUser(id string) (User, error) {
	userData, err := makeRequest("GET", "https://bonus.ly/api/v1/users/me", []byte(""))
	if err != nil {
		fmt.Printf("something went wrong during get request")
		return User{}, err
	}
	userApiResponse := userApiResponse{}
	if err := json.Unmarshal(userData, &userApiResponse); err != nil {
		fmt.Printf("Error unmarshaling user data")
		return User{}, err
	}

	return userApiResponse.Result, nil
}

func readApiToken() string {
	dat, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	config := config{}

	if err := yaml.Unmarshal([]byte(dat), &config); err != nil {
		panic(err)
	}
	return config.ApiToken
}

func SaveUserDataToDisk(userData UserData) {
	data, err := json.Marshal(userData)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("userData.txt", data, 0666); err != nil {
		panic(err)
	}
}
func ReadUserDataFromDisk() UserData {
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

type Bonus struct {
	GiverEmail    string   `json:"giver_email"`
	ReceiverEmail string   `json:"receiver_email"`
	Amount        int      `json:"amount"`
	Hashtag       []string `json:"hashtag"`
	Reason        string   `json:"reason"`
}

func CreateBonus(payload Bonus) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error occurred")
	}
	return makeRequest("POST", "https://bonus.ly/api/v1/bonuses", data)
}
