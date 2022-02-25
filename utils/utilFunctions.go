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
	ApiToken  string    `json:"apiToken"`
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
	apiToken, err := readApiToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+apiToken)
	req.Header.Add("HTTP_APPLICATION_NAME", "bonuslyCLI")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func GetLocalUser(verbose bool) (User, error) {
	userData, err := ReadUserDataFromDisk(verbose)
	if err != nil {
		return User{}, err
	}
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
		fmt.Printf("Something went wrong during get request")
		return User{}, err
	}
	userApiResponse := userApiResponse{}
	if err := json.Unmarshal(userData, &userApiResponse); err != nil {
		fmt.Printf("Error unmarshaling user data")
		return User{}, err
	}

	return userApiResponse.Result, nil
}

func readApiToken() (string, error) {
	dat, err := os.ReadFile(USER_DATA_FILE)
	if err != nil {
		fmt.Println("Error reading user data file.")
		return "", err
	}

	userData := UserData{}

	if err := json.Unmarshal([]byte(dat), &userData); err != nil {
		fmt.Println("Error unmarshaling user data.")
		return "", err
	}
	return userData.ApiToken, nil
}

func SaveUserDataToDisk(userData UserData) error {
	data, err := json.Marshal(userData)
	if err != nil {
		return errors.New("Error marshaling user data to byte string.")
	}
	if err := os.WriteFile(USER_DATA_FILE, data, 0666); err != nil {
		return errors.New("Error writing new data to disk")
	}
	return nil
}
func ReadUserDataFromDisk(verbose bool) (UserData, error) {
	dat, err := os.ReadFile(USER_DATA_FILE)
	if err != nil {
		return UserData{}, err
	}
	userData := UserData{}
	if err := json.Unmarshal(dat, &userData); err != nil {
		return UserData{}, err
	}
	return userData, nil
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

func CheckApiTokenExists() bool {
	_, err := readApiToken()
	if err != nil {
		fmt.Println("Can't find API token. Did you forget to call `bonusly config --token`?")
		return false
	}
	return true
}
