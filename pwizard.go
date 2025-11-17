package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	nums    = "1234567890"
	symbols = "!@#$%^&*()-_=+/?[]{}`~"
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Passwords []Password

type Password struct {
	Length        int    `json:"length"`
	EnableSymbols bool   `json:"symbols"`
	Type          string `json:"name"`
	Charset       string
	Password      string
}

func (pws *Passwords) PasswordsConfigure(configPath string) error {
	configJson, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(configJson, pws)
	if err != nil {
		return err
	}
	return nil
}

func (pw *Password) PasswordValidate() error {
	if pw.Length < 4 {
		return errors.New("Password length can not be less then 4 symbols")
	}
	return nil
}

func AggregateString(substrings ...string) (aggregatedString string) {
	return strings.Join(substrings, "")
}

func CryptoPasswordGenerate(charset string, length int) (string, error) {
	randomBytes, secureString := make([]byte, length), make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	charsetLength := len(charset)
	for i := range randomBytes {
		secureString[i] = charset[randomBytes[i]%byte(charsetLength)]
	}
	return string(secureString), nil
}

func main() {
	configPath := flag.String("config", "/etc/pwizard/passwords.json", "config file path")
	flag.Parse()

	var passwords Passwords

	if err := passwords.PasswordsConfigure(*configPath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i := range passwords {
		if err := passwords[i].PasswordValidate(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	for i := range passwords {
		var err error
		if passwords[i].EnableSymbols {
			passwords[i].Charset = AggregateString(nums, symbols, letters)
		} else {
			passwords[i].Charset = AggregateString(nums, letters)
		}
		if passwords[i].Password, err = CryptoPasswordGenerate(passwords[i].Charset, passwords[i].Length); err != nil {
			fmt.Println("Can not generate password for", passwords[i].Type, ":", err)
			continue
		}
		fmt.Printf("%v: %v\n\n", passwords[i].Type, passwords[i].Password)
	}
}
