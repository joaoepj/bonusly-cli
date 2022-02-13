package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

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

func fetchCurrentGivingBalance() int {
	return 5
}

func GetUser(id, apiToken string) ([]byte, error) {
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

func ReadApiToken() string {
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
