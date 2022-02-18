package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"time"

	"gopkg.in/yaml.v3"
)

const BASE_URL = "https://bonus.ly/api/v1"
const USER_DATA_FILE = ".userdata"

type userApiResponse struct {
	Success bool   `json:"success"`
	Result  User   `json:"result"`
	Message []byte `json:"message"`
}

type bonusApiResponse struct {
	Success bool `json:"success"`
	Result  struct {
		Id string `json:"id"`
	} `json:"result"`
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

func FetchCurrentGivingBalance() int {
	user, err := GetUser("me")
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return user.GivingBalance
}

func makeRequest(method, url string, payload []byte) ([]byte, error) {
	client := &http.Client{}

	var req *http.Request
	var err error
	if method == http.MethodGet {
		req, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Printf("error occured")
		}
	} else if method == http.MethodPost {
		payloadBody := bytes.NewBuffer(payload)
		req, err = http.NewRequest(http.MethodPost, url, payloadBody)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", "application/json")
	} else {
		fmt.Printf("Unknown method %s", method)
		return nil, errors.New("unknown method")
	}
	req.Header.Add("Authorization", "Bearer "+readApiToken())
	req.Header.Add("HTTP_APPLICATION_NAME", "bonuslyCLI")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func GetLocalUser(verbose bool) (User, error) {
	userData := ReadUserDataFromDisk(verbose)
	// TODO: handle non-existant user data
	user := User{}
	if err := json.Unmarshal(userData.Data, &user); err != nil {
		return User{}, err
	}
	return user, nil
}

func GetUser(id string) (User, error) {
	userData, err := makeRequest(http.MethodGet, BASE_URL+"/users/me", []byte(""))
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
	if err := os.WriteFile(USER_DATA_FILE, data, 0666); err != nil {
		panic(err)
	}
}
func ReadUserDataFromDisk(verbose bool) UserData {
	dat, err := os.ReadFile(USER_DATA_FILE)
	if err != nil {
		if verbose {
			fmt.Println("No userdata file exists. Returning empty user.")
		}
		return UserData{}
	}
	userData := UserData{}
	if err := json.Unmarshal(dat, &userData); err != nil {
		panic(err)
	}
	return userData
}

type Bonus struct {
	Amount int    `json:"amount"`
	Reason string `json:"reason"`
}

// returns ID of post if created successfully
func CreateBonus(payload Bonus) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	bonusData, err := makeRequest(http.MethodPost, BASE_URL+"/bonuses", data)
	bonus := bonusApiResponse{}
	if err := json.Unmarshal(bonusData, &bonus); err != nil {
		return "", err
	}
	return bonus.Result.Id, nil
}
