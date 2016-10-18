package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/howeyc/gopass"
)

// String prompt.
func String(prompt string, args ...interface{}) string {
	fmt.Printf(prompt+": ", args...)
	reader := bufio.NewReader(os.Stdin)
	bytes, _, _ := reader.ReadLine()
	return string(bytes)
}

// StringRequired prompt.
func StringRequired(prompt string, args ...interface{}) string {
	var s string

	for {
		s = String(prompt, args...)
		if strings.Trim(s, " ") == "" {
			continue
		}

		break
	}

	return s
}

// Confirm continues prompting until the input is boolean-ish.
func Confirm(prompt string, args ...interface{}) bool {
	for {
		s := StringRequired(prompt, args...)

		switch s {
		case "yes", "y", "Y":
			return true
		case "no", "n", "N":
			return false
		default:
			continue
		}
	}
}

// Password prompt.
func Password(prompt string, args ...interface{}) string {
	fmt.Printf(prompt+": ", args...)
	password, _ := gopass.GetPasswd()
	s := string(password[0:])
	return s
}

// PasswordMasked is a password prompt with a mask.
func PasswordMasked(prompt string, args ...interface{}) string {
	fmt.Printf(prompt+": ", args...)
	password, _ := gopass.GetPasswdMasked()
	s := string(password[0:])
	return s
}
