package utils

import (
	"fmt"
	"github.com/tcnksm/go-input"
	"gopkg.in/ldap.v3"
	"os"
)

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
