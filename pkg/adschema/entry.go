package adschema

import (
	"encoding/base64"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ADEntry struct {
	*ldap.Entry
}

func (e *ADEntry) String() string {
	return e.DN
}

func (e *ADEntry) LDAPFormat() string {
	var sb strings.Builder
	if e.DN != "" {
		sb.WriteString(fmt.Sprintf("dn: %s\n", e.DN))
	}
	for _, attribute := range e.Attributes {
		for _, value := range attribute.ByteValues {
			//valueString := HandleLDAPBytes(attribute.Name, value)
			sb.WriteString(fmt.Sprintf("%s: %v\n", attribute.Name, printable(value)))
		}
	}
	return sb.String()
}

// HandleLDAPBytes takes a byte slice from a raw attribute value and returns either a UTF8 string (if it's a string),
// or GUID or timestamp
func HandleLDAPBytes(name string, b []byte) interface{} {
	if name == "objectGUID" {
		g, err := WindowsGuidFromBytes(b)
		if err != nil {
			return b
		}
		return g
	}
	if name == "objectSid" {
		s, err := WindowsSIDFromBytes(b)
		if err != nil {
			return b
		}
		return s
	}

	if name == "domainFunctionality" {
		return FunctionalityLevelsMapping[string(b)]
	}
	if name == "forestFunctionality" {
		return FunctionalityLevelsMapping[string(b)]
	}
	if name == "domainControllerFunctionality" {
		return FunctionalityLevelsMapping[string(b)]
	}

	if utf8.Valid(b) {
		s := string(b)
		if s == "9223372036854775807" { //max int64 size
			return 0 //basically a no-value (e.g. never expires)
		}
		if NTFileTimeRegex.Match(b) {
			timeStamp, err := NTFileTimeToTimestamp(s)
			if err != nil {
				return s
			}
			return timeStamp
		}
		if ADLdapTimeRegex.Match(b) {
			timeStamp, err := ADLdapTimeToTimestamp(s)
			if err != nil {
				return s
			}
			return timeStamp
		}
		if i, err := strconv.Atoi(s); err == nil {
			return i
		}
		return s
	}
	return base64.StdEncoding.EncodeToString(b)
}
