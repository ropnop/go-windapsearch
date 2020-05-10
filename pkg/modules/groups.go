package modules

import (
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/ropnop/go-windapsearch/pkg/utils"
	"github.com/spf13/pflag"
)

type GroupsModule struct {
	SearchTerm string
}

func init() {
	AllModules = append(AllModules, new(GroupsModule))
}

func (g *GroupsModule) Name() string {
	return "groups"
}

func (g *GroupsModule) Description() string {
	return "List all AD groups"
}

func (g *GroupsModule) Filter() string {
	filter := "(objectcategory=group)"
	if g.SearchTerm != "" {
		filter = utils.AddAndFilter(filter, utils.CreateANRSearch(g.SearchTerm))
	}
	return filter
}

func (g *GroupsModule) FlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet(g.Name(), pflag.ExitOnError)
	flags.StringVarP(&g.SearchTerm, "search", "s", "", "Search term to filter on")
	return flags
}

func (g *GroupsModule) DefaultAttrs() []string {
	return []string{"distinguishedName", "cn"}
}

func (g *GroupsModule) Run(lSession *ldapsession.LDAPSession, attrs []string) error {
	searchReq := lSession.MakeSimpleSearchRequest(g.Filter(), attrs)
	return lSession.ExecuteSearchRequest(searchReq)
}