package utils

import (
	"encoding/json"
	"github.com/ropnop/go-windapsearch/pkg/adschema"
	"gopkg.in/ldap.v3"
)



type LDAPEntryJSON map[string]interface{}

func SearchResultToJSON(result *ldap.SearchResult) (jResponse []byte, err error) {
	var ldapResponsesJSON []LDAPEntryJSON
	for _, entry := range result.Entries {
		jEntry := make(LDAPEntryJSON)
		for _, attribute := range entry.Attributes {
			if len(attribute.Values) == 1 {
				jEntry[attribute.Name] = adschema.HandleLDAPBytes(attribute.Name, attribute.ByteValues[0])
			} else {
				var vals []interface{}
				for _, val := range attribute.ByteValues {
					vals = append(vals, adschema.HandleLDAPBytes(attribute.Name, val))
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

func EntryToJSON(entry *ldap.Entry) (jResponse []byte, err error) {
	jEntry := make(LDAPEntryJSON)
	for _, attribute := range entry.Attributes {
		if len(attribute.Values) == 1 {
			jEntry[attribute.Name] = adschema.HandleLDAPBytes(attribute.Name, attribute.ByteValues[0])
		} else {
			var vals []interface{}
			for _, val := range attribute.ByteValues {
				vals = append(vals, adschema.HandleLDAPBytes(attribute.Name, val))
			}
			jEntry[attribute.Name] = vals
		}
	}
	return json.Marshal(jEntry)
}


