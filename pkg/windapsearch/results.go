package windapsearch

import (
	"encoding/json"
	"fmt"
	"github.com/ropnop/go-windapsearch/pkg/adschema"
	"gopkg.in/ldap.v3"
	"io"
	"os"
	"sync"
)




func (w *WindapSearchSession) outputWorker(in chan []byte) {
	entryDelimiter := "\n"
	if w.Options.JSON {
		entryDelimiter = ","
		io.WriteString(w.OutputWriter, "[")
	}
	firstEntry := <-in
	w.OutputWriter.Write(firstEntry)
	for b := range in {
		io.WriteString(w.OutputWriter, entryDelimiter)
		w.OutputWriter.Write(b)
	}
	if w.Options.JSON {
		io.WriteString(w.OutputWriter, "]")
	}
	w.doneChan <- true
}

func (w *WindapSearchSession) stringWorker(in chan *ldap.Entry, out chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	for r := range in {
		e := &adschema.ADEntry{r}
		out <- []byte(e.LDAPFormat())
	}
}

func (w *WindapSearchSession) jsonWorker(in chan *ldap.Entry, out chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	for r := range in {
		e := &adschema.ADEntry{r}
		b, err := json.Marshal(e)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error marshaling %s\n", e.DN)
		}
		out <- b
	}
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

	var wg sync.WaitGroup
	w.doneChan = make(chan bool, 1)

	worker := w.stringWorker
	if w.Options.JSON {
		worker = w.jsonWorker
	}
	outputChan := make(chan []byte)
	go w.outputWorker(outputChan)
	for i := 0; i <= w.Workers; i ++ {
		wg.Add(1)
		go worker(resultsChan, outputChan, &wg)
	}

	err := w.Module.Run(w.LDAPSession, attrs)
	if err != nil  { return err }

	wg.Wait()
	close(outputChan)
	<-w.doneChan

	return nil
}
