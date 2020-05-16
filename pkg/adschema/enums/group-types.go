package enums

import (
	"fmt"
	"github.com/audibleblink/bamflags"
	"os"
)

// https://docs.microsoft.com/en-us/windows/win32/adschema/a-grouptype#remarks

//1 (0x00000001)	Specifies a group that is created by the system.
//2 (0x00000002)	Specifies a group with global scope.
//4 (0x00000004)	Specifies a group with domain local scope.
//8 (0x00000008)	Specifies a group with universal scope.
//16 (0x00000010)	Specifies an APP_BASIC group for Windows Server Authorization Manager.
//32 (0x00000020)	Specifies an APP_QUERY group for Windows Server Authorization Manager.
//2147483648 (0x80000000)	Specifies a security group. If this flag is not set, then the group is a distribution group.

const (
	SystemGroup = 1 << iota
	GlobalScope
	DomainScope
	UniversalScope
	APPBasic
	APPQuery
)

var GroupTypeMap = map[int]string{
	SystemGroup: "Created by system",
	GlobalScope: "Global Scope",
	DomainScope: "Domain Local Scope",
	UniversalScope: "Universal Scope",
	APPBasic: "APP_BASIC group for Windows Server Authorization Manager",
	APPQuery: "APP_QUERY group for Windows Server Authorization Manager",
}

func ConvertGroupType(groupType int64) interface{} {
	fmt.Fprintf(os.Stderr, "groupType: %d\n", groupType)
	values, err := bamflags.ParseInt(groupType)
	if err != nil {
		return []string{fmt.Sprintf("%s", groupType)}
	}
	var flags []string
	for _, value := range values {
		if propName, ok := GroupTypeMap[value]; ok {
			flags = append(flags, propName)
		}
	}
	return flags
}

