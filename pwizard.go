package main

import (
	"fmt"
	"math/rand"
)

const (
	nums    = "1234567890"
	symbols = "!@#$%^&*()-_=+/?[]{}`~"
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Password struct {
	Length        int
	EnableSymbols bool
	Type          string
	Charset       string
	Password      string
}

func AddPasswordConfiguration(Length int, EnableSymbols bool, Type string) Password {
	return Password{
		Length:        Length,
		EnableSymbols: EnableSymbols,
		Type:          Type,
	}
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

	passwords = append(passwords, AddPasswordConfiguration(8, false, "Test8"))
	passwords = append(passwords, AddPasswordConfiguration(12, false, "Test12"))
	passwords = append(passwords, AddPasswordConfiguration(18, true, "Complex"))

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
