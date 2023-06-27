package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/kitd/webjacker"
)

var (
	//go:embed templates
	templateFiles embed.FS
	words         []string
	templates     *template.Template
)

// An autocompletion service
type AutoCompleter struct {
	*webjacker.HttpResource
}

func main() {

	// Load the words list
	wordsRaw, _ := os.ReadFile("words.txt")
	wordsList := bytes.Split(wordsRaw, []byte("\n"))
	words = make([]string, len(wordsList))
	for i, b := range wordsList {
		words[i] = strings.TrimSpace(string(b))
	}

	// Load the templates
	templates = template.New("Main")
	if _, err := templates.ParseFS(templateFiles, "templates/*.gohtml"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a word search autocompletion service named `words`
	wordSearch := AutoCompleter{
		webjacker.NewHttpResource("words"),
	}

	// On GET, retrieve the input box text, look up the words prefixed by that text,
	// then run the results through the data template and send to the response
	wordSearch.On(http.MethodGet,
		func(rw http.ResponseWriter, r *http.Request, params url.Values) {
			prefix := params.Get(wordSearch.Id)
			results := searchWords(prefix)
			templates.ExecuteTemplate(rw, "autoc_data", results)
		})

	// Load our HttpResource. It will be available on the path `/words`
	webjacker.RegisterHttpResource(wordSearch.HttpResource)

	// Handle the main page
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if err := templates.ExecuteTemplate(rw, "index", wordSearch); err != nil {
			log.Fatal(err)
		}
	})

	// Run the server
	http.ListenAndServe(":8080", nil)
}

// Search the words. Note that we need to handle the case of an empty string, which might occur if the user
// types a backspace to remove the first char, leaving an empty input box. In this case, return no results.
func searchWords(prefix string) []string {
	var results []string
	if prefix != "" {
		pfx := strings.ToLower(prefix)
		for _, w := range words {
			word := strings.ToLower(w)
			if strings.HasPrefix(word, pfx) {
				results = append(results, w)
			}
		}
	}
	return results
}
