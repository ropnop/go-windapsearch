package modules

import (
	"github.com/ropnop/go-windapsearch/pkg/adschema"
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


type ModuleWithDisplayFilter interface {
	//DisplayFilter is an optional function for a module used to filter what results get written to the output
	//it's best to try and filter the request, but when that's not possible, a custom function can be provided
	//that takes an LDAP entry and returns true/false (whether to display or not)
	DisplayFilter(e *adschema.ADEntry) bool
}

// Hacky way to have "DisplayFilter" be optional for the Module interface - if it exists, call it, if not, just return true
func DisplayFilter(m Module, e *adschema.ADEntry) bool {
	if moduleWithDisplayFilter, ok := m.(ModuleWithDisplayFilter); ok {
		return moduleWithDisplayFilter.DisplayFilter(e)
	} else {
		return true
	}
}

var AllModules []Module
