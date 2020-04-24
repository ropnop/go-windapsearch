package modules

import (
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/spf13/pflag"
	"gopkg.in/ldap.v3"
)

type Module interface {
	Name() string
	Description() string
	Filter() string
	FlagSet() *pflag.FlagSet
	DefaultAttrs() []string
	Run(session *ldapsession.LDAPSession, attrs []string) (*ldap.SearchResult, error)
}

var AllModules []Module

