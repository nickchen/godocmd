package templates

import "text/template"

var templates = map[string]string{"markdown.tmpl": `
# {{.Name}}

import "{{.ImportPath}}"

- [Overview](#Overview)
- [Index](#Index)
- [Examples](#Examples)

## Index
{{if .Consts}}
* [Constants](#Constants){{end}}{{/* END if .Consts */}}
{{if .FuncsFiltered}}
* [Functions](#Functions){{range $func := .FuncsFiltered}}
    * [{{functionSignature .}}]({{anchorFunc "link" .}}){{end}}{{/* END range .FuncsFiltered */}}
{{end}}{{/* END if .Funcs */}}{{if .Types}}
* [Types](#Types){{range $types := .Types}}
    * [type {{ .Name }}]({{anchorFunc "link" .}}){{range $typeFunc := .Funcs }}
        * [{{functionSignature .}}]({{anchorFunc "link" .}}){{end}}{{/* END range .Funcs */}}{{end}}{{/* END range .Types */}}
{{end}}{{/* END if .Types */}}

### Package files

{{ range $file := .Filenames }} {{sourceFileLink .}} {{ end }}

{{if .Consts}}
## Constants
{{range $const := .Consts}}
{{ codeBlock }}
{{ anyTypeSourceString .Decl }}
{{ codeBlock }}
{{end}}{{/* END range .Consts */}}{{end}}{{/* END if .Consts */}}
{{if .FuncsFiltered}}
## Functions
{{range $func := .FuncsFiltered}}
### func {{.Name}}
{{.Doc}}
{{ codeBlock }}
{{functionSignature .}}
{{ codeBlock }}
{{end}}{{/* END range .FuncsFiltered */}}
{{end}}{{/* END if .FuncsFiltered */}}

{{if .Types}}
## Types
{{range $type := .Types}}
### type {{ .Name }}
{{ codeBlock }}
{{ anyTypeSourceString .Decl }}
{{ codeBlock }}
{{ with .Consts }}
{{ range $const := . }}
{{.Doc}}
{{ codeBlock }}
{{ anyTypeSourceString .Decl }}
{{ codeBlock }}

{{ end }} {{/* END range $const := . */}}
{{ end }} {{/* END with .Const */}}

{{ with .Funcs }}
{{ range $func := . }}
#### {{ functionSignature . }}
{{.Doc}}
{{ codeBlock }}
{{ functionSignature . }}
{{ codeBlock }}
{{ range $example := getExampleForFunc $type . }}
##### Example ({{$func.Name}})
{{ codeBlock }}
{{ anyTypeSourceString .Decl }}
{{ codeBlock }}
###### type 2 {{ .Name }}
{{ end }} {{/* END range $example := getExampleForFunc $type . */}}
{{ end }} {{/* END range $func := . */}}
{{ end }} {{/* END with .Funcs */}}

{{ with .Methods }}
{{ range $func := . }}
#### {{ functionSignature . }}
{{.Doc}}
{{ codeBlock }}
{{ functionSignature . }}
{{ codeBlock }}
{{ range $example := getExampleForFunc $type . }}

##### Example ({{$func.Name}})
{{ codeBlock }}
{{ anyTypeSourceString .Decl }}
{{ codeBlock }}
{{ end }} {{/* END with example */}}
{{ end }} {{/* END range .(Methods) */}}{{ end }} {{/* END with .Methods */}}

{{ range $example := getOtherExamplesForType $type }}

##### Example ({{$type.Name}})
{{ codeBlockGolang }}
{{ anyTypeSourceString .Decl }}
{{ codeBlock }}
{{ end }} {{/* END range $example := getOtherExamplesForType $type */}}

{{ end }} {{/* END range .(Types) */}}
{{ end }} {{/* END if .Types */}}
<p align="center" ><small>automatically generated</small></p>
`,
}

// Parse parses declared templates.
func Parse(t *template.Template) (*template.Template, error) {
	for name, s := range templates {
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		if _, err := tmpl.Parse(s); err != nil {
			return nil, err
		}
	}
	return t, nil
}

