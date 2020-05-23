package modules

import (
	"fmt"
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/ropnop/go-windapsearch/pkg/utils"
	"github.com/spf13/pflag"
	"os"
)

type MembersModule struct {
	Recursive bool
	Search    string
	DN        string
	OnlyUsers bool
}

func init() {
	AllModules = append(AllModules, new(MembersModule))
}

func (m MembersModule) Name() string {
	return "members"
}

func (m MembersModule) Description() string {
	return "Query for members of a group"
}

func (m *MembersModule) FlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("members-module", pflag.ExitOnError)
	flags.BoolVarP(&m.Recursive, "recursive", "r", false, "Perform recursive lookup")
	flags.StringVarP(&m.Search, "search", "s", "", "Search for group name")
	flags.StringVarP(&m.DN, "group", "g", "", "Full DN of group to enumerate")
	flags.BoolVar(&m.OnlyUsers, "users", false, "Only return user objects")
	return flags
}

func (m MembersModule) DefaultAttrs() []string {
	return []string{"cn", "sAMAccountName"}
}

func (m *MembersModule) ChooseGroup(session *ldapsession.LDAPSession) (dn string, err error) {
	filter := "(objectcategory=group)"
	filter = utils.AddAndFilter(filter, utils.CreateANRSearch(m.Search))
	sr := session.MakeSimpleSearchRequest(filter, []string{})
	matchResults, err := session.GetPagedSearchResults(sr)
	if err != nil {
		return
	}
	return utils.ChooseDN(matchResults)
}

func (m MembersModule) Filter() string {
	var filter string
	if m.Recursive {
		filter = fmt.Sprintf("(memberof:1.2.840.113556.1.4.1941:=%s)", m.DN)
	} else {
		filter = fmt.Sprintf("(memberOf=%s)", m.DN)
	}
	if m.OnlyUsers {
		filter = utils.AddAndFilter(filter, "(objectcategory=user)")
	}
	return filter
}

func (m *MembersModule) Run(session *ldapsession.LDAPSession, attrs []string) error {
	if m.DN == "" && m.Search == "" {
		return fmt.Errorf("must provide a group or a search term")
	}
	if m.DN == "" {
		dn, err := m.ChooseGroup(session)
		if err != nil {
			return err
		}
		m.DN = dn
		fmt.Fprintf(os.Stderr, "[+] Using group: %s\n\n", m.DN)
	}
	sr := session.MakeSimpleSearchRequest(m.Filter(), attrs)
	return session.ExecuteSearchRequest(sr)

}
