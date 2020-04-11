package utils

import (
	"fmt"
	"gopkg.in/ldap.v3"
	"io"
)

func WriteSearchResults(result *ldap.SearchResult, writer io.Writer) {
	if result == nil {
		io.WriteString(writer, "[-] No results")
		return
	}
	for _, entry := range result.Entries {
		for _, attribute := range entry.Attributes {
			for _, value := range attribute.Values {
				io.WriteString(writer, fmt.Sprintf("%v: %v\n", attribute.Name, value))
			}
		}
		io.WriteString(writer, "\n")

	}
	return
}
