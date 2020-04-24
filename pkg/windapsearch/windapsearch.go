package windapsearch

import (
	"errors"
	"fmt"
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/ropnop/go-windapsearch/pkg/modules"
	"github.com/ropnop/go-windapsearch/pkg/utils"
	"github.com/spf13/pflag"
	"gopkg.in/ldap.v3"
	"io"
	"os"
	"strings"
)

type WindapSearchSession struct {
	Options      CommandLineOptions
	LDAPSession  *ldapsession.LDAPSession
	Module modules.Module
	AllModules   []modules.Module
	OutputWriter io.Writer
}


var RootFlagSet *pflag.FlagSet

func init() {
	RootFlagSet = pflag.NewFlagSet("windapsearch", pflag.ExitOnError)
	RootFlagSet.SortFlags = false
	RootFlagSet.StringP("domain", "d", "", "The FQDN of the domain (e.g. 'lab.example.com'). Only needed if dc not provided")
	RootFlagSet.String("dc", "", "The Domain Controller to query against")
	RootFlagSet.StringP( "username", "u", "", "The full username with domain to bind with (e.g. 'ropnop@lab.example.com' or 'LAB\\ropnop')\n If not specified, will attempt anonymous bind")
	RootFlagSet.StringP( "password", "p", "", "Password to use. If not specified, will be prompted for")
	RootFlagSet.Int( "port", 0, "Port to connect to (if non standard)")
	RootFlagSet.Bool( "secure", false, "Use LDAPS. This will not verify TLS certs, however. (default: false)" )
	RootFlagSet.BoolP( "resolve", "r", false, "Resolve IP addresses for enumerated computer names. Will make DNS queries against system NS")
	RootFlagSet.StringSlice( "attrs", nil, "Comma separated custom atrribute names to display (e.g. 'badPwdCount,lastLogon')")
	RootFlagSet.Bool( "full", false, "Output all attributes from LDAP")
	RootFlagSet.StringP( "output", "o", "", "Save results to file")
	RootFlagSet.BoolP( "json", "j", false, "Convert LDAP output to JSON" )
	RootFlagSet.BoolP( "interactive", "i", false, "Start in interactive mode")
}

type CommandLineOptions struct {
	FlagSet *pflag.FlagSet
	Help	bool
	Domain           string
	DomainController string
	Username string
	Password string
	Port int
	Secure bool
	ResolveHosts   bool
	Attributes     []string
	FullAttributes bool
	Output      string
	JSON bool
	Module string
	Interactive bool
}


func NewSession() *WindapSearchSession {
	var w WindapSearchSession

	wFlags := pflag.NewFlagSet("WindapSearch", pflag.ContinueOnError)
	wFlags.SortFlags = false
	wFlags.StringVarP(&w.Options.Domain, "domain", "d", "", "The FQDN of the domain (e.g. 'lab.example.com'). Only needed if dc not provided")
	wFlags.StringVar(&w.Options.DomainController, "dc", "", "The Domain Controller to query against")
	wFlags.StringVarP(&w.Options.Username, "username", "u", "", "The full username with domain to bind with (e.g. 'ropnop@lab.example.com' or 'LAB\\ropnop')\n If not specified, will attempt anonymous bind")
	wFlags.StringVarP(&w.Options.Password, "password", "p", "", "Password to use. If not specified, will be prompted for")
	wFlags.IntVar(&w.Options.Port, "port", 0, "Port to connect to (if non standard)")
	wFlags.BoolVar(&w.Options.Secure, "secure", false, "Use LDAPS. This will not verify TLS certs, however. (default: false)" )
	wFlags.BoolVarP(&w.Options.ResolveHosts, "resolve", "r", false, "Resolve IP addresses for enumerated computer names. Will make DNS queries against system NS")
	wFlags.StringSliceVar(&w.Options.Attributes, "attrs", nil, "Comma separated custom atrribute names to display (e.g. 'badPwdCount,lastLogon')")
	wFlags.BoolVar(&w.Options.FullAttributes, "full", false, "Output all attributes from LDAP")
	wFlags.StringVarP(&w.Options.Output, "output", "o", "", "Save results to file")
	wFlags.BoolVarP(&w.Options.JSON, "json", "j", false, "Convert LDAP output to JSON" )
	wFlags.BoolVarP(&w.Options.Interactive, "interactive", "i", false, "Start in interactive mode")
	wFlags.BoolVarP(&w.Options.Help, "help", "h", false, "Show this help")

	pflag.ErrHelp = errors.New("")
	wFlags.Usage = w.ShowUsage

	for _, m := range modules.AllModules {
		w.RegisterModule(m)
	}

	//wFlags.StringP("module", "m", "", fmt.Sprintf("Module to use. Available modules: \n[ %s ]", w.ModuleListString()))
	wFlags.StringVarP(&w.Options.Module, "module", "m", "", "Module to use")

	w.Options.FlagSet = wFlags

	w.OutputWriter = os.Stdout //default to stdout
	return &w
}

func (w *WindapSearchSession) RegisterModule(mod modules.Module) {
	w.AllModules = append(w.AllModules, mod)
}

func (w *WindapSearchSession) ModuleListString() string {
	var sb strings.Builder
	for _, mod := range w.AllModules {
		sb.WriteString(mod.Name())
		sb.WriteString(", ")
	}
	listString := sb.String()
	return strings.TrimSuffix(listString, ", ")
}

func (w *WindapSearchSession) ModuleDescriptionString() string {
	var sb strings.Builder
	for _, mod := range w.AllModules {
		sb.WriteString(fmt.Sprintf("\t%s\t\t%s\n", mod.Name(), mod.Description()))
	}
	return sb.String()
}

func (w *WindapSearchSession) GetModuleByName(name string) modules.Module {
	for _, m := range w.AllModules {
		if m.Name() == name {
			return m
		}
	}
	return nil
}

func (w *WindapSearchSession) ShowUsage() {
	fmt.Fprintf(os.Stderr, "windapsearch: a tool to perform Windows domain enumeration through LDAP queries\n\nUsage: %s [options] -m [module]\n\nOptions:\n", os.Args[0])
	w.Options.FlagSet.PrintDefaults()
	if w.Module == nil {
		fmt.Fprintf(os.Stderr, "\nAvailable modules:\n%s", w.ModuleDescriptionString())
	} else {
		fmt.Fprintf(os.Stderr, "\nOptions for %q module:\n", w.Module.Name())
		w.Module.FlagSet().PrintDefaults()
	}
}

func (w *WindapSearchSession) Run() (err error) {
	defer func() {
		err = wrap(err)
	}()

	w.Options.FlagSet.Parse(os.Args[:])

	w.LoadModule()


	if w.Options.Help {
		w.ShowUsage()
		return
	}

	if w.Options.Output != "" {
		fp, err2 := os.Create(w.Options.Output)
		if err2 != nil { err = err2; return }
		w.OutputWriter = fp
		defer fp.Close()
	}

	if w.Options.Domain == "" && w.Options.DomainController == "" {
		w.ShowUsage()
		fmt.Println()
		fmt.Println("[!] You must specify either a domain or an IP address of a domain controller")
		return
	}
	password := w.Options.Password
	if w.Options.Username != "" && password == "" {
		password, err = utils.SecurePrompt(fmt.Sprintf("Password for [%s]", w.Options.Username))
		if err != nil { return err }
	}

	ldapOptions := ldapsession.LDAPSessionOptions{
		Domain:           w.Options.Domain,
		DomainController: w.Options.DomainController,
		Username:         w.Options.Username,
		Password:         password,
		Port:             w.Options.Port,
		Secure:           w.Options.Secure,
	}

	w.LDAPSession, err = ldapsession.NewLDAPSession(&ldapOptions)
	if err != nil { return }
	defer w.LDAPSession.Close()

	if w.Options.Interactive {
		return w.StartTUI()
	} else {
		return w.StartCLI()
	}
}

func (w *WindapSearchSession) LoadModule() {
	mod := w.GetModuleByName(w.Options.Module)
	if mod != nil {
		w.Module = mod
	}
}

func (w *WindapSearchSession) StartCLI() error {
	if w.Module == nil {
		fmt.Printf("[!] You must specify a valid module to use\n")
		fmt.Printf(" Available modules: \n%s", w.ModuleDescriptionString())
		return nil
	}

	mod := w.GetModuleByName(w.Options.Module)
	if mod == nil {
		w.ShowUsage()
		fmt.Println()
		fmt.Printf("[!] Module %q not found\n", w.Options.Module)
		return nil
	}

	modFlags := mod.FlagSet()
	modFlags.AddFlagSet(w.Options.FlagSet)
	modFlags.Parse(os.Args[:])



	var attrs []string
	if w.Options.FullAttributes {
		attrs = []string{"*"}
	} else if len(w.Options.Attributes) > 0 {
		attrs = w.Options.Attributes
	} else {
		attrs = mod.DefaultAttrs()
	}
	results, err := mod.Run(w.LDAPSession, attrs)
	if err != nil  { return err }


	err = w.handleResults(results)
	if err != nil { return err }
	if w.Options.Output != "" {
		fmt.Printf("[+] %s written\n", w.Options.Output)
	}
	return nil
}

func (w *WindapSearchSession) StartTUI() error {
	return nil
}

func (w *WindapSearchSession) handleResults(results *ldap.SearchResult) error {
	if w.Options.JSON {
		jResults, err := utils.SearchResultToJSON(results)
		if err != nil {
			return err
		}
		w.OutputWriter.Write(jResults)
	} else {
		utils.WriteSearchResults(results, w.OutputWriter)
	}
	return nil

}



