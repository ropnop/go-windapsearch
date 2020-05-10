// +build ignore

// This code consumes the JSON created from scraping the MS documentation and generates the map of Attribute names
// to Syntax
// to update the JSON, run scrapeAttributesFromMS.py
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"text/template"
	"time"
)

type ADAttributeJSON struct {
	CN              string `json:"CN"`
	LdapDisplayName string `json:"Ldap-Display-Name"`
	AttributeId     string `json:"Attribute-Id"`
	SystemIDGuid    string `json:"System-Id-Guid"`
	Syntax          string `json:"Syntax"`
	IsSingleValue   bool   `json:"Is-Single-Valued"`
}

type data struct {
	ADAttributes []ADAttributeJSON
	Syntaxes     []string
	Timestamp    string
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var d data
	d.Timestamp = time.Now().Format(time.RFC3339)

	file, err := ioutil.ReadFile("ADAttributes.json")
	must(err)

	err = json.Unmarshal([]byte(file), &d.ADAttributes)
	must(err)

	uniqueSyntaxes := make(map[string]bool)
	var multiValueAttributes []string

	for _, attr := range d.ADAttributes {
		if _, ok := uniqueSyntaxes[attr.Syntax]; !ok {
			uniqueSyntaxes[attr.Syntax] = true
			d.Syntaxes = append(d.Syntaxes, attr.Syntax)
		}
		if !attr.IsSingleValue {
			multiValueAttributes = append(multiValueAttributes, attr.LdapDisplayName)
		}
	}

	tmpl, err := template.New("attributeSyntax").Parse(attributeSyntaxTemplate)
	must(err)

	f, err := os.Create("attr_syntaxes.go")
	must(err)
	defer f.Close()

	tmpl.Execute(f, d)

}

var attributeSyntaxTemplate = `
// This file was automatically generated at
// {{ .Timestamp }}
// from AD Schema documentation here https://docs.microsoft.com/en-us/windows/win32/adschema/
// length of unique attributes: {{len .ADAttributes}}
// length of unique syntaxes: {{len .Syntaxes}}
package adschema

type ADAttributeInfo struct {
	Syntax string
	IsSingleValue bool
}

var AttributeMap = map[string]*ADAttributeInfo{
{{range $attr := .ADAttributes}}
"{{$attr.LdapDisplayName}}": &ADAttributeInfo{Syntax: "{{$attr.Syntax}}", IsSingleValue: {{$attr.IsSingleValue}}},{{end}}
}
`
