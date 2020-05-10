package adschema

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

var (
	_syntaxNameToValue = map[string]syntax{
		"Boolean":                      Boolean,
		"Enumeration":                  Enumeration,
		"Interval":                     Interval,
		"Object(Access-Point)":         Object_Access_Point,
		"Object(DN-Binary)":            Object_DN_Binary,
		"Object(DS-DN)":                Object_DS_DN,
		"Object(Presentation-Address)": Object_Presentation_Address,
		"Object(Replica-Link)":         Object_Replica_Link,
		"String(Generalized-Time)":     String_Generalized_Time,
		"String(IA5)":                  String_IA5,
		"String(NT-Sec-Desc)":          String_NT_Sec_Desc,
		"String(Numeric)":              String_Numeric,
		"String(Object-Identifier)":    String_Object_Identifier,
		"String(Sid)":                  String_Sid,
		"String(Teletex)":              String_Teletex,
		"String(Unicode)":              String_Unicode,
	}

	_syntaxValueToName = map[syntax]string{
		Boolean:                     "Boolean",
		Enumeration:                 "Enumeration",
		Interval:                    "Interval",
		Object_Access_Point:         "Object_Access_Point",
		Object_DN_Binary:            "Object_DN_Binary",
		Object_DS_DN:                "Object_DS_DN",
		Object_Presentation_Address: "Object_Presentation_Address",
		Object_Replica_Link:         "Object_Replica_Link",
		String_Generalized_Time:     "String_Generalized_Time",
		String_IA5:                  "String_IA5",
		String_NT_Sec_Desc:          "String_NT_Sec_Desc",
		String_Numeric:              "String_Numeric",
		String_Object_Identifier:    "String_Object_Identifier",
		String_Sid:                  "String_Sid",
		String_Teletex:              "String_Teletex",
		String_Unicode:              "String_Unicode",
	}
)
