package adschema

import (
	"encoding/binary"
	"fmt"
	"github.com/bwmarrin/go-objectsid"
	"github.com/ropnop/go-windapsearch/pkg/adschema/enums"
	"strconv"
	"time"
)

// Unique syntaxes
//"Boolean",
//"Enumeration",
//"Interval",
//"Object(Access-Point)",
//"Object(DN-Binary)",
//"Object(DS-DN)",
//"Object(Presentation-Address)",
//"Object(Replica-Link)",
//"String(Generalized-Time)",
//"String(IA5)",
//"String(NT-Sec-Desc)",
//"String(Numeric)",
//"String(Object-Identifier)",
//"String(Sid)",
//"String(Teletex)",
//"String(Unicode)"

type ConvertBytes = func(string, []byte) (interface{}, error)

// these are custon functions for converting LDAP bytes to a more readable form
// I don't account for all of them - if the syntax isn't listed they just default to a printable string
var SyntaxFunctions = map[string]ConvertBytes{
	"Boolean":                  ConvertBool,
	"String(Generalized-Time)": ConvertGeneralizedTime,
	"Interval":                 ConvertInterval,
	"String(Sid)":              ConvertSid,
	"Object(Replica-Link)":     ConvertObjectReplicaLink,
	"Enumeration":              ConvertEnumeration,
}

func DefaultPrint(name string, b []byte) (interface{}, error) {
	return printable(b), nil
}

func ConvertBool(name string, b []byte) (interface{}, error) {
	return strconv.ParseBool(string(b))
}

func ConvertGeneralizedTime(name string, b []byte) (interface{}, error) {
	// https://docs.microsoft.com/en-us/windows/win32/adschema/s-string-generalized-time
	timestamp := string(b)
	return time.Parse("20060102150405.0Z0700", timestamp)
}

// these attrbitures are longs which represent "number of 100 nanosecond intervals since January 1, 1601 (UTC)"
var NTFiletimeAttributes = map[string]bool{
	"accountExpires":     true,
	"badPasswordTime":    true,
	"lastLogoff":         true,
	"lastLogon":          true,
	"lastLogonTimestamp": true,
	"lastSetTime":        true,
	"lockoutTime":        true,
	"pwdLastSet":         true,
}

func ConvertInterval(name string, b []byte) (interface{}, error) {
	// https://docs.microsoft.com/en-us/windows/win32/adschema/s-interval
	timestamp := string(b)
	if timestamp == "9223372036854775807" || timestamp == "0" { // indicates a "never", I chose to represent this as a 0
		return "0", nil
	}

	if _, ok := NTFiletimeAttributes[name]; ok {
		return NTFileTimeToTimestamp(timestamp)
	}
	return timestamp, nil
}

func ConvertSid(name string, b []byte) (interface{}, error) {
	if len(b) < 12 {
		return "", fmt.Errorf("windows SID seems too short")
	}
	sid := objectsid.Decode(b)
	return sid.String(), nil
}

func ConvertObjectReplicaLink(name string, b []byte) (interface{}, error) {
	if len(b) != 16 {
		return printable(b), nil
	}
	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		binary.LittleEndian.Uint32(b[:4]),
		binary.LittleEndian.Uint16(b[4:6]),
		binary.LittleEndian.Uint16(b[6:8]),
		b[8:10],
		b[10:]), nil
}

func ConvertEnumeration(name string, b []byte) (interface{}, error) {
	// https://docs.microsoft.com/en-us/windows/win32/adschema/s-enumeration
	// Active Directory treats this as an integer.
	val, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return 0, err
	}
	if _, ok := enums.EnumFuncs[name]; ok {
		return enums.EnumFuncs[name](val), nil
	}
	return val, nil
}
