package windapsearch

import (
	"fmt"
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/ropnop/go-windapsearch/pkg/modules"
	"github.com/ropnop/go-windapsearch/pkg/utils"
	flag "github.com/spf13/pflag"
	"gopkg.in/ldap.v3"
	"io"
	"os"
	"strings"
)

type WindapSearchSession struct {
	Options CommandLineOptions
	ModuleOptions    ModuleOptions
	LDAPSession *ldapsession.LDAPSession
	Modules []modules.Module
	OutputWriter io.Writer
}

var (
	Domain string
	DomainController string
	Username string
	Password string
	Port int
	Secure bool
	Module string
)

type CommandLineOptions struct {
	FlagSet *flag.FlagSet
	Domain           string `group:"Domain Options" short:"d" long:"domain" description:"The FQDN of the domain (e.g. 'lab.example.com'). Only needed if dc not provided"`
	DomainController string `group:"Domain Options" long:"dc" description:"The Domain Controller to query against"`
	Username string `group:"Bind Options" short:"u" long:"username" description:"The full username with domain to bind with (e.g. 'ropnop@lab.example.com' or 'LAB\\ropnop') If not specified, will attempt anonymous bind"`
	Password string `group:"Bind Options" short:"p" long:"password" description:"Password to use. If not specified, will be prompted for"`
	Port int `group:"Bind Options" long:"port" description:"Port to connect to (if non standard)"`
	Secure bool `group:"Bind Options" long:"secure" description:"Use LDAPS. This will not verify TLS certs, however. (default: false)"`
	ResolveHosts   bool   `short:"r" long:"resolve" description:"Resolve IP addresses for enumerated computer names. Will make DNS queries against system NS"`
	Attributes     []string `long:"attrs" description:"Comma separated custom atrribute names to display (e.g. 'badPwdCount,lastLogon')"`
	FullAttributes bool   `long:"full" description:"Output all attributes from LDAP"`
	Output      string `short:"o" long:"output" description:"Save results to file"`
	JSON bool `short:"j" long:"json" description:"Output as JSON format"`
	Module string
	Interactive bool
}


type ModuleOptions struct {
	FlagSet *flag.FlagSet
}

func NewSession() *WindapSearchSession {
	var w WindapSearchSession

	wFlags := flag.NewFlagSet("WindapSearch", flag.ExitOnError)
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

	for _, m := range modules.AllModules {
		w.RegisterModule(m)
	}

	wFlags.StringVarP(&w.Options.Module, "module", "m", "", fmt.Sprintf("Module to use. Available modules: \n[ %s ]", w.ModuleListString()))

	w.Options.FlagSet = wFlags

	w.OutputWriter = os.Stdout //default to stdout
	return &w
}

func (w *WindapSearchSession) RegisterModule(mod modules.Module) {
	w.Modules = append(w.Modules, mod)
}

func (w *WindapSearchSession) ModuleListString() string {
	var sb strings.Builder
	for _, mod := range w.Modules {
		sb.WriteString(mod.Name())
		sb.WriteString(", ")
	}
	listString := sb.String()
	return strings.TrimSuffix(listString, ", ")
}

func (w *WindapSearchSession) ModuleDescriptionString() string {
	var sb strings.Builder
	for _, mod := range w.Modules {
		sb.WriteString(fmt.Sprintf("\t%s\t\t%s\n", mod.Name(), mod.Description()))
	}
	return sb.String()
}

func (w *WindapSearchSession) GetModuleByName(name string) modules.Module {
	for _, m := range w.Modules {
		if m.Name() == name {
			return m
		}
	}
	return nil
}

func (w *WindapSearchSession) Run() (err error) {
	defer func() {
		err = wrap(err)
	}()

	w.Options.FlagSet.Parse(os.Args[:])

	if w.Options.Output != "" {
		fp, err2 := os.Create(w.Options.Output)
		if err2 != nil { err = err2; return }
		w.OutputWriter = fp
		defer fp.Close()
	}

	if w.Options.Domain == "" && w.Options.DomainController == "" {
		w.Options.FlagSet.PrintDefaults()
		fmt.Println("[!] You must specify either a domain or an IP address of a domain controller")
		return
	}
	password := w.Options.Password
	if w.Options.Username != "" && password == "" {
		password, err = utils.SecurePrompt(fmt.Sprintf("Password for [%s]", w.Options.Username))
		if err != nil { return }
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

func (w *WindapSearchSession) StartCLI() error {
	if w.Options.Module == "" {
		fmt.Printf("[!] You must specify a module to use\n")
		fmt.Printf(" Available modules: \n%s", w.ModuleDescriptionString())
		return nil
	}

	mod := w.GetModuleByName(w.Options.Module)
	if mod == nil {
		fmt.Printf("[!] Module %q not found\n", w.Options.Module)
		fmt.Printf(" Available modules: \n%s", w.ModuleDescriptionString())
		return nil
	}

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



