package modules

import (
	"fmt"
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/ropnop/go-windapsearch/pkg/utils"
	"github.com/spf13/pflag"
	"gopkg.in/ldap.v3"
)

type SearchModule struct {
	SearchTerm string
	AllResults bool
}

func init() {
	AllModules = append(AllModules, new(SearchModule))
}

func (s SearchModule) Name() string {
	return "search"
}

func (s SearchModule) Description() string {
	return "Perform an ANR Search and return the results"
}

func (s *SearchModule) FlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("search-mdoule", pflag.ExitOnError)
	flags.StringVarP(&s.SearchTerm, "search", "s", "", "Search term")
	flags.BoolVar(&s.AllResults, "all", false, "Output attrs for all matching search results")
	return flags
}

func (s SearchModule) DefaultAttrs() []string {
	return []string{"*"}
}

func (s *SearchModule) SearchFilter() string {
	return fmt.Sprintf("(%s)", utils.CreateANRSearch(s.SearchTerm))
}

func (s *SearchModule) Run(session *ldapsession.LDAPSession, attrs []string) error {
	if s.SearchTerm == "" {
		return fmt.Errorf("must include a search term")
	}
	if s.AllResults {
		sr := session.MakeSimpleSearchRequest(s.SearchFilter(), attrs)
		return session.ExecuteSearchRequest(sr)
	}
	searchRequest := session.MakeSimpleSearchRequest(s.SearchFilter(), []string{"distinguishedName"})
	searchResults, err := session.GetPagedSearchResults(searchRequest)
	if err != nil {
		return err
	}
	dn, err := utils.ChooseDN(searchResults)
	if err != nil {
		return err
	}

	sr := ldap.NewSearchRequest(
		dn,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(cn=*)",
		attrs,
		nil)
	return session.ExecuteSearchRequest(sr)

}
