package adschema

import (
	"encoding/base64"
	"encoding/json"
	"gopkg.in/ldap.v3"
	"unicode/utf8"
)

type LDAPAttribute ldap.EntryAttribute
type LDAPEntryJSON map[string]interface{}

type ADAttribute struct {
	*ldap.EntryAttribute
}

//func (e *ADEntry) MarshalJSON() ([]byte, error) {
//	jEntry := make(LDAPEntryJSON)
//	for _, attribute := range e.Attributes {
//		if len(attribute.Values) == 1 {
//			jEntry[attribute.Name] = HandleLDAPBytes(attribute.Name, attribute.ByteValues[0])
//		} else {
//			var vals []interface{}
//			for _, val := range attribute.ByteValues {
//				vals = append(vals, HandleLDAPBytes(attribute.Name, val))
//			}
//			jEntry[attribute.Name] = vals
//		}
//	}
//	return json.Marshal(jEntry)
//}

func (e *ADEntry) MarshalJSON() ([]byte, error) {
	jEntry := make(map[string]*ADAttribute)
	for _, attribute := range e.Attributes {
		jEntry[attribute.Name] = &ADAttribute{attribute}
	}
	return json.Marshal(jEntry)
}

func (e *ADAttribute) MarshalJSON() ([]byte, error) {
	// Look up syntax for attribute name
	info, ok := AttributeMap[e.Name]
	if !ok {
		return marshalUnknownAttribute(e)
	}
	convert, ok := SyntaxFunctions[info.Syntax]
	if !ok {
		convert = DefaultPrint
	}
	var vals []interface{}
	for _, v := range e.ByteValues {
		i, err := convert(e.Name, v)
		if err != nil {
			return nil, err
		}
		vals = append(vals, i)
	}
	if info.IsSingleValue && len(vals) == 1 {
		return json.Marshal(vals[0])
	}
	return json.Marshal(vals)

}

func marshalUnknownAttribute(e *ADAttribute) ([]byte, error) {
	var vals []string
	for _, val := range e.ByteValues {
		vals = append(vals, printable(val))
	}
	info, ok := AttributeMap[e.Name]
	if ok {
		if info.IsSingleValue && len(vals) == 1 {
			return json.Marshal(vals[0])
		}
	}

	return json.Marshal(vals)
}

func printable(b []byte) string {
	if utf8.Valid(b) {
		return string(b)
	}
	return base64.StdEncoding.EncodeToString(b)
}
