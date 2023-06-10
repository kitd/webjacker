package main

import (
	"net/http"
	"net/url"

	"github.com/kitd/webjacker"
)

type AutoCompleter struct {
	*webjacker.HttpResource
}

func (ac *AutoCompleter) Handle(rw http.ResponseWriter, r *http.Request, params url.Values) {
	prefix := params.Get(ac.Id)
	results := searchWords(prefix)
	templates.ExecuteTemplate(rw, "autoc_data", results)
}
