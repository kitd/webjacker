package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kitd/webjacker"
)

var words []string

func main() {

	wordsRaw, _ := os.ReadFile("words.txt")
	wordsList := bytes.Split(wordsRaw, []byte("\n"))
	words = make([]string, len(wordsList))
	for i, b := range wordsList {
		words[i] = strings.TrimSpace(string(b))
	}

	templates := template.New("Main")
	if _, err := templates.ParseFiles("./index.html"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wordSearch := AutoCompleter{
		webjacker.NewHttpResource("words"),
	}

	if _, err := templates.Parse(autocViews); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wordSearch.On(http.MethodGet, func(rw http.ResponseWriter, r *http.Request, params url.Values) {
		prefix := params.Get(wordSearch.Id)
		results := searchWords(prefix)
		templates.ExecuteTemplate(rw, "autoc_data", results)
	})

	webjacker.RegisterHttpResource(*wordSearch.HttpResource, http.DefaultServeMux)

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if err := templates.ExecuteTemplate(rw, "index", wordSearch); err != nil {
			log.Fatal(err)
		}
	})

	http.ListenAndServe(":8080", nil)
}

func searchWords(prefix string) []string {
	var results []string
	pfx := strings.ToLower(prefix)
	for _, w := range words {
		word := strings.ToLower(w)
		if strings.HasPrefix(word, pfx) {
			results = append(results, w)
		}
	}
	return results
}
