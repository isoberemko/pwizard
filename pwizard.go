package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

const (
	nums    = "1234567890"
	symbols = "!@#$%^&*()-_=+/?[]{}`~"
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Password struct {
	Length        int    `json:"length"`
	EnableSymbols bool   `json:"symbols"`
	Type          string `json:"name"`
	Charset       string
	Password      string
}

func AggregateString(substrings ...string) (aggregatedString string) {
	for _, substring := range substrings {
		aggregatedString += substring
	}
	return aggregatedString
}

func PasswordGenerate(charset string, length int) (password string) {
	charsetLen := len(charset)
	for range length {
		password += string(charset[rand.Intn(charsetLen)])
	}
	return password
}

func main() {
	passwords := []Password{}

	configJson, err := os.ReadFile("passwords.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configJson, &passwords)
	if err != nil {
		panic(err)
	}

	for i := range passwords {
		if passwords[i].EnableSymbols {
			passwords[i].Charset = AggregateString(nums, symbols, letters)
		} else {
			passwords[i].Charset = AggregateString(nums, letters)
		}
		passwords[i].Password = PasswordGenerate(passwords[i].Charset, passwords[i].Length)
		fmt.Println(
			passwords[i].Type, "\n",
			passwords[i].EnableSymbols, "\n",
			passwords[i].Length, "\n",
			passwords[i].Password)
		fmt.Println()
	}
}
