package cmd

import (
	"fmt"
	"github.com/ropnop/go-windapsearch/pkg/modules"
	"github.com/ropnop/go-windapsearch/pkg/utils"
	"github.com/ropnop/go-windapsearch/pkg/windapsession"
	"io"
	"log"
	"os"
	"strings"
)

func init() {

}

func must(err error) {
	if err == nil {
		return
	}
	if strings.Contains(err.Error(), "Invalid Credentials") {
		log.Fatalf("[!] Invalid Credentials")
	}
	if strings.Contains(err.Error(), "to perform this operation a successful bind must be completed") {
		log.Fatal("[!] A successful bind is required for this operation. Please provide valid credentials")
	}
	log.Fatal(err)
}

func Run() {
	if _, err := OptionsParser.Parse(); err != nil {
		os.Exit(1)
	}
	if DomainOptions.Domain == "" && DomainOptions.DomainController == "" {
		fmt.Println("[!] You must specify either a domain or an IP address of a domain controller")
		os.Exit(1)
	}

	password := BindOptions.Password
	var err error
	if BindOptions.Username != "" && password == "" {
		password, err = utils.SecurePrompt(fmt.Sprintf("Password for [%s]", BindOptions.Username))
		must(err)
	}
	options := windapsession.LDAPSessionOptions{
		Domain:           DomainOptions.Domain,
		DomainController: DomainOptions.DomainController,
		Username:         BindOptions.Username,
		Password:         password,
		Port:             BindOptions.Port,
		Secure:           BindOptions.Secure,
	}
	wSession, err := windapsession.NewWindapSession(&options)
	must(err)
	defer wSession.Close()

	var outputWriter io.Writer
	if OutputOptions.Output != "" {
		fd, err := os.Create(OutputOptions.Output)
		must(err)
		outputWriter = fd
	} else {
		outputWriter = os.Stdout
	}
	
	outputOptions := modules.OutputOptions{
		ResolveHosts: OutputOptions.ResolveHosts,
		Attributes:   strings.Split(OutputOptions.Attributes, ","),
		Full:         OutputOptions.FullAttributes,
		JSON:         OutputOptions.JSON,
		Output:       outputWriter,
	}
	
	if EnumerationOptions.Users {
		modules.GetAllUsers.OutputOptions = outputOptions
		modules.GetAllUsers.SetAttrs([]string{"*"})
		wSession.ExecuteModule(modules.GetAllUsers)
	}

	//modules.GetAllUsers.AddFilter("cn=Trevor Hoffman")
	//modules.GetAllUsers.SetAttrs([]string{"*"})

	//userResponse, err := wSession.ExecuteModule(modules.GetAllUsers)
	//must(err)
	////g := userResponse.Entries[0].GetAttributeValue("objectGUID")
	////windapsession.PrintSearchResults(userResponse)
	//jResponse := utils.SearchResultToJSON(userResponse)
	//jData, err := json.Marshal(jResponse)
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Println(string(jData))

}
