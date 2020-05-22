package utils

import (
	"fmt"
	"github.com/tcnksm/go-input"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/ldap.v3"
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

func ChooseDN(results *ldap.SearchResult) (dn string, err error) {
	var options []string
	for _, result := range results.Entries {
		options = append(options, result.DN)
	}
	if len(options) == 0 {
		return "", fmt.Errorf("no results")
	}
	if len(options) == 1 {
		return options[0], nil
	}
	ui := &input.UI{
		Writer: os.Stderr,
		Reader: os.Stdin,
	}
	query := "What DN do you want to use?"
	return ui.Select(query, options, &input.Options{})
}
