package main

import (
	"html/template"
	"log"
	"os"
)

// Person Name
type Person struct {
	Name string
}

func main() {
	p := Person{Name: "Ankur"}

	// “New” allocates a new template with the given name.
	tmp := template.New("new_html_template")

	// “Parse” parses a string into a template.
	tmp, err := tmp.Parse("Hello {{.Name}}")
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
