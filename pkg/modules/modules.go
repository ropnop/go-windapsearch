package modules

import (
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"gopkg.in/ldap.v3"
)

type Module interface {
	Name() string
	Description() string
	Filter() string
	Options() interface{}
	DefaultAttrs() []string
	Run(session *ldapsession.LDAPSession, attrs []string) (*ldap.SearchResult, error)
}

var AllModules []Module

