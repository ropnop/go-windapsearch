package windapsearch

import (
	"fmt"
	"github.com/ropnop/go-windapsearch/pkg/utils"
	"gopkg.in/ldap.v3"
	"os"
)

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

func (w *WindapSearchSession) printWorker(ch chan *ldap.Entry) {
	for r := range ch {
		utils.WriteEntry(r, w.OutputWriter)
	}
}

func (w *WindapSearchSession) jsonWorker(ch chan *ldap.Entry) {
	// TODO make this work properly and spit out an array
	//io.WriteString(w.OutputWriter, "[")
	for r := range ch {
		//io.WriteString(w.OutputWriter, fmt.Sprintf("TOOD: %s\n", r.DN))
		j, err := utils.EntryToJSON(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[!] Error marshalling %s\n", r.DN)
			continue
		}
		w.OutputWriter.Write(j)
		//io.WriteString(w.OutputWriter, ",")
	}
	//io.WriteString(w.OutputWriter, "]")
}

func (w *WindapSearchSession) runModule() error {
	var attrs []string
	if w.Options.FullAttributes {
		attrs = []string{"*"}
	} else {
		attrs = w.Options.Attributes
	}

	resultsChan := make(chan *ldap.Entry)
	w.LDAPSession.SetChannel(resultsChan)
	if w.Options.JSON {
		go w.jsonWorker(resultsChan)
	} else {
		go w.printWorker(resultsChan)
	}

	err := w.Module.Run(w.LDAPSession, attrs)
	if err != nil  { return err }

	return nil
}
