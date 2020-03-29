package utils

// All this code was copied from the Windows syscall code here: https://golang.org/src/syscall/types_windows.go?s=10082:10127#L354
// I copied it here so we can build on non Windows OS

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var NTFileTimeRegex *regexp.Regexp
var ADLdapTimeRegex *regexp.Regexp

func init() {
	NTFileTimeRegex = regexp.MustCompile(`^[0-9]{18}$`)
	ADLdapTimeRegex = regexp.MustCompile(`^[0-9]{14}\.[0-9]Z$`)
}


func NTFileTimeToTimestamp(s string) (timestamp time.Time, err error) {
	ticks, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}

	secs := (int64)((ticks/10000000) -  11644473600)
	nsecs := (int64)((ticks % 10000000) * 100)

	return time.Unix(secs, nsecs), nil
}

func ADLdapTimeToTimestamp(s string) (timestamp time.Time, err error) {
	s = strings.TrimSuffix(s, ".0Z")
	return time.Parse("20060102150405", s)
}
