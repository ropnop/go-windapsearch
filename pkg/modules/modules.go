package modules

import (
	"fmt"
	"io"
)

var AllModules []WindapModule

type WindapModule struct {
	Name string
	Description string
	Filter string
	DefaultAttrs []string
	OutputOptions OutputOptions
}

type OutputOptions struct {
	ResolveHosts bool
	Attributes []string
	Full bool
	JSON bool
	Output io.Writer
}

type Module interface {
	GetHelp() string
	GetFilter() string
}

func (m *WindapModule) SetAttrs(attrs []string) {
	m.DefaultAttrs = attrs
}

func (m *WindapModule) AddFilter(filter string) {
	newFilter := fmt.Sprintf("(&(%s)(%s))", m.Filter, filter)
	m.Filter = newFilter
}

