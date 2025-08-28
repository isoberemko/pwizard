package main

import (
	"fmt"
	"math/rand"
)

const nums = "1234567890"
const symbols = "!@#$%^&*()-_=+/?[]{}`~"
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func aggregateString(substrings ...string) (aggregatedString string) {
	for _, substring := range substrings {
		aggregatedString += substring
	}
	return aggregatedString
}

func passwordGenerate(charset string, length int) (password string) {
	charsetLen := len(charset)
	for range length {
		password += string(charset[rand.Intn(charsetLen)])
	}
	return password
}

func main() {
	var (
		charset  string
		password struct {
			length        int
			enableSymbols bool
		}
	)

	password.length = 18

	if password.enableSymbols {
		charset = aggregateString(nums, symbols, letters)
	} else {
		charset = aggregateString(nums, letters)
	}

	fmt.Println(passwordGenerate(charset, password.length))

}
