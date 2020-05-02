package modules

import (
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/spf13/pflag"
)

type AdminObjects struct{}

func init() {
	AllModules = append(AllModules, new(AdminObjects))
}

func (AdminObjects) Name() string {
	return "admin-objects"
}

func (AdminObjects) Description() string {
	return "Enumerate all objects with protected ACLs (i.e admins)"
}

func (AdminObjects) FlagSet() *pflag.FlagSet {
	return pflag.NewFlagSet("adminobjects", pflag.ExitOnError)
}

func (AdminObjects) DefaultAttrs() []string {
	return []string{"distinguishedName"}
}

func (AdminObjects) Run(session *ldapsession.LDAPSession, attrs []string) (error) {
	filter := "(adminCount=1)"
	sr := session.MakeSimpleSearchRequest(filter, attrs)
	return session.ExecuteSearchRequest(sr)
}

