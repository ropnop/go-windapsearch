package windapsearch

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"text/tabwriter"

	"github.com/ropnop/go-windapsearch/pkg/buildinfo"
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"github.com/ropnop/go-windapsearch/pkg/modules"
	"github.com/ropnop/go-windapsearch/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type WindapSearchSession struct {
	Options      CommandLineOptions
	LDAPSession  *ldapsession.LDAPSession
	Module       modules.Module
	AllModules   []modules.Module
	Log          *logrus.Entry
	OutputWriter io.Writer
	workers      int
	ctx          context.Context
	cancel       context.CancelFunc
}

type CommandLineOptions struct {
	FlagSet          *pflag.FlagSet
	Help             bool
	Domain           string
	DomainController string
	Username             string
	Password             string
	NTLMHash             string
	UseNTLM              bool
	Port                 int
	Proxy                string
	Secure               bool
	ResolveHosts         bool
	Attributes           []string
	FullAttributes       bool
	IgnoreDisplayFilters bool
	Output               string
	JSON                 bool
	Module               string
	Interactive          bool
	Version              bool
	Verbose              bool
	Debug                bool
	PageSize             int
	ModuleFlags          *pflag.FlagSet
}

func NewSession() *WindapSearchSession {
	var w WindapSearchSession

	wFlags := pflag.NewFlagSet("WindapSearch", pflag.ContinueOnError)
	wFlags.SortFlags = false
	wFlags.StringVarP(&w.Options.Domain, "domain", "d", "", "The FQDN of the domain (e.g. 'lab.example.com'). Only needed if dc not provided")
	wFlags.StringVar(&w.Options.DomainController, "dc", "", "The Domain Controller to query against")
	wFlags.StringVarP(&w.Options.Username, "username", "u", "", "The full username with domain to bind with (e.g. 'ropnop@lab.example.com' or 'LAB\\ropnop')\n If not specified, will attempt anonymous bind")
	wFlags.StringVarP(&w.Options.Password, "password", "p", "", "Password to use. If not specified, will be prompted for")
	wFlags.StringVar(&w.Options.NTLMHash, "hash", "", "NTLM Hash to use instead of password (i.e. pass-the-hash)")
	wFlags.BoolVar(&w.Options.UseNTLM, "ntlm", false, "Use NTLM auth (automatic if hash is set)")
	wFlags.IntVar(&w.Options.Port, "port", 0, "Port to connect to (if non standard)")
	wFlags.BoolVar(&w.Options.Secure, "secure", false, "Use LDAPS. This will not verify TLS certs, however. (default: false)")
	wFlags.StringVar(&w.Options.Proxy, "proxy", "", "SOCKS5 Proxy to use (e.g. 127.0.0.1:9050)")
	wFlags.BoolVar(&w.Options.FullAttributes, "full", false, "Output all attributes from LDAP")
	wFlags.BoolVar(&w.Options.IgnoreDisplayFilters, "ignore-display-filters", false, "Ignore any display filters set by the module and always output every entry")
	wFlags.StringVarP(&w.Options.Output, "output", "o", "", "Save results to file")
	wFlags.BoolVarP(&w.Options.JSON, "json", "j", false, "Convert LDAP output to JSON")
	wFlags.IntVar(&w.Options.PageSize, "page-size", 1000, "LDAP page size to use")
	//wFlags.BoolVarP(&w.Options.Interactive, "interactive", "i", false, "Start in interactive mode") //TODO
	wFlags.BoolVar(&w.Options.Version, "version", false, "Show version info and exit")
	wFlags.BoolVarP(&w.Options.Verbose, "verbose", "v", false, "Show info logs")
	wFlags.BoolVar(&w.Options.Debug, "debug", false, "Show debug logs")
	wFlags.BoolVarP(&w.Options.Help, "help", "h", false, "Show this help")

	pflag.ErrHelp = errors.New("")
	wFlags.Usage = w.ShowUsage

	for _, m := range modules.AllModules {
		w.RegisterModule(m)
	}

	wFlags.StringVarP(&w.Options.Module, "module", "m", "", "Module to use")

	w.Options.FlagSet = wFlags

	w.OutputWriter = os.Stdout //default to stdout
	w.workers = 5              //concurrent workers for marshaling entries. 5 seems reasonable

	logger := logrus.New()

	logger.Out = os.Stderr // default log to stderr
	logger.SetLevel(logrus.ErrorLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})
	w.Log = logger.WithFields(logrus.Fields{"package": "windapsearch"})

	return &w
}

func (w *WindapSearchSession) handleInterrupt() {
	// set up cancelling, catch SIGINT
	w.ctx, w.cancel = context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		w.cancel()
	}()
}

func (w *WindapSearchSession) RegisterModule(mod modules.Module) {
	w.AllModules = append(w.AllModules, mod)
}

func (w *WindapSearchSession) LoadModule() {
	mod := w.GetModuleByName(w.Options.Module)
	if mod != nil {
		w.Module = mod
		w.Options.ModuleFlags = mod.FlagSet()
		w.Options.ModuleFlags.StringSliceVar(&w.Options.Attributes, "attrs", mod.DefaultAttrs(), "Comma separated custom atrributes to display")
	}
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
	sb := &strings.Builder{}
	tw := tabwriter.NewWriter(sb, 0, 0, 4, ' ', 0)
	for _, mod := range w.AllModules {
		fmt.Fprintf(tw, "\t%s\t%s\n", mod.Name(), mod.Description())
	}
	tw.Flush()
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
	fmt.Fprintf(os.Stderr, "windapsearch: a tool to perform Windows domain enumeration through LDAP queries\n%s\nUsage: %s [options] -m [module] [module options]\n\nOptions:\n", buildinfo.FormatVersionString(), os.Args[0])
	w.Options.FlagSet.PrintDefaults()
	if w.Module == nil {
		fmt.Fprintf(os.Stderr, "\nAvailable modules:\n%s", w.ModuleDescriptionString())
	} else {
		fmt.Fprintf(os.Stderr, "\nOptions for %q module:\n", w.Module.Name())
		w.Options.ModuleFlags.PrintDefaults()
	}
}

func (w *WindapSearchSession) Run() (err error) {
	defer func() {
		err = wrap(err)
	}()

	w.Options.FlagSet.Parse(os.Args[:])

	w.LoadModule()

	//w.Options.ModuleFlags.AddFlagSet(w.Options.FlagSet)
	w.Options.FlagSet.AddFlagSet(w.Options.ModuleFlags)
	w.Options.FlagSet.Parse(os.Args[:])

	if w.Options.Help {
		w.ShowUsage()
		return
	}

	if w.Options.Version {
		fmt.Println(buildinfo.FormatVersionString())
		return
	}

	if w.Options.Verbose {
		w.Log.Logger.SetLevel(logrus.InfoLevel)
	}
	if w.Options.Debug {
		w.Log.Logger.SetLevel(logrus.DebugLevel)
	}

	if w.Options.Output != "" {
		fp, err2 := os.Create(w.Options.Output)
		if err2 != nil {
			err = err2
			return
		}
		w.OutputWriter = fp
		defer fp.Close()
		w.Log.Infof("Saving output to %q", fp.Name())
	} else {
		w.Log.Infof("Saving output to STDOUT")
	}

	if w.Options.Domain == "" && w.Options.DomainController == "" {
		w.ShowUsage()
		fmt.Fprintf(os.Stderr, "\n[!] You must specify either a domain or an IP address of a domain controller\n")
		return
	}
	password := w.Options.Password
	username := w.Options.Username

	if w.Options.UseNTLM && username == "" {
		return fmt.Errorf("must provide username for NTLM authentication")
	}

	if username != "" { // only prompt for password if username is provided
		if len(strings.Split(w.Options.Username, "@")) == 1 {
			username = fmt.Sprintf("%s@%s", w.Options.Username, w.Options.Domain)
		} else {
			username = w.Options.Username
		}
		if username != "" && password == "" && w.Options.NTLMHash == "" {
			password, err = utils.SecurePrompt(fmt.Sprintf("Password for [%s]", username))
			if err != nil {
				return err
			}
		}
	}

	// now that ldap connections are opened, handle interrupts gracefully
	w.handleInterrupt()

	ldapOptions := ldapsession.LDAPSessionOptions{
		Domain:           w.Options.Domain,
		DomainController: w.Options.DomainController,
		Username:         username,
		Password:         password,
		Hash:             w.Options.NTLMHash,
		UseNTLM:          w.Options.UseNTLM,
		Port:             w.Options.Port,
		Proxy:            w.Options.Proxy,
		Secure:           w.Options.Secure,
		PageSize:         w.Options.PageSize,
		Logger:           w.Log.Logger,
	}

	w.LDAPSession, err = ldapsession.NewLDAPSession(&ldapOptions, w.ctx)
	if err != nil {
		return
	}
	defer w.LDAPSession.Close()

	if w.Options.Interactive {
		return w.StartTUI()
	} else {
		return w.StartCLI()
	}
}

func (w *WindapSearchSession) StartCLI() error {
	if w.Module == nil {
		fmt.Fprintf(os.Stderr, "[!] You must specify a valid module to use\n")
		fmt.Fprintf(os.Stderr, " Available modules: \n%s", w.ModuleDescriptionString())
		return nil
	}
	err := w.runModule()
	if err != nil {
		return err
	}
	if w.Options.Output != "" {
		fmt.Printf("[+] %s written\n", w.Options.Output)
	}
	return nil
}

func (w *WindapSearchSession) StartTUI() error {
	return nil
}
