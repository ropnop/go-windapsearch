package modules

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/spf13/pflag"
)

type CustomSearch struct {
	CustomFilter string
	CustomBaseDN string
}

func init() {
	AllModules = append(AllModules, new(CustomSearch))
}

func (c *CustomSearch) Name() string {
	return "custom"
}

func (c *CustomSearch) Description() string {
	return "Run a custom LDAP syntax filter"
}

func (c *CustomSearch) Filter() string {
	return c.CustomFilter
}

func (c *CustomSearch) FlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("custom", pflag.ExitOnError)
	flags.StringVar(&c.CustomFilter, "filter", "", "LDAP syntax filter")
	flags.StringVar(&c.CustomBaseDN, "base", "", "Custom base DN to search from")
	return flags
}

func (c *CustomSearch) DefaultAttrs() []string {
	return []string{"*"}
}

func (c *CustomSearch) Run(lSession *ldapsession.LDAPSession, attrs []string) error {
	if c.Filter() == "" {
		return fmt.Errorf("must provide a filter to run")
	}
	var searchReq *ldap.SearchRequest
	if c.CustomBaseDN != "" {
		searchReq = lSession.MakeSearchRequestWithDN(c.CustomBaseDN, c.Filter(), attrs)
	} else {
		searchReq = lSession.MakeSimpleSearchRequest(c.Filter(), attrs)
	}
	return lSession.ExecuteSearchRequest(searchReq)
}
