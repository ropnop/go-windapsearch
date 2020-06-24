package modules

import (
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/pflag"
)

type FunctionalityModule struct{}

func init() {
	AllModules = append(AllModules, new(FunctionalityModule))
}

func (FunctionalityModule) Name() string {
	return "metadata"
}

func (FunctionalityModule) Description() string {
	return "Print LDAP server metadata"
}

func (FunctionalityModule) FlagSet() *pflag.FlagSet {
	return pflag.NewFlagSet("metadata", pflag.ExitOnError)
}

func (FunctionalityModule) DefaultAttrs() []string {
	return []string{
		"defaultNamingContext",
		"domainFunctionality",
		"forestFunctionality",
		"domainControllerFunctionality",
		"dnsHostName",
	}
}

func (FunctionalityModule) Run(session *ldapsession.LDAPSession, attrs []string) error {
	sr := ldap.NewSearchRequest(
		"",
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		attrs,
		nil)
	//res, err := session.LConn.Search(sr)
	res, err := session.GetSearchResults(sr)
	if err != nil {
		return err
	}
	session.ManualWriteSearchResultsToChan(res)
	return nil
}
