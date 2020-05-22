package modules

import (
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/spf13/pflag"
)

type GPOsModule struct{}

func init() {
	AllModules = append(AllModules, new(GPOsModule))
}

func (g GPOsModule) Name() string {
	return "gpos"
}

func (g GPOsModule) Description() string {
	return "Enumerate Group Policy Objects"
}

func (g *GPOsModule) FlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("gpos", pflag.ExitOnError)
	return flags
}

func (g GPOsModule) DefaultAttrs() []string {
	return []string{"displayName", "gPCFileSysPath"}
}

func (g GPOsModule) Filter() string {
	return "(objectClass=groupPolicyContainer)"
}

func (g *GPOsModule) Run(session *ldapsession.LDAPSession, attrs []string) error {
	sr := session.MakeSimpleSearchRequest(g.Filter(), attrs)
	return session.ExecuteSearchRequest(sr)
}
