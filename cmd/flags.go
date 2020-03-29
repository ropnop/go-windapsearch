package cmd

import (
	"github.com/jessevdk/go-flags"
)

var DomainOptions struct {
	Domain           string `group:"Domain Options" short:"d" long:"domain" description:"The FQDN of the domain (e.g. 'lab.example.com'). Only needed if dc not provided"`
	DomainController string `group:"Domain Options" long:"dc" description:"The Domain Controller to query against"`
}

var BindOptions struct {
	Username string `group:"Bind Options" short:"u" long:"username" description:"The full username with domain to bind with (e.g. 'ropnop@lab.example.com' or 'LAB\\ropnop'"`
	Password string `group:"Bind Options" short:"p" long:"password" description:"Password to use. If not specified, will be prompted for"`
	Port int `group:"Bind Options" long:"port" description:"Port to connect to (if non standard)"`
	Secure bool `group:"Bind Options" long:"secure" description:"Use LDAPS. This will not verify TLS certs, however. (default: false)"`
}

var EnumerationOptions struct {
	Groups       bool   `short:"G" long:"groups" description:"Enumerate all AD Groups"`
	Users        bool   `short:"U" long:"users" description:"Enumerate all AD Users"`
	Computers    bool   `short:"C" long:"computers" description:"Enumerate all AD Computers"`
	GroupName    string `short:"m" long:"members" description:"Enumerate all members of a group"`
	DomainAdmins bool   `long:"da" description:"Shortcut for enumerate all members of group 'Domain Admins'. Performs recursive lookups for nested members."`
	SearchTerm   string `short:"s" long:"search" description:"Fuzzy search for all matching LDAP entries"`
	LookupDN     string `short:"l" long:"lookup" description:"Search through LDAP and lookup entry. Works with fuzzy search. Defaults to printing all attributes, but honors '--attrs'"`
	CustomFilter string `short:"f" long:"filter" description:"Search with a fully custom filter. Must be valid LDAP Filter Syntax"`
}

var OutputOptions struct {
	ResolveHosts   bool   `short:"r" long:"resolve" description:"Resolve IP addresses for enumerated computer names. Will make DNS queries against system NS"`
	Attributes     string `long:"attrs" description:"Comma separated custom atrribute names to display (e.g. 'badPwdCount,lastLogon')"`
	FullAttributes bool   `long:"full" description:"Output all attributes from LDAP"`
	Output      string `short:"o" long:"output" description:"Save results to file"`
	JSON bool `short:"j" long:"json" description:"Output as JSON format"`
}

var OptionsParser = flags.NewNamedParser("windapsearchgo", flags.Default)

func init() {
	OptionsParser.AddGroup("Domain Options", "", &DomainOptions)
	OptionsParser.AddGroup("Bind Options", "Specify bind account. If not specified, anonymous bind will be attempted", &BindOptions)
	OptionsParser.AddGroup("Enumeration Options", "Data to enumerate from LDAP", &EnumerationOptions)
	OptionsParser.AddGroup("Output Options", "Display and output options for results", &OutputOptions)
}