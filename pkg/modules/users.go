package modules

var GetAllUsers = WindapModule{
	Name:         "users",
	Description:  "Dump all user objects",
	Filter:       "(objectCategory=user)",
	DefaultAttrs: []string{"cn", "userPrincipalName"},
}



