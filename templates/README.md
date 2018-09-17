## Golang templates

### variable

When the template is rendered in golang, it can accept a variable of type interface{}.

**commonly used types of incoming parameters**

1. struct
2. map[string]interface{}

syntax support embedded in the template, all need to be `{{}}`.

`.` represents the current variable, within the template file.

```Go
type Article struct {
    ArticleId int
    ArticleContent string
}
```

Then we can pass in the template.

```Go
<p>{{.ArticleContent}}<span>{{.ArticleId}}</span></p>
```

**template supports the range loop to traverse the contents of the map and slice.**

```Go
{{range $i, $v := .slice}}
{{end}}

{{range .slice}}
{{end}}

{{range .slice}}
{{.field}}
{{end}}
```

### Nesting of templates

named template that we want to use. This named template needs to be defined inside
the define block as shown below.

```
{{template "navbar"}}
```

This is the named template with name navbar that we can reuse
```
{{define "navbar"}}
{{end}}
```

Content between definitions will override {{template "navbar"}}

get the variable of the parent template

```
{{template "navbar" .}}
```
