package cmd

import (
	"github.com/ropnop/go-windapsearch/pkg/windapsearch"
	"log"
)

func Run() {
	w := windapsearch.NewSession()
	err := w.Run()
	if err != nil {
		log.Fatal(err)
	}
}
