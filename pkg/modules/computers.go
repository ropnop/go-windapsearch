package modules

import (
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/spf13/pflag"
)

type ComputersModule struct{}

func init() {
	AllModules = append(AllModules, new(ComputersModule))
}

func (c ComputersModule) Name() string {
	return "computers"
}

func (c ComputersModule) Description() string {
	return "Enumerate AD Computers"
}

func (c ComputersModule) FlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("computers-module", pflag.ExitOnError)
	return flags
}

func (c ComputersModule) DefaultAttrs() []string {
	return []string{"cn", "dNSHostName", "operatingSystem", "operatingSystemVersion", "operatingSystemServicePack"}
}

func (c ComputersModule) Run(session *ldapsession.LDAPSession, attrs []string) error {
	filter := "(objectClass=Computer)"
	searchReq := session.MakeSimpleSearchRequest(filter, attrs)
	return session.ExecuteSearchRequest(searchReq)
}
