package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

const (
	digits       = "1234567890"
	lettersLower = "abcdefghijklmnopqrstuvwxyz"
	lettersUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letters      = lettersLower + lettersUpper

	symbolsSafe      = "!@%*_-."
	symbolsExtended  = symbolsSafe + "#$^&()=+[]{}|;:,<>?~"
	symbolsDangerous = symbolsExtended + "`\"'\\/"
)

type Passwords []Password

type Password struct {
	Length                 *int    `json:"length"`
	Name                   *string `json:"name"`
	EnableSymbolsSafe      *bool   `json:"symbols"`
	EnableSymbolsExtended  *bool   `json:"symbols_extended"`
	EnableSymbolsDangerous *bool   `json:"symbols_dangerous"`
	CustomSymbolsSet       *string `json:"custom_symbols_set"`
	Digits                 *bool   `json:"digits"`
	Letters                *bool   `json:"letters"`
	UppercaseOnly          *bool   `json:"uppercase_only"`
	LowercaseOnly          *bool   `json:"lowercase_only"`
	Charset                string
	Password               string
}

func ErrorCritical(err error) {
	fmt.Println(err)
	os.Exit(1)
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

func (pw *Password) PasswordConfigNormalizeDigits() {
	if pw.Digits == nil {
		pw.Digits = new(bool)
		*pw.Digits = true
	}
}

func (pw *Password) PasswordConfigNormalizeLetters() {
	if pw.LowercaseOnly == nil {
		pw.LowercaseOnly = new(bool)
	}
	if pw.UppercaseOnly == nil {
		pw.UppercaseOnly = new(bool)
	}
	if pw.Letters == nil {
		pw.Letters = new(bool)
		*pw.Letters = true
	}
}

func (pw *Password) PasswordConfigNormalizeSymbols() {
	if pw.EnableSymbolsSafe == nil {
		pw.EnableSymbolsSafe = new(bool)
	}
	if pw.EnableSymbolsExtended == nil {
		pw.EnableSymbolsExtended = new(bool)
	}
	if pw.EnableSymbolsDangerous == nil {
		pw.EnableSymbolsDangerous = new(bool)
	}
}

func (pw *Password) PasswordConfigNormalize() {
	pw.PasswordConfigNormalizeDigits()
	pw.PasswordConfigNormalizeLetters()
	pw.PasswordConfigNormalizeSymbols()
}

func (pw *Password) PasswordConfigValidateName() error {
	switch {
	case pw.Name == nil:
		return errors.New("password configuration name is not specified")
	case len(*pw.Name) == 0:
		return errors.New("password configuration name must not be empty")
	}
	return nil
}

func (pw *Password) PasswordConfigValidateLength() error {
	switch {
	case pw.Length == nil:
		return errors.New("password length is not specified")
	case *pw.Length < 4:
		return errors.New("password length must not be less than 4 characters")
	}
	return nil
}

func (pw *Password) PasswordConfigValidateLetters() error {
	switch {
	case *pw.LowercaseOnly && *pw.UppercaseOnly:
		return errors.New("lowercase-only and uppercase-only modes cannot be enabled at the same time")
	case !*pw.Letters && *pw.LowercaseOnly:
		return errors.New("lowercase-only mode cannot be enabled when letters are disabled")
	case !*pw.Letters && *pw.UppercaseOnly:
		return errors.New("uppercase-only mode cannot be enabled when letters are disabled")
	}
	return nil
}

func (pw *Password) PasswordConfigValidateSymbols() error {
	if pw.CustomSymbolsSet != nil {
		switch {
		case *pw.EnableSymbolsDangerous:
			return errors.New("custom symbol set cannot be combined with the dangerous symbol set")
		case *pw.EnableSymbolsExtended:
			return errors.New("custom symbol set cannot be combined with the extended symbol set")
		case *pw.EnableSymbolsSafe:
			return errors.New("custom symbol set cannot be combined with the safe symbol set")
		}
	}
	return nil
}

func (pw *Password) PasswordConfigValidate() {
	if err := pw.PasswordConfigValidateName(); err != nil {
		ErrorCritical(err)
	}
	if err := pw.PasswordConfigValidateLength(); err != nil {
		ErrorCritical(err)
	}
	if err := pw.PasswordConfigValidateLetters(); err != nil {
		ErrorCritical(err)
	}
	if err := pw.PasswordConfigValidateSymbols(); err != nil {
		ErrorCritical(err)
	}
}

func (pw *Password) CreateCharset() error {
	if *pw.Digits {
		pw.Charset += digits
	}
	if *pw.Letters {
		switch {
		case *pw.LowercaseOnly:
			pw.Charset += lettersLower
		case *pw.UppercaseOnly:
			pw.Charset += lettersUpper
		default:
			pw.Charset += letters
		}
	}
	switch {
	case pw.CustomSymbolsSet != nil:
		pw.Charset += *pw.CustomSymbolsSet
	case *pw.EnableSymbolsDangerous:
		pw.Charset += symbolsDangerous
	case *pw.EnableSymbolsExtended:
		pw.Charset += symbolsExtended
	case *pw.EnableSymbolsSafe:
		pw.Charset += symbolsSafe
	}
	if len(pw.Charset) == 0 {
		return errors.New("character set is empty")
	}
	return nil
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
		ErrorCritical(err)
	}

	for i := range passwords {
		var err error

		passwords[i].PasswordConfigNormalize()
		passwords[i].PasswordConfigValidate()

		if err = passwords[i].CreateCharset(); err != nil {
			ErrorCritical(err)
		}
		if passwords[i].Password, err = CryptoPasswordGenerate(passwords[i].Charset, *passwords[i].Length); err != nil {
			fmt.Println("Cannot generate password for", *passwords[i].Name, ":", err)
			continue
		}
		fmt.Printf("%v: %v\n\n", *passwords[i].Name, passwords[i].Password)
	}
}
