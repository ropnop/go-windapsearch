package adschema

import (
	"encoding/json"
	"gopkg.in/ldap.v3"
)

type LDAPAttribute ldap.EntryAttribute
type LDAPEntryJSON map[string]interface{}

func (e *ADEntry) MarshalJSON() ([]byte, error) {
	jEntry := make(LDAPEntryJSON)
	for _, attribute := range e.Attributes {
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
	return json.Marshal(jEntry)
}

