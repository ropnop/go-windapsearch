package enums

import (
	uac "github.com/audibleblink/msldapuac"
)

type ConvertEnum func(int64) interface{}

var EnumFuncs = map[string]ConvertEnum{
	"sAMAccountType": func(i int64) interface{} {
		val, ok := SamAccountTypeEnum[i]
		if !ok {
			return i
		}
		return val
	},
	"userAccountControl": ConvertUAC,
	"groupType":          ConvertGroupType,
}

// SAM-Account-Type
// https://docs.microsoft.com/en-us/windows/win32/adschema/a-samaccounttype
var SamAccountTypeEnum = map[int64]string{
	0x0:        "SAM_DOMAIN_OBJECT",
	0x10000000: "SAM_GROUP_OBJECT",
	0x10000001: "SAM_NON_SECURITY_GROUP_OBJECT",
	0x20000000: "SAM_ALIAS_OBJECT",
	0x20000001: "SAM_NON_SECURITY_ALIAS_OBJECT",
	0x30000000: "SAM_USER_OBJECT",
	0x30000001: "SAM_MACHINE_ACCOUNT",
	0x30000002: "SAM_TRUST_ACCOUNT",
	0x40000000: "SAM_APP_BASIC_GROUP",
	0x40000001: "SAM_APP_QUERY_GROUP",
	0x7fffffff: "SAM_ACCOUNT_TYPE_MAX",
}

func ConvertUAC(i int64) interface{} {
	flags, err := uac.ParseUAC(i)
	if err != nil {
		return i
	}
	return flags
}
