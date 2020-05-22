package windapsearch

import (
	"encoding/json"
	"fmt"
	"github.com/ropnop/go-windapsearch/pkg/adschema"
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"io"
	"os"
	"sync"
)

func (w *WindapSearchSession) outputWorker(input chan []byte, done chan struct{}) {
	defer func() {
		if w.Options.JSON {
			io.WriteString(w.OutputWriter, "]")
		}
		// notify we're done writing
		close(done)
	}()

	entryDelimiter := "\n"
	if w.Options.JSON {
		entryDelimiter = ","
		io.WriteString(w.OutputWriter, "[")
	}
	firstEntry, ok := <-input
	if !ok {
		return
	}
	w.OutputWriter.Write(firstEntry)
	for {
		select {
		case b, ok := <-input:
			if !ok {
				return
			}
			io.WriteString(w.OutputWriter, entryDelimiter)
			w.OutputWriter.Write(b)
		}
	}
}

func (w *WindapSearchSession) searchResultWorker(chans *ldapsession.ResultChannels, out chan []byte, wg *sync.WaitGroup) {
	defer func() {
		fmt.Fprintf(os.Stderr, "worker closing...\n")
		wg.Done()
	}()
	for {
		select {
		case entry, ok := <-chans.Entries:
			if !ok {
				return
			}
			e := &adschema.ADEntry{entry}
			if !w.Options.JSON {
				out <- []byte(e.LDAPFormat())
			} else {
				b, err := json.Marshal(e)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error marshaling %s\n", e.DN)
				}
				out <- b
			}
		// these do nothing, but we need have something receiving these channels, or else the program wil
		case <- chans.Referrals:
		case <- chans.Controls:
			continue
		}
	}
}



func (w *WindapSearchSession) runModule() error {
	var attrs []string
	if w.Options.FullAttributes {
		attrs = []string{"*"}
	} else {
		attrs = w.Options.Attributes
	}


	// Set up our write worker, used to write stuff to stdout or file
	// doneChan is used to indicate the module is completely done and results are written
	doneWriting := make(chan struct{})
	outputChan := make(chan []byte)
	go w.outputWorker(outputChan, doneWriting)

	// set up our result workers, used to translate/marshal entries
	var wg sync.WaitGroup
	for i := 0; i < w.Options.Workers; i++ {
		wg.Add(1)
		go w.searchResultWorker(w.LDAPSession.Channels, outputChan, &wg)
	}

	err := w.Module.Run(w.LDAPSession, attrs)
	if err != nil {
		return err
	}

	// wait for the search to be done and workers to finish
	wg.Wait()
	fmt.Fprintf(os.Stderr, "waitgroup finished!\n")

	// when workers are done, nothing left to write
	close(outputChan)
	<-doneWriting

	return nil
}
