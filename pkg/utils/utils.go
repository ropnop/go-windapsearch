package utils

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

func SecurePrompt(message string) (response string, err error) {
	fmt.Printf("%s: ", message)
	securebytes, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return
	}
	fmt.Println()
	return string(securebytes), nil
}
