package modules

import (
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/ropnop/go-windapsearch/pkg/utils"
	"github.com/spf13/pflag"
)

type UnconstrainedModule struct {
	Users     bool
	Computers bool
}

func init() {
	AllModules = append(AllModules, new(UnconstrainedModule))
}

func (u UnconstrainedModule) Name() string {
	return "unconstrained"
}

func (u UnconstrainedModule) Description() string {
	return "find objects that allow unconstrained delegation"
}

func (u *UnconstrainedModule) FlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("unconstrained-module", pflag.ExitOnError)
	flags.BoolVar(&u.Users, "users", false, "Only show users")
	flags.BoolVar(&u.Computers, "computers", false, "Only show computers")
	return flags
}

func (u UnconstrainedModule) DefaultAttrs() []string {
	return []string{"cn", "sAMAccountName"}
}

func (u *UnconstrainedModule) Filter() string {
	filter := "(userAccountControl:1.2.840.113556.1.4.803:=524288)"
	if u.Users {
		usersFilter := utils.AddAndFilter("(objectClass=user)", "(objectCategory=user)")
		filter = utils.AddAndFilter(filter, usersFilter)
	}
	if u.Computers {
		compFilter := utils.AddAndFilter("(objectCategory=computer)", "(objectClass=computer)")
		filter = utils.AddAndFilter(filter, compFilter)
	}
	return filter
}

func (u *UnconstrainedModule) Run(session *ldapsession.LDAPSession, attrs []string) error {
	sr := session.MakeSimpleSearchRequest(u.Filter(), attrs)
	return session.ExecuteSearchRequest(sr)
}
