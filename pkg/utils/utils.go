package utils

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"
)

func SecurePrompt(message string) (response string, err error) {
	fmt.Fprintf(os.Stderr, "%s: ", message)
	securebytes, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return
	}
	fmt.Fprint(os.Stderr, "\n")
	return string(securebytes), nil
}
