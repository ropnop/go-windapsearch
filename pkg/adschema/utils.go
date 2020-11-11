package adschema

import (
	"encoding/binary"
	"fmt"
	"github.com/bwmarrin/go-objectsid"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

var FunctionalityLevelsMapping = map[string]string{
	"0": "2000",
	"1": "2003 Interim",
	"2": "2003",
	"3": "2008",
	"4": "2008 R2",
	"5": "2012",
	"6": "2012 R2",
	"7": "2016",
	"":  "Unknown",
}

var NTFileTimeRegex *regexp.Regexp
var ADLdapTimeRegex *regexp.Regexp

func init() {
	NTFileTimeRegex = regexp.MustCompile(`^[0-9]{18}$`)
	ADLdapTimeRegex = regexp.MustCompile(`^[0-9]{14}\.[0-9]Z$`)
}

func WindowsGuidFromBytes(b []byte) (string, error) {
	if len(b) != 16 {
		return "", fmt.Errorf("GUID must be 16 bytes")
	}
	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		binary.LittleEndian.Uint32(b[:4]),
		binary.LittleEndian.Uint16(b[4:6]),
		binary.LittleEndian.Uint16(b[6:8]),
		b[8:10],
		b[10:]), nil
}

func WindowsSIDFromBytes(b []byte) (string, error) {
	if len(b) < 12 {
		return "", fmt.Errorf("windows SID seems too short")
	}
	sid := objectsid.Decode(b)
	return sid.String(), nil
}

func NTFileTimeToTimestamp(s string) (timestamp time.Time, err error) {
	ticks, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}

	secs := (int64)((ticks / 10000000) - 11644473600)
	nsecs := (int64)((ticks % 10000000) * 100)

	return time.Unix(secs, nsecs), nil
}

func ADLdapTimeToTimestamp(s string) (timestamp time.Time, err error) {
	s = strings.TrimSuffix(s, ".0Z")
	return time.Parse("20060102150405", s)
}

// credit: https://golang.org/src/syscall/syscall_windows.go
func UTF16ToString(s []uint16) string {
	for i, v := range s {
		if v == 0 {
			s = s[0:i]
			break
		}
	}
	return string(utf16.Decode(s))
}
