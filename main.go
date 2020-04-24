package main

import (
	"fmt"
	"github.com/ropnop/go-windapsearch/pkg/windapsearch"
	"os"
)

func main() {
	w := windapsearch.NewSession()
	err := w.Run()
	if err != nil {
		fmt.Printf("[!] %s\n", err)
		os.Exit(-1)
	}
}