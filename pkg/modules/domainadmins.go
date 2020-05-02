package modules

import (
	"fmt"
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/spf13/pflag"
	"strings"
)

var DomainAdminGroups = []string{
"Domain Admins",
"Domain-Admins",
"Domain Administrators",
"Domain-Administrators",
"Domänen Admins",
"Domänen-Admins",
"Domain Admins",
"Domain-Admins",
"Domänen Administratoren",
"Domänen-Administratoren",
}

type DAModule struct{}

func init() {
	AllModules = append(AllModules, new(DAModule))
}

func (DAModule) Name() string {
	return "domain-admins"
}

func (DAModule) Description() string {
	return "Recursively list all users objects in Domain Admins group"
}

func (DAModule) Filter(baseDN string) string {
	var sb strings.Builder
	sb.WriteString("(&(objectClass=user)(|")
	for _, group := range DomainAdminGroups {
		filter := fmt.Sprintf("(memberof:1.2.840.113556.1.4.1941:=CN=%s,CN=Users,%s)",group,baseDN)
		sb.WriteString(filter)
	}
	sb.WriteString("))")
	return sb.String()
}

func (DAModule) FlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("domain-admins", pflag.ExitOnError)
	return flags
}

func (DAModule) DefaultAttrs() []string {
	return []string{"cn", "sAMAccountName"}
}

func (d DAModule) Run(session *ldapsession.LDAPSession, attrs []string) error {
	filter := d.Filter(session.BaseDN)
	searchReq := session.MakeSimpleSearchRequest(filter, attrs)
	return session.ExecuteSearchRequest(searchReq)
}

