package windapsession

import (
	"fmt"
	"github.com/ropnop/go-windapsearch/pkg/modules"
	"github.com/ropnop/go-windapsearch/pkg/utils"
	"gopkg.in/ldap.v3"
)

func (w *Windapsession) MakeSimpleSearchRequest(filter string, attrs []string) *ldap.SearchRequest {
	return ldap.NewSearchRequest(
		w.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		filter,
		attrs,
		nil)
}

func (w *Windapsession) GetSearchResults(request *ldap.SearchRequest) (result *ldap.SearchResult, err error) {
	return w.lConn.SearchWithPaging(request, 1000)
}

func (w *Windapsession) ExecuteModule(mod modules.WindapModule) (err error) {
	req := w.MakeSimpleSearchRequest(mod.Filter, mod.DefaultAttrs)
	rawResults, err := w.GetSearchResults(req)
	if err != nil {
		return
	}
	var results []byte
	if mod.OutputOptions.JSON {
		results, err = utils.SearchResultToJSON(rawResults)
		if err != nil {
			return
		}
	} else {
		results = []byte("Not implemented yet")
	}
	mod.OutputOptions.Output.Write(results)
	return nil
}

func PrintSearchResults(result *ldap.SearchResult) {
	for _, entry := range result.Entries {
		for _, attribute := range entry.Attributes {
			for _, value := range attribute.Values {
				fmt.Printf("%v: %v\n", attribute.Name, value)
			}
		}
		fmt.Println()
	}
}
