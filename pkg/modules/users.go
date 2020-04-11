package modules

import (
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"gopkg.in/ldap.v3"
)

type UsersModule struct {
	Flags usersFlags
}

type usersFlags struct {
	ExtraFilter string `short:"f" long:"filter" description:"Extra LDAP syntax search filter"`
}

func init() {
	AllModules = append(AllModules, new(UsersModule))
}

func (u UsersModule) Name() string {
	return "users"
}

func (u UsersModule) Description() string {
	return "List all user objects"
}

func (u UsersModule) Filter() string {
	return "(objectcategory=user)"
}

func (u UsersModule) Options() interface{} {
	return u.Flags
}

func (u UsersModule) DefaultAttrs() []string {
	return []string{"cn", "sAMAccountName"}
}

func (u UsersModule) Run(lSession *ldapsession.LDAPSession, attrs []string) (results *ldap.SearchResult, err error) {
	searchReq := lSession.MakeSimpleSearchRequest(u.Filter(), attrs)
	return lSession.GetSearchResults(searchReq)
}





