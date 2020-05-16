package modules

import (
	"fmt"
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/spf13/pflag"
	"strings"
)

var PrivilegedGroups = append([]string{
	"Administrators",  // Builtin administrators group for the domain
	"Enterprise Admins",
	"Schema Admins",  // Highly privileged builtin group
	"Account Operators",
	"Backup Operators",
	"Server Management",
	"Konten-Operatoren",
	"Sicherungs-Operatoren",
	"Server-Operatoren",
	"Schema-Admins",
}, DomainAdminGroups...)

type PrivilegedObjectsModule struct{}

func init() {
	AllModules = append(AllModules, new(PrivilegedObjectsModule))
}

func (p PrivilegedObjectsModule) Name() string {
	return "privileged-users"
}

func (p PrivilegedObjectsModule) Description() string {
	return "Recursively list members of all highly privileged groups (i.e. Domain Admins, Enterprise Admins, Schema Admins, etc...)"
}

func (p PrivilegedObjectsModule) FlagSet() *pflag.FlagSet {
	return pflag.NewFlagSet("privileged-objects", pflag.ExitOnError)
}

func (p PrivilegedObjectsModule) DefaultAttrs() []string {
	return []string{"cn", "sAMAccountName"}
}

func (PrivilegedObjectsModule) Filter(baseDN string) string {
	var sb strings.Builder
	sb.WriteString("(&(objectClass=user)(|")
	for _, group := range PrivilegedGroups{
		filter := fmt.Sprintf("(memberof:1.2.840.113556.1.4.1941:=CN=%s,CN=Users,%s)",group,baseDN)
		sb.WriteString(filter)
	}
	sb.WriteString("))")
	return sb.String()

}

func (p PrivilegedObjectsModule) Run(session *ldapsession.LDAPSession, attrs []string) error {
	filter := p.Filter(session.BaseDN)
	searchReq := session.MakeSimpleSearchRequest(filter, attrs)
	return session.ExecuteSearchRequest(searchReq)
}

