package modules

import (
	"fmt"
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/spf13/pflag"
	"gopkg.in/ldap.v3"
)

type UsersModule struct {
	ExtraFilter string
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
		return fmt.Sprintf("(&%s(%s))", filter, u.ExtraFilter)
	}
	return filter

}


func (u *UsersModule) FlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet(u.Name(), pflag.ExitOnError)
	flags.StringVar(&u.ExtraFilter, "filter", "", "Extra LDAP syntax filter to use")
	return flags
}

func (u *UsersModule) DefaultAttrs() []string {
	return []string{"cn", "sAMAccountName"}
}

func (u *UsersModule) Run(lSession *ldapsession.LDAPSession, attrs []string) (results *ldap.SearchResult, err error) {
	searchReq := lSession.MakeSimpleSearchRequest(u.Filter(), attrs)
	//return lSession.GetSearchResults(searchReq)
	ch := make(chan *ldap.Entry)
	go func() {
		for r := range ch {
				fmt.Printf("[+] Got an entry! %s\n", r.DN)
		}
	}()
	err = lSession.SearchWithPagingToChannel(searchReq, ch, 1000)
	return
}





