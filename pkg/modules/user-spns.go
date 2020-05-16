package modules

import (
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/spf13/pflag"
)

type UserSPNsModule struct{}

func init() {
	AllModules = append(AllModules, new(UserSPNsModule))
}

func (u UserSPNsModule) Name() string {
	return "user-spns"
}

func (u UserSPNsModule) Description() string {
	return "Enumerate all users objects with Service Principal Names (for kerberoasting)"
}

func (u UserSPNsModule) FlagSet() *pflag.FlagSet {
	return pflag.NewFlagSet("user-spns", pflag.ExitOnError)
}

func (u UserSPNsModule) DefaultAttrs() []string {
	return []string{"cn", "servicePrincipalName"}
}

func (u UserSPNsModule) Filter() string {
	return "(&(&(servicePrincipalName=*)(UserAccountControl:1.2.840.113556.1.4.803:=512))(!(UserAccountControl:1.2.840.113556.1.4.803:=2)))"

}

func (u UserSPNsModule) Run(session *ldapsession.LDAPSession, attrs []string) error {
	sr := session.MakeSimpleSearchRequest(u.Filter(), attrs)
	return session.ExecuteSearchRequest(sr)
}

