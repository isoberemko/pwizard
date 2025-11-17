package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
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

func CryptoPasswordGenerate(charset string, length int) (string, error) {
	bytes, secureString := make([]byte, length), make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	charsetLength := len(charset)
	for i := range bytes {
		secureString[i] = charset[bytes[i]%byte(charsetLength)]
	}
	return string(secureString), nil
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
		if passwords[i].Password, err = CryptoPasswordGenerate(passwords[i].Charset, passwords[i].Length); err != nil {
			log.Println(err)
		}
		fmt.Println(
			passwords[i].Type, "\n",
			passwords[i].EnableSymbols, "\n",
			passwords[i].Length, "\n",
			passwords[i].Password)
		fmt.Println()
	}
}
