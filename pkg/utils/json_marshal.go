package utils

import (
	"encoding/base64"
	"encoding/json"
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"gopkg.in/ldap.v3"
	"strconv"
	"unicode/utf8"
)



type LDAPEntryJSON map[string]interface{}

func SearchResultToJSON(result *ldap.SearchResult) (jResponse []byte, err error) {
	var ldapResponsesJSON []LDAPEntryJSON
	for _, entry := range result.Entries {
		jEntry := make(LDAPEntryJSON)
		for _, attribute := range entry.Attributes {
			if len(attribute.Values) == 1 {
				jEntry[attribute.Name] = HandleLDAPBytes(attribute.Name, attribute.ByteValues[0])
			} else {
				var vals []interface{}
				for _, val := range attribute.ByteValues {
					vals = append(vals, HandleLDAPBytes(attribute.Name, val))
				}
				jEntry[attribute.Name] = vals
			}
		}
		ldapResponsesJSON = append(ldapResponsesJSON, jEntry)
	}
	if len(ldapResponsesJSON) == 1 {
		return json.Marshal(ldapResponsesJSON[0])
	}
	return json.Marshal(ldapResponsesJSON)
}

// HandleLDAPBytes takes a byte slice from a raw attribute value and returns either a UTF8 string (if it's a string),
// or GUID or timestampgit s
func HandleLDAPBytes(name string, b []byte) interface{} {
	if name == "objectGUID" {
		g, err := WindowsGuidFromBytes(b); if err != nil {
			return b
		}
		return g
	}
	if name == "objectSid" {
		s, err := WindowsSIDFromBytes(b); if err != nil {
			return b
		}
		return s
	}

	if name == "domainFunctionality" {
		return ldapsession.FunctionalityLevelsMapping[string(b)]
	}
	if name == "forestFunctionality" {
		return ldapsession.FunctionalityLevelsMapping[string(b)]
	}
	if name == "domainControllerFunctionality" {
		return ldapsession.FunctionalityLevelsMapping[string(b)]
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


