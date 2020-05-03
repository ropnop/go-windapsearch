package adschema

import (
	"encoding/json"
	"strconv"
)

type ADAttribute struct {
	CN string `json:"CN"`
	LdapDisplayName string `json:"Ldap-Display-Name"`
	AttributeId string `json:"Attribute-Id"`
	SystemIDGuid string `json:"System-Id-Guid"`
	Syntax AttributeSyntax `json:"Syntax"`
}

type syntax int
const (
	Boolean syntax = iota
	Enumeration
	Interval
	Object_Access_Point
	Object_DN_Binary
	Object_DS_DN
	Object_Presentation_Address
	Object_Replica_Link
	String_Generalized_Time
	String_IA5
	String_NT_Sec_Desc
	String_Numeric
	String_Object_Identifier
	String_Sid
	String_Teletex
	String_Unicode
)

type AttributeSyntax struct {
	Name string
	ConvertBytes ConvertBytes
}

type ConvertBytes func([]byte) interface{}


var BooleanSyntax = AttributeSyntax{
	Name:    "Boolean",
	ConvertBytes: ConvertBool,
}
func ConvertBool(b []byte) interface{} {
	v, err := strconv.ParseBool(string(b))
	if err != nil {
		return nil
	}
	return v
}

func (a *AttributeSyntax) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s{
	case "Boolean":
		a.Name = "Boolean"
		a.ConvertBytes = ConvertBool
	case "Interval":
		a.Name = "Interval"
		a.ConvertBytes = func(b []byte) interface{} { return nil }
	default:
		a.Name = "Unknown"
		a.ConvertBytes = func(b []byte) interface{} { return nil }
	}
	return nil
}

//var EnumerationSyntax = AttributeSyntax{
//	Name:    "Enumeration",
//	Convert: ParseEnumeration,
//}
//
//var IntervalSyntax = AttributeSyntax{
//	Name:    "Interval",
//	Convert: ParseInterval,
//}
//
//var ObjectAccessPoint = AttributeSyntax{
//	Name:    "Object(Access-Point)",
//	Convert: ParseObject,
//}
