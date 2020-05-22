package modules

import (
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/ropnop/go-windapsearch/pkg/utils"
	"github.com/spf13/pflag"
)

type UsersModule struct {
	ExtraFilter string
	SearchTerm  string
}

func init() {
	AllModules = append(AllModules, new(UsersModule))
}

func (u *UsersModule) Name() string {
	return "users"
}

func (u *UsersModule) Description() string {
	return "List all user objects"
}

func (u *UsersModule) Filter() string {
	filter := "(objectcategory=user)"
	if u.ExtraFilter != "" {
		//return fmt.Sprintf("(&%s(%s))", filter, u.ExtraFilter)
		filter = utils.AddAndFilter(filter, u.ExtraFilter)
	}
	if u.SearchTerm != "" {
		filter = utils.AddAndFilter(filter, utils.CreateANRSearch(u.SearchTerm))
	}
	return filter

}

func (u *UsersModule) FlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet(u.Name(), pflag.ExitOnError)
	flags.StringVar(&u.ExtraFilter, "filter", "", "Extra LDAP syntax filter to use")
	flags.StringVarP(&u.SearchTerm, "search", "s", "", "Search term to filter on")
	return flags
}

func (u *UsersModule) DefaultAttrs() []string {
	return []string{"cn", "sAMAccountName"}
}

func (u *UsersModule) Run(lSession *ldapsession.LDAPSession, attrs []string) error {
	searchReq := lSession.MakeSimpleSearchRequest(u.Filter(), attrs)
	return lSession.ExecuteSearchRequest(searchReq)

}
