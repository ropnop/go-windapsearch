package modules

import (
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/spf13/pflag"
)

type Module interface {
	Name() string
	Description() string
	FlagSet() *pflag.FlagSet
	DefaultAttrs() []string
	Run(session *ldapsession.LDAPSession, attrs []string) error
}

var AllModules []Module

