package main

import (
	"os"
	"text/template"
)

func main() {
	temp := `{{"\x1b[1m\x1b[31mColored Text\x1b[0m"}}`
	t := template.Must(template.New("colored text").Parse(temp))
	t.Execute(os.Stdout, nil)
}
