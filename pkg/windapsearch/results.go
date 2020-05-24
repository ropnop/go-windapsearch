package windapsearch

import (
	"encoding/json"
	"github.com/ropnop/go-windapsearch/pkg/adschema"
	"github.com/ropnop/go-windapsearch/pkg/ldapsession"
	"io"
	"sync"
)

func (w *WindapSearchSession) outputWorker(input chan []byte, done chan struct{}) {
	w.Log.Debugf("outputWorker started")
	defer func() {
		if w.Options.JSON {
			io.WriteString(w.OutputWriter, "]")
		}
		// notify we're done writing by closing channel
		close(done)
		w.Log.Debugf("outputWorker closing, finished writing")
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
	for b := range input {
		io.WriteString(w.OutputWriter, entryDelimiter)
		w.OutputWriter.Write(b)
	}
}

func (w *WindapSearchSession) searchResultWorker(chans *ldapsession.ResultChannels, out chan []byte, wg *sync.WaitGroup) {
	w.Log.Debugf("searchResultsWorker started")
	defer func() {
		w.Log.Debugf("searchResultsWorker closing")
		wg.Done()
	}()
	for {
		select {
		case entry, ok := <-chans.Entries:
			if !ok {
				return
			}
			w.Log.WithField("DN", entry.DN).Debug("parsing entry")
			e := &adschema.ADEntry{entry}
			if !w.Options.JSON {
				out <- []byte(e.LDAPFormat())
			} else {
				b, err := json.Marshal(e)
				if err != nil {
					w.Log.WithField("DN", e.DN).Warn("error marshaling entry")
				}
				out <- b
			}
		// these do nothing, but we need have something receiving these channels, or else the program will freeze
		case <-chans.Referrals:
		case <-chans.Controls:
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
	for i := 0; i < w.workers; i++ {
		wg.Add(1)
		go w.searchResultWorker(w.LDAPSession.Channels, outputChan, &wg)
	}

	err := w.Module.Run(w.LDAPSession, attrs)
	if err != nil {
		return err
	}

	// wait for the search to be done and workers to finish
	wg.Wait()
	w.Log.Debug("waitgroup finished, all entry workers done")

	// when workers are done, nothing left to write
	close(outputChan)
	w.Log.Debug("output channel closed. waiting for writer to finish")

	<-doneWriting

	return nil
}
