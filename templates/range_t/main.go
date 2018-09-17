package main

import (
	"html/template"
	"log"
	"os"
)

// Person Name
type Person struct {
	Name   string
	Emails []string
}

func main() {
	p := Person{
		Name:   "Ankur",
		Emails: []string{"1@g.com", "2@g.com"},
	}

	// “New” allocates a new template with the given name.
	tmp := template.New("new_html_template")

	// With “range” the current object “.” is set
	// to the successive elements of the array or slice Emails.
	var t = `Hello {{.Name}}
	{{range .Emails}}
	Your Email's {{.}}
	{{end}}
	`
	// “Parse” parses a string into a template.
	tmp, err := tmp.Parse(t)
	if err != nil {
		log.Fatal(err)
		return
	}

	// merge template 'tmpl' with content of 's'
	// “Execute” applies a parsed template to the specified data object,
	// and writes the output to “os.Stdout”.
	err = tmp.Execute(os.Stdout, p)
	if err != nil {
		log.Fatal(err)
		return
	}
}
