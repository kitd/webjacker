package main

import (
	"github.com/kitd/webjacker"
)

var autocViews string = `
	{{define "autoc_data" -}}
		{{- range . -}}
			<option>{{.}}</option>
		{{- end -}}
	{{- end}}

	{{define "autoc" -}}
		<input name="{{.Id}}" id="{{.Id}}" type="text" list="{{.Id}}_data"
			hx-get="{{.Path}}" 
			hx-trigger="keyup delay:500ms"
			hx-target="#{{.Id}}_data">
		<datalist id="{{.Id}}_data">
			{{- template "autoc_data" -}}
		</datalist>
	{{- end}}`

type AutoCompleter struct {
	*webjacker.HttpResource
}
