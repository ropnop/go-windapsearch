package main

import (
	"github.com/ropnop/go-windapsearch/pkg/windapsearch"
)

func main() {
	w := windapsearch.NewSession()
	err := w.Run()
	if err != nil {
		w.Log.Fatalf(err.Error())
	}
}
